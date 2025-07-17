package sysdig

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sysdiglabs/agent-kilt/pkg/cfnpatcher"
)

const agentinoKiltDefinition = `build {
    entry_point: ["/opt/draios/bin/instrument"]
    command: ${?original.entry_point} ${?original.command}
    environment_variables: {
        "SYSDIG_ORCHESTRATOR": ${config.orchestrator_host}
        "SYSDIG_ORCHESTRATOR_PORT": ${config.orchestrator_port}
        "SYSDIG_COLLECTOR": ${config.collector_host}
        "SYSDIG_COLLECTOR_PORT": ${config.collector_port}
        "SYSDIG_ACCESS_KEY": ${config.sysdig_access_key}
        "SYSDIG_LOGGING": ${config.sysdig_logging}
        "SYSDIG_SIDECAR": ${config.sidecar}
        "SYSDIG_PRIORITY": ${config.priority}
    }
    capabilities: ["SYS_PTRACE"]
    mount: [
        {
            name: "SysdigInstrumentation"
            image: ${config.agent_image}
            volumes: ["/opt/draios"]
            entry_point: ["/opt/draios/bin/logwriter"]
        }
    ]
}`

func dataSourceSysdigFargateWorkloadAgent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigFargateWorkloadAgentRead,

		Schema: map[string]*schema.Schema{
			"container_definitions": {
				Type:        schema.TypeString,
				Description: "the input Fargate container definitions to instrument with the Sysdig workload agent",
				Required:    true,
			},
			"sysdig_access_key": {
				Type:        schema.TypeString,
				Description: "the Sysdig access key",
				Optional:    true,
			},
			"workload_agent_image": {
				Type:        schema.TypeString,
				Description: "the Sysdig workload agent image",
				Required:    true,
			},
			"image_auth_secret": {
				Type:        schema.TypeString,
				Description: "registry authentication secret",
				Optional:    true,
			},
			"orchestrator_host": {
				Type:        schema.TypeString,
				Description: "the orchestrator host to connect to",
				Optional:    true,
				Deprecated:  "The orchestrator agent is deprecated and no longer supported as of workload-agent version 6.0.0. Please use the collector_host parameter instead.",
			},
			"orchestrator_port": {
				Type:        schema.TypeString,
				Description: "the orchestrator port to connect to",
				Optional:    true,
				Deprecated:  "The orchestrator agent is deprecated and no longer supported as of workload-agent version 6.0.0. Please use the collector_port parameter instead.",
			},
			"collector_host": {
				Type:        schema.TypeString,
				Description: "the collector host to connect to",
				Optional:    true,
			},
			"collector_port": {
				Type:        schema.TypeString,
				Description: "the collector port to connect to",
				Optional:    true,
			},
			"ignore_containers": {
				Type:        schema.TypeList,
				Description: "list of containers to not add instrumentation to",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"bare_pdig_on_containers": {
				Type:        schema.TypeList,
				Description: "use bare pdig to instrument the containers in the list",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"log_configuration": {
				Type:        schema.TypeSet,
				MaxItems:    1,
				Description: "configuration for instrumentation logs using the awslogs driver",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:        schema.TypeString,
							Description: "The log group where the awslogs driver will send log streams",
							Required:    true,
						},
						"stream_prefix": {
							Type:        schema.TypeString,
							Description: "Prefix for the instrumentation log stream",
							Required:    true,
						},
						"region": {
							Type:        schema.TypeString,
							Description: "Region for the log group",
							Required:    true,
						},
					},
				},
			},
			"sysdig_logging": {
				Type:        schema.TypeString,
				Description: "the instrumentation logging level",
				Optional:    true,
			},
			"sidecar": {
				Type:        schema.TypeString,
				Description: "Sidecar mode: auto/force/(empty string)",
				Default:     "auto",
				Optional:    true,
			},
			"priority": {
				Type:        schema.TypeString,
				Description: "The priority of the agent. Can be 'security' or 'availability'",
				Default:     "availability",
				Optional:    true,
			},
			"instrumentation_essential": {
				Type:        schema.TypeBool,
				Description: "Should the instrumentation container be marked as essential",
				Default:     false,
				Optional:    true,
			},
			"instrumentation_cpu": {
				Type:        schema.TypeInt,
				Description: "The number of cpu units dedicated to the instrumentation container",
				Default:     0,
				Optional:    true,
			},
			"instrumentation_memory_limit": {
				Type:        schema.TypeInt,
				Description: "The maximum amount (in MiB) of memory used by the instrumentation container",
				Default:     0,
				Optional:    true,
			},
			"instrumentation_memory_reservation": {
				Type:        schema.TypeInt,
				Description: "The minimum amount (in MiB) of memory reserved for the instrumentation container",
				Default:     0,
				Optional:    true,
			},

			"output_container_definitions": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type cfnTag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type cfnProperties struct {
	RequiresCompatibilities []string         `json:"RequiresCompatibilities"`
	ContainerDefinitions    []map[string]any `json:"ContainerDefinitions"`
	Tags                    []cfnTag         `json:"Tags"`
}

type cfnResource struct {
	ResourceType string        `json:"Type"`
	Properties   cfnProperties `json:"Properties"`
}

type cfnStack struct {
	Resources map[string]cfnResource `json:"Resources"`
}

func contains(items []string, target string) bool {
	return slices.Contains(items, target)
}

// fargatePostKiltModifications performs any additional changes needed after Kilt has applied it's transformations
func fargatePostKiltModifications(patchedBytes []byte, patchOpts *patchOptions) ([]byte, error) {
	if len(patchOpts.LogConfiguration) == 0 && len(patchOpts.BarePdigOnContainers) == 0 {
		// nothing to do
		return patchedBytes, nil
	}

	containers, err := gabs.ParseJSON(patchedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse containers for post-processing: %s", err)
	}

	for _, container := range containers.Children() {
		// Skip unnamed containers
		// Note that lowercase "name" tags have been replaced by "Name" during TaskDefinition patching
		containerName, ok := container.Search("Name").Data().(string)
		if !ok {
			continue
		}

		if containerName == "SysdigInstrumentation" {
			// Add log configuration to the SysdigInstrumentation sidecar container
			if len(patchOpts.LogConfiguration) != 0 {
				awsLogConfig := &ecs.LogConfiguration{
					LogDriver: aws.String("awslogs"),
					Options: map[string]*string{
						"awslogs-group":         aws.String(patchOpts.LogConfiguration["group"].(string)),
						"awslogs-stream-prefix": aws.String(patchOpts.LogConfiguration["stream_prefix"].(string)),
						"awslogs-region":        aws.String(patchOpts.LogConfiguration["region"].(string)),
					},
				}
				_, err = container.Set(awsLogConfig, "LogConfiguration")
				if err != nil {
					return nil, fmt.Errorf("failed to set log configuration: %s", err)
				}
			}

			if !patchOpts.Essential {
				_, err := container.Set(false, "essential")
				if err != nil {
					return nil, fmt.Errorf("failed to set essential flag: %s", err)
				}
			}

			if patchOpts.CPUShares != 0 {
				_, err := container.Set(patchOpts.CPUShares, "cpu")
				if err != nil {
					return nil, fmt.Errorf("failed to set cpu shares: %s", err)
				}
			}

			if patchOpts.MemoryLimit != 0 {
				_, err := container.Set(patchOpts.MemoryLimit, "memory")
				if err != nil {
					return nil, fmt.Errorf("failed to set memory limit: %s", err)
				}
			}

			if patchOpts.MemoryReservation != 0 {
				_, err := container.Set(patchOpts.MemoryReservation, "memoryReservation")
				if err != nil {
					return nil, fmt.Errorf("failed to set memory reservation: %s", err)
				}
			}
		} else {
			// Use bare pdig in the current workload container if instrumented
			if contains(patchOpts.BarePdigOnContainers, containerName) && !contains(patchOpts.IgnoreContainers, containerName) {
				envars := map[string]any{
					"Name":  "__INSTRUMENTATION_WRAPPER",
					"Value": "/opt/draios/bin/pdig,-C,-t,-1",
				}
				err := container.ArrayAppend(envars, "Environment")
				if err != nil {
					return nil, fmt.Errorf("failed to extend environment variables: %s", err)
				}
			}
		}
	}

	return containers.Bytes(), nil
}

// PatchFargateTaskDefinition modifies the container definitions
func patchFargateTaskDefinition(ctx context.Context, containerDefinitions string, kiltConfig *cfnpatcher.Configuration, patchOpts *patchOptions) (patched *string, err error) {
	var cdefs []map[string]any
	err = json.Unmarshal([]byte(containerDefinitions), &cdefs)
	if err != nil {
		return nil, err
	}

	// Convert the ignore containers list into Kilt tags for the patcher
	tags := []cfnTag{}
	if len(patchOpts.IgnoreContainers) > 0 {
		containerTagValue := strings.Join(patchOpts.IgnoreContainers, ":")
		tags = append(tags, cfnTag{
			Key:   "kilt-ignore-containers",
			Value: containerTagValue,
		})
	}

	stack := cfnStack{
		Resources: map[string]cfnResource{
			"kilt": {
				ResourceType: "AWS::ECS::TaskDefinition",
				Properties: cfnProperties{
					RequiresCompatibilities: []string{"FARGATE"},
					ContainerDefinitions:    cdefs,
					Tags:                    tags,
				},
			},
		},
	}

	patchedStack, err := json.Marshal(stack)
	if err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			patched = nil
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation
				err = errors.New("unknown panic")
			}
		}
	}()

	// ECS JSON modifications
	patchedStack, _ = terraformPreModifications(ctx, patchedStack)

	templateParameters := make([]byte, 0)
	patchedBytes, err := cfnpatcher.Patch(ctx, kiltConfig, patchedStack, templateParameters)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(patchedBytes, &stack)
	if err != nil {
		return nil, err
	}

	patchedBytes, err = json.Marshal(stack.Resources["kilt"].Properties.ContainerDefinitions)
	if err != nil {
		return nil, err
	}

	patchedBytes, err = fargatePostKiltModifications(patchedBytes, patchOpts)

	patchedString := string(patchedBytes)
	return &patchedString, nil
}

type KiltRecipeConfig struct {
	SysdigAccessKey  string `json:"sysdig_access_key"`
	AgentImage       string `json:"agent_image"`
	OrchestratorHost string `json:"orchestrator_host"`
	OrchestratorPort string `json:"orchestrator_port"`
	CollectorHost    string `json:"collector_host"`
	CollectorPort    string `json:"collector_port"`
	SysdigLogging    string `json:"sysdig_logging"`
	Sidecar          string `json:"sidecar"`
	Priority         string `json:"priority"`
}

type patchOptions struct {
	BarePdigOnContainers []string
	IgnoreContainers     []string
	LogConfiguration     map[string]any
	Essential            bool
	CPUShares            int
	MemoryLimit          int
	MemoryReservation    int
}

func newPatchOptions(d *schema.ResourceData) *patchOptions {
	opts := &patchOptions{
		BarePdigOnContainers: []string{},
		IgnoreContainers:     []string{},
		LogConfiguration:     map[string]any{},
	}

	if items := d.Get("bare_pdig_on_containers"); items != nil {
		for _, itemRaw := range items.([]any) {
			if itemStr, ok := itemRaw.(string); ok {
				opts.BarePdigOnContainers = append(opts.BarePdigOnContainers, strings.TrimSpace(itemStr))
			}
		}
	}

	if items := d.Get("ignore_containers"); items != nil {
		for _, itemRaw := range items.([]any) {
			if itemStr, ok := itemRaw.(string); ok {
				opts.IgnoreContainers = append(opts.IgnoreContainers, strings.TrimSpace(itemStr))
			}
		}
	}

	if logConfiguration := d.Get("log_configuration").(*schema.Set).List(); len(logConfiguration) > 0 {
		opts.LogConfiguration = logConfiguration[0].(map[string]any)
	}

	if essential := d.Get("instrumentation_essential"); essential != nil {
		opts.Essential = essential.(bool)
	} else {
		priority := d.Get("priority").(string)
		opts.Essential = priority == "security"
	}

	if cpuShares := d.Get("instrumentation_cpu"); cpuShares != nil {
		opts.CPUShares = cpuShares.(int)
	} else {
		opts.CPUShares = 0
	}

	if memoryLimit := d.Get("instrumentation_memory_limit"); memoryLimit != nil {
		opts.MemoryLimit = memoryLimit.(int)
	} else {
		opts.MemoryLimit = 0
	}

	if memoryReservation := d.Get("instrumentation_memory_reservation"); memoryReservation != nil {
		opts.MemoryReservation = memoryReservation.(int)
	} else {
		opts.MemoryReservation = 0
	}

	return opts
}

func dataSourceSysdigFargateWorkloadAgentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	priority := d.Get("priority").(string)
	if priority != "security" && priority != "availability" {
		return diag.Errorf("Invalid priority: %s. must be either \"security\" or \"availability\"", priority)
	}

	recipeConfig := KiltRecipeConfig{
		SysdigAccessKey:  d.Get("sysdig_access_key").(string),
		AgentImage:       d.Get("workload_agent_image").(string),
		OrchestratorHost: d.Get("orchestrator_host").(string),
		OrchestratorPort: d.Get("orchestrator_port").(string),
		CollectorHost:    d.Get("collector_host").(string),
		CollectorPort:    d.Get("collector_port").(string),
		SysdigLogging:    d.Get("sysdig_logging").(string),
		Sidecar:          d.Get("sidecar").(string),
		Priority:         priority,
	}

	jsonConf, err := json.Marshal(&recipeConfig)
	if err != nil {
		return diag.Errorf("Failed to serialize configuration: %v", err.Error())
	}

	scObj := gabs.New()
	imageAuth := d.Get("image_auth_secret").(string)
	if imageAuth != "" {
		_, err := scObj.Set(imageAuth, "RepositoryCredentials", "CredentialsParameter")
		if err != nil {
			return diag.Errorf("cannot set image auth secret in sidecar config: %v", err.Error())
		}
	}

	sc, err := json.Marshal(scObj)
	if err != nil {
		panic("cannot marshal sidecar config: " + err.Error())
	}
	sidecarConfig := string(sc)

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConf),
		SidecarConfig:      sidecarConfig,
	}

	containerDefinitions := d.Get("container_definitions").(string)

	patchOpts := newPatchOptions(d)

	outputContainerDefinitions, err := patchFargateTaskDefinition(ctx, containerDefinitions, kiltConfig, patchOpts)
	if err != nil {
		return diag.Errorf("Error applying configuration patch: %v", err.Error())
	}

	cdefChecksum := sha256.Sum256([]byte(containerDefinitions))
	d.SetId(fmt.Sprintf("%x", cdefChecksum))
	_ = d.Set("output_container_definitions", *outputContainerDefinitions)
	return nil
}

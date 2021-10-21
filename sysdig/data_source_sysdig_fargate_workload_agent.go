package sysdig

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
        "SYSDIG_LOGGING": ""
    }
    mount: [
        {
            name: "SysdigInstrumentation"
            image: ${config.agent_image}
            volumes: ["/opt/draios"]
            entry_point: ["/opt/draios/bin/waitforever"]
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
				Required:    true,
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
			},
			"orchestrator_port": {
				Type:        schema.TypeString,
				Description: "the orchestrator port to connect to",
				Optional:    true,
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
			"output_container_definitions": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

type cfnProperties struct {
	RequiresCompatibilities []string                 `json:"RequiresCompatibilities"`
	ContainerDefinitions    []map[string]interface{} `json:"ContainerDefinitions"`
}

type cfnResource struct {
	ResourceType string        `json:"Type"`
	Properties   cfnProperties `json:"Properties"`
}

type cfnStack struct {
	Resources map[string]cfnResource `json:"Resources"`
}

func patchFargateTaskDefinition(ctx context.Context, containerDefinitions string, kiltConfig *cfnpatcher.Configuration) (patched *string, err error) {
	var cdefs []map[string]interface{}
	err = json.Unmarshal([]byte(containerDefinitions), &cdefs)
	if err != nil {
		return nil, err
	}

	stack := cfnStack{
		Resources: map[string]cfnResource{
			"kilt": {
				ResourceType: "AWS::ECS::TaskDefinition",
				Properties: cfnProperties{
					RequiresCompatibilities: []string{"FARGATE"},
					ContainerDefinitions:    cdefs,
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

	patchedBytes, err := cfnpatcher.Patch(ctx, kiltConfig, patchedStack)
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
}

func dataSourceSysdigFargateWorkloadAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	recipeConfig := KiltRecipeConfig{
		SysdigAccessKey:  d.Get("sysdig_access_key").(string),
		AgentImage:       d.Get("workload_agent_image").(string),
		OrchestratorHost: d.Get("orchestrator_host").(string),
		OrchestratorPort: d.Get("orchestrator_port").(string),
		CollectorHost:    d.Get("collector_host").(string),
		CollectorPort:    d.Get("collector_port").(string),
	}

	jsonConf, err := json.Marshal(&recipeConfig)
	if err != nil {
		return diag.Errorf("Failed to serialize configuration: %v", err.Error())
	}

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		ImageAuthSecret:    d.Get("image_auth_secret").(string),
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConf),
	}

	containerDefinitions := d.Get("container_definitions").(string)

	outputContainerDefinitions, err := patchFargateTaskDefinition(ctx, containerDefinitions, kiltConfig)
	if err != nil {
		return diag.Errorf("Error applying configuration patch: %v", err.Error())
	}

	cdefChecksum := sha256.Sum256([]byte(containerDefinitions))
	d.SetId(fmt.Sprintf("%x", cdefChecksum))
	err = d.Set("output_container_definitions", *outputContainerDefinitions)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

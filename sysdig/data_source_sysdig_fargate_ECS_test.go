//go:build unit

package sysdig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/sysdiglabs/agent-kilt/runtimes/cloudformation/cfnpatcher"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// sortContainerEnv goes into a container definition and sorts the environment variables
func sortContainerEnv(json []byte) string {
	jsonObject, _ := gabs.ParseJSON(json)
	containers, _ := jsonObject.Data().([]interface{})
	for _, container := range containers {
		if env, ok := container.(map[string]interface{})["Environment"]; ok {
			envSort := env.([]interface{})
			sort.Slice(envSort, func(i, j int) bool {
				return gabs.Wrap(envSort[i]).S("Name").Data().(string) < gabs.Wrap(envSort[j]).S("Name").Data().(string)
			})
		}
	}
	return jsonObject.String()
}

func sortAndCompare(t *testing.T, expected []byte, actual []byte) {
	expectedJSON := sortContainerEnv(expected)
	actualJSON := sortContainerEnv(actual)
	assert.JSONEq(t, expectedJSON, actualJSON)
}

// getKiltRecipe returns the default json Kilt recipe
func getKiltRecipe(t *testing.T) string {
	recipeConfig := KiltRecipeConfig{
		SysdigAccessKey:  "sysdig_access_key",
		AgentImage:       "workload_agent_image",
		OrchestratorHost: "orchestrator_host",
		OrchestratorPort: "orchestrator_port",
		CollectorHost:    "collector_host",
		CollectorPort:    "collector_port",
		SysdigLogging:    "sysdig_logging",
	}

	jsonRecipeConfig, err := json.Marshal(&recipeConfig)
	if err != nil {
		t.Fatalf("Failed to serialize configuration: %v", err.Error())
	}

	return string(jsonRecipeConfig)
}

func TestContains(t *testing.T) {
	tests := []struct {
		slice  []string
		target string
		result bool
	}{
		{
			slice:  []string{"gimme", "fried", "chicken"},
			target: "chicken",
			result: true,
		},
		{
			slice:  []string{"the", "answer", "is"},
			target: "42",
			result: false,
		},
		{
			slice:  []string{""},
			target: "empty",
			result: false,
		},
	}
	for idx, tc := range tests {
		t.Run(fmt.Sprintf("%d", idx), func(t *testing.T) {
			result := contains(tc.slice, tc.target)
			assert.Equal(t, tc.result, result, "Error, expected: %t, got: %t", tc.result, result)
		})
	}
}

func TestNewPatchOptions(t *testing.T) {
	newMockResource := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ignore_containers": {
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
				"bare_pdig_on_containers": {
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
				"log_configuration": {
					Type:     schema.TypeSet,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"group": {
								Type:     schema.TypeString,
								Required: true,
							},
							"stream_prefix": {
								Type:     schema.TypeString,
								Required: true,
							},
							"region": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
			},
		}
	}

	// Create a mock resource
	resource := newMockResource()
	data := resource.Data(nil)

	var err error
	err = data.Set("bare_pdig_on_containers", []interface{}{
		"gimme", "fried", "chicken",
	})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Could not set bare_pdig_on_containers, got error: %v", err))
	}

	err = data.Set("ignore_containers", []interface{}{
		"gimme", "fried", "chicken",
	})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Could not set ignore_containers, got error: %v", err))
	}

	err = data.Set("log_configuration", []interface{}{
		map[string]interface{}{
			"group":         "gimme",
			"stream_prefix": "fried",
			"region":        "chicken",
		},
	})
	if err != nil {
		assert.FailNow(t, fmt.Sprintf("Could not set log_configuration, got error: %v", err))
	}

	// Expected vs actual
	expectedPatchOptions := &patchOptions{
		BarePdigOnContainers: []string{"gimme", "fried", "chicken"},
		IgnoreContainers:     []string{"gimme", "fried", "chicken"},
		LogConfiguration: map[string]interface{}{
			"group":         "gimme",
			"stream_prefix": "fried",
			"region":        "chicken",
		},
		Essential: true,
	}
	actualPatchOptions := newPatchOptions(data)

	if !reflect.DeepEqual(expectedPatchOptions, actualPatchOptions) {
		t.Errorf("patcConfigurations are not equal. Expected: %v, Actual: %v", expectedPatchOptions, actualPatchOptions)
	}
}

func getSidecarConfig() string {
	scObj := gabs.New()
	_, err := scObj.Set("image_auth_secret", "RepositoryCredentials", "CredentialsParameter")
	if err != nil {
		panic("cannot set image auth secret in sidecar config: " + err.Error())
	}
	sc, _ := json.Marshal(scObj)
	return string(sc)
}

func TestECStransformation(t *testing.T) {
	inputfile, err := os.ReadFile("testfiles/ECSinput.json")
	if err != nil {
		t.Fatalf("Cannot find testfiles/ECSinput.json")
	}

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       getKiltRecipe(t),
		SidecarConfig:      getSidecarConfig(),
	}

	patchOpts := &patchOptions{}

	patchedOutput, err := patchFargateTaskDefinition(context.Background(), string(inputfile), kiltConfig, patchOpts)
	if err != nil {
		t.Fatalf("Cannot execute PatchFargateTaskDefinition : %v", err.Error())
	}

	expectedOutput, err := os.ReadFile("testfiles/ECSInstrumented.json")
	if err != nil {
		t.Fatalf("Cannot find testfiles/ECSinput.json")
	}

	sortAndCompare(t, expectedOutput, []byte(*patchedOutput))
}

func TestPatchFargateTaskDefinition(t *testing.T) {
	// Kilt Configuration, test invariant
	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       getKiltRecipe(t),
		SidecarConfig:      getSidecarConfig(),
	}

	// File readers
	readFile := func(fileName string) string {
		content, _ := os.ReadFile("testfiles/" + fileName + ".json")
		return string(content)
	}

	getContainerDefinitionOriginal := func(fileName string) string {
		return readFile(fileName)
	}

	getContainerDefinitionPatched := func(fileName string) string {
		return readFile(fileName + "_expected")
	}

	tests := []struct {
		testName  string
		patchOpts *patchOptions
	}{
		{
			testName:  `fargate_entrypoint_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_env_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_cmd_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_linuxparameters_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_combined_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_volumesfrom_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName:  `fargate_field_case_test`,
			patchOpts: &patchOptions{},
		},
		{
			testName: `fargate_log_group`,
			patchOpts: &patchOptions{
				LogConfiguration: map[string]interface{}{
					"group":         "test_log_group",
					"stream_prefix": "test_prefix",
					"region":        "test_region",
				},
				Essential: true,
			},
		},
		{
			testName: `fargate_ignore_container_test`,
			patchOpts: &patchOptions{
				IgnoreContainers: []string{"other", "another"},
			},
		},
		{
			testName: `fargate_bare_pdig`,
			patchOpts: &patchOptions{
				BarePdigOnContainers: []string{"barePdig"},
				IgnoreContainers:     []string{"skipped"},
				Essential:            true,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			patched, err := patchFargateTaskDefinition(
				context.Background(),
				getContainerDefinitionOriginal(tc.testName),
				kiltConfig,
				tc.patchOpts)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Could not patch task definition, got error: %v", err))
			}
			sortAndCompare(t, []byte(getContainerDefinitionPatched(tc.testName)), []byte(*patched))
		})
	}
}

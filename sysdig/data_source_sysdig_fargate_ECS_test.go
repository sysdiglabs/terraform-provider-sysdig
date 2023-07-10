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
	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/stretchr/testify/assert"

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

func TestNewPatchOptions(t *testing.T) {
	newMockResource := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ignore_containers": {
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
	data.Set("ignore_containers", []interface{}{
		"gimme", "fried", "chicken",
	})
	data.Set("log_configuration", []interface{}{
		map[string]interface{}{
			"group":         "gimme",
			"stream_prefix": "fried",
			"region":        "chicken",
		},
	})

	// Expected vs actual
	expectedPatchOptions := &patchOptions{
		IgnoreContainers: []string{"gimme", "fried", "chicken"},
		LogConfiguration: map[string]interface{}{
			"group":         "gimme",
			"stream_prefix": "fried",
			"region":        "chicken",
		},
	}
	actualPatchOptions := newPatchOptions(data)

	if !reflect.DeepEqual(expectedPatchOptions, actualPatchOptions) {
		t.Errorf("patcConfigurations are not equal. Expected: %v, Actual: %v", expectedPatchOptions, actualPatchOptions)
	}
}

func TestECStransformation(t *testing.T) {
	inputfile, err := os.ReadFile("testfiles/ECSinput.json")

	if err != nil {
		t.Fatalf("Cannot find testfiles/ECSinput.json")
	}

	recipeConfig := KiltRecipeConfig{
		SysdigAccessKey:  "sysdig_access_key",
		AgentImage:       "workload_agent_image",
		OrchestratorHost: "orchestrator_host",
		OrchestratorPort: "orchestrator_port",
		CollectorHost:    "collector_host",
		CollectorPort:    "collector_port",
		SysdigLogging:    "sysdig_logging",
	}

	jsonConf, err := json.Marshal(&recipeConfig)
	if err != nil {
		t.Fatalf("Failed to serialize configuration: %v", err.Error())
	}

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		ImageAuthSecret:    "image_auth_secret",
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConf),
	}

	ignoreContainers := []string{}

	patchedOutput, err := patchFargateTaskDefinition(context.Background(), string(inputfile), kiltConfig, nil, &ignoreContainers)
	if err != nil {
		t.Fatalf("Cannot execute PatchFargateTaskDefinition : %v", err.Error())
	}

	expectedOutput, err := os.ReadFile("testfiles/ECSInstrumented.json")
	if err != nil {
		t.Fatalf("Cannot find testfiles/ECSinput.json")
	}

	type cdef struct {
		Command          []string            `json:"Command"`
		EntryPoint       []string            `json:"EntryPoint"`
		Environment      []map[string]string `json:"Environment"`
		Image            string              `json:"Image"`
		Linuxparameters  interface{}         `json:"LinuxParameters"`
		VolumesFrom      []interface{}       `json:"VolumesFrom"`
		LogConfiguration interface{}         `json:"LogConfiguration"`
		Name             string              `json:"Name"`
		Name2            string              `json:"name"`
		Image2           string              `json:"image"`
		EntryPoint2      string              `json:"entryPoint"`
	}

	var patchedContainerDefinitions, expectedContainerDefinitions []cdef
	err = json.Unmarshal([]byte(*patchedOutput), &patchedContainerDefinitions)
	if err != nil {
		t.Fatalf("Error Unmarshaling patched Container definitions : %v", err.Error())
	}

	err = json.Unmarshal([]byte(expectedOutput), &expectedContainerDefinitions)
	if err != nil {
		t.Fatalf("Error Unmarshaling expected Container definitions: %v", err.Error())
	}

	// Check if Name key is correct
	assert.Equal(t, expectedContainerDefinitions[0].Name, patchedContainerDefinitions[0].Name)
	assert.Equal(t, expectedContainerDefinitions[0].Name2, "")

	// The order received from patchedOutput changes continuously hence it is important to check if the arrays of expected and actual are equal without order being correct. This check also
	// helps with checking if key/value is named "Name" and "Value" accordingly.
	assert.ElementsMatch(t, expectedContainerDefinitions[0].Environment, patchedContainerDefinitions[0].Environment)

	// Check if Image key is correct
	assert.Equal(t, expectedContainerDefinitions[0].Image, patchedContainerDefinitions[0].Image)
	assert.Equal(t, patchedContainerDefinitions[0].Image2, "")

	// Check if entryPoint key is correct
	assert.Equal(t, expectedContainerDefinitions[0].EntryPoint, patchedContainerDefinitions[0].EntryPoint)
	assert.Equal(t, patchedContainerDefinitions[0].EntryPoint2, "")
}

func TestPatchFargateTaskDefinition(t *testing.T) {
	// Kilt Configuration, test invariant
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

	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		ImageAuthSecret:    "image_auth_secret",
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonRecipeConfig),
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
		testName         string
		logConfig        map[string]interface{}
		ignoreContainers []string
	}{
		{
			testName:         `fargate_entrypoint_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_env_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_cmd_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_linuxparameters_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_combined_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_volumesfrom_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_field_case_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{},
		},
		{
			testName: `fargate_log_group`,
			logConfig: map[string]interface{}{
				"group":         "test_log_group",
				"stream_prefix": "test_prefix",
				"region":        "test_region",
			},
			ignoreContainers: []string{},
		},
		{
			testName:         `fargate_ignore_container_test`,
			logConfig:        map[string]interface{}{},
			ignoreContainers: []string{"other", "another"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			patched, err := patchFargateTaskDefinition(
				context.Background(),
				getContainerDefinitionOriginal(tc.testName),
				kiltConfig,
				tc.logConfig,
				&tc.ignoreContainers)
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("Could not patch task definition, got error: %v", err))
			}
			sortAndCompare(t, []byte(getContainerDefinitionPatched(tc.testName)), []byte(*patched))
		})
	}
}

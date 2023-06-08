//go:build unit

package sysdig

import (
	"context"
	"encoding/json"
	"os"
	"sort"
	"testing"

	"github.com/Jeffail/gabs/v2"
	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/stretchr/testify/assert"
)

var (
	testKiltDefinition = KiltRecipeConfig{
		SysdigAccessKey:  "sysdig_access_key",
		AgentImage:       "workload_agent_image",
		OrchestratorHost: "orchestrator_host",
		OrchestratorPort: "orchestrator_port",
		CollectorHost:    "collector_host",
		CollectorPort:    "collector_port",
		SysdigLogging:    "sysdig_logging",
	}

	testIgnoreContainers = []string{}

	testContainerDefinitionFiles = []string{
		"fargate_entrypoint_test",
		"fargate_env_test",
		"fargate_cmd_test",
		"fargate_linuxparameters_test",
		"fargate_combined_test",
		"fargate_volumesfrom_test",
		"fargate_field_case_test",
	}
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

	patchedOutput, err := patchFargateTaskDefinition(context.Background(), string(inputfile), kiltConfig, nil, &testIgnoreContainers)
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

func TestTransform(t *testing.T) {
	for _, testName := range testContainerDefinitionFiles {
		t.Run(testName, func(t *testing.T) {
			jsonConfig, _ := json.Marshal(testKiltDefinition)
			kiltConfig := &cfnpatcher.Configuration{
				Kilt:               agentinoKiltDefinition,
				ImageAuthSecret:    "image_auth_secret",
				OptIn:              false,
				UseRepositoryHints: true,
				RecipeConfig:       string(jsonConfig),
			}

			inputContainerDefinition, _ := os.ReadFile("testfiles/" + testName + ".json")
			patched, _ := patchFargateTaskDefinition(context.Background(), string(inputContainerDefinition), kiltConfig, nil, &testIgnoreContainers)
			expectedContainerDefinition, _ := os.ReadFile("testfiles/" + testName + "_expected.json")

			sortAndCompare(t, expectedContainerDefinition, []byte(*patched))
		})
	}
}

func TestLogGroup(t *testing.T) {
	jsonConfig, _ := json.Marshal(testKiltDefinition)
	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		ImageAuthSecret:    "image_auth_secret",
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConfig),
	}

	logConfig := map[string]interface{}{
		"group":         "test_log_group",
		"stream_prefix": "test_prefix",
		"region":        "test_region",
	}

	inputContainerDefinition, _ := os.ReadFile("testfiles/fargate_log_group.json")
	patched, _ := patchFargateTaskDefinition(context.Background(), string(inputContainerDefinition), kiltConfig, logConfig, &testIgnoreContainers)
	expectedContainerDefinition, _ := os.ReadFile("testfiles/fargate_log_group_expected.json")

	sortAndCompare(t, expectedContainerDefinition, []byte(*patched))
}

func TestIgnoreContainers(t *testing.T) {
	jsonConfig, _ := json.Marshal(testKiltDefinition)
	kiltConfig := &cfnpatcher.Configuration{
		Kilt:               agentinoKiltDefinition,
		ImageAuthSecret:    "image_auth_secret",
		OptIn:              false,
		UseRepositoryHints: true,
		RecipeConfig:       string(jsonConfig),
	}

	fileTemplate := "fargate_ignore_container_test"
	ignoreContainers := []string{"other", "another"}

	inputContainerDefinition, _ := os.ReadFile("testfiles/" + fileTemplate + ".json")
	patched, _ := patchFargateTaskDefinition(context.Background(), string(inputContainerDefinition), kiltConfig, nil, &ignoreContainers)
	expectedContainerDefinition, _ := os.ReadFile("testfiles/" + fileTemplate + "_expected.json")

	sortAndCompare(t, expectedContainerDefinition, []byte(*patched))
}

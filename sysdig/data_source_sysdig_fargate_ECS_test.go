package sysdig

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/falcosecurity/kilt/runtimes/cloudformation/cfnpatcher"
	"github.com/stretchr/testify/assert"
)

func TestECStransformation(t *testing.T) {
	inputfile, err := ioutil.ReadFile("testfiles/ECSinput.json")

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

	patchedOutput, err := patchFargateTaskDefinition(context.Background(), string(inputfile), kiltConfig)
	if err != nil {
		t.Fatalf("Cannot execute PatchFargateTaskDefinition : %v", err.Error())
	}

	expectedOutput, err := ioutil.ReadFile("testfiles/ECSInstrumented.json")
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

	assert.Equal(t, expectedContainerDefinitions[0].Name, patchedContainerDefinitions[0].Name)

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

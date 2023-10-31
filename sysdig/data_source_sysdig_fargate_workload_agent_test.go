//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure

package sysdig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

// Direct connection mode has been deprecated in Prod envs
func TestAccSysdigFargateWorkloadAgentDirectConnection(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getFargateWorkloadAgentDirectConnection(),
			},
		},
	})
}

func getFargateWorkloadAgentDirectConnection() string {
	return `
data "sysdig_fargate_workload_agent" "test" {
	container_definitions = "[]"

	image_auth_secret = ""
	collector_host = "foo"
	collector_port = 1234
	sysdig_access_key = "abcdef"
	workload_agent_image = "busybox"
	sysdig_logging = "info"
}
`
}

func TestAccSysdigFargateWorkloadAgentOrchestrated(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: getFargateWorkloadAgentOrchestrated(),
			},
		},
	})
}

func getFargateWorkloadAgentOrchestrated() string {
	return `
data "sysdig_fargate_workload_agent" "test" {
	container_definitions = "[]"

	orchestrator_host = "sysdig.orchestrator.agent.com"
	orchestrator_port = 6667
	workload_agent_image = "quay.io/sysdig/workload-agent:latest"
	sysdig_logging = "info"
}
`
}

//go:build tf_acc_sysdig_secure

package sysdig_test

import (
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAcceptSecurePostureRisk(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: acceptPostureRiskResource(),
			},
			{
				Config: acceptPostureRiskZone(),
			},
		},
	})
}

func acceptPostureRiskResource() string {
	return `
resource "sysdig_secure_posture_accept_risk" "accept_resource" {
    description = "test accept posture risk resource"
    control_name = "ServiceAccounts with cluster access"
    reason = "Risk Transferred"
    expires_in = "30 Days"
    filter = "name in ('system:controller:daemon-set-s') and kind in ('ClusterRole')"
}`
}

func acceptPostureRiskZone() string {
	return `
resource "sysdig_secure_posture_accept_risk" "accept_resource" {
    description = "test accept posture risk resource"
    control_name = "ServiceAccounts with cluster access"
    reason = "Risk Transferred"
    expires_in = "30 Days"
    filter = "name in ('system:controller:daemon-set-s') and kind in ('ClusterRole')"
}`
}

package sysdig_test

import (
	"fmt"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCreateCustomControlResource(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: preCheckAnyEnv(t, SysdigSecureApiTokenEnv),
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: createControlResource(rText()),
			},
		},
	})
}

func createControlResource(name string) string {
	return fmt.Sprintf(`resource "sysdig_secure_posture_control" "test" {
        name = "S3 - Enabled Versioning-test-%s"
        description = "S3 - Enabled Versioning"
        resource_kind = "AWS_S3_BUCKET"
        severity = "Low"
        rego          = <<-EOF

            package sysdig

            import future.keywords.if
            import future.keywords.in

            default risky := false

            risky if {
              count(input.Versioning) == 0
            }

            risky if {
              some version in input.Versioning
              lower(version.Status) != "enabled"
            }
        EOF
     
     remediation_details = <<-EOF 
      **Using AWS CLI**\n1. 
    EOF
	}`, name)
}

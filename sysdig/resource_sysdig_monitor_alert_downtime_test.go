package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccAlertDowntime(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: alertDowntimeWithName(rText()),
			},
		},
	})
}

func alertDowntimeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_downtime" "sample" {
	name = "TERRAFORM TEST - DOWNTIME %s"
	description = "TERRAFORM TEST - DOWNTIME %s"
	severity = 2

	entities_to_monitor = ["host.hostName", "host.mac"]
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10
	trigger_after_pct = 100

	enabled = false

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

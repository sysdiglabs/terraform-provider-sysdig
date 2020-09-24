package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestAccAlertGroupOutlier(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: alertGroupOutlierWithName(rText()),
			},
		},
	})
}

func alertGroupOutlierWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_group_outlier" "sample" {
	name = "TERRAFORM TEST - GROUP OUTLIER %s"
	description = "TERRAFORM TEST - GROUP OUTLIER %s"
	severity = 6

	monitor = ["cpu.cores.used", "cpu.cores.used.percent", "cpu.stolen.percent", "cpu.used.percent"]

	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

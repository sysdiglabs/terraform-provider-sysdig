//go:build tf_acc_sysdig || tf_acc_sysdig_monitor || tf_acc_ibm || tf_acc_ibm_monitor

package sysdig_test

import "fmt"

func monitorTeamMinimumConfiguration(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name      = "sample-%s"

  entrypoint {
	type = "Explore"
  }
}`, name)
}

func monitorTeamWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_team" "sample" {
  name               = "sample-%s"
  description        = "%s"
  scope_by           = "container"
  filter             = "container.image.repo = \"sysdig/agent\""

  entrypoint {
	type = "Explore"
  }
}`, name, name)
}

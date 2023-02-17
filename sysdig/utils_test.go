package sysdig_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"testing"
)

func randomText() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
}

func preCheckMonitorToken(t *testing.T) func() {
	return func() {
		if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
			t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
		}
	}
}

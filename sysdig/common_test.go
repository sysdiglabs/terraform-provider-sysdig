//go:build tf_acc_sysdig || tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_ibm || tf_acc_ibm_monitor || tf_acc_policies || unit

package sysdig_test

import (
	"os"
	"strings"
	"testing"
)

const (
	SysdigMonitorApiTokenEnv  = "SYSDIG_MONITOR_API_TOKEN"
	SysdigSecureApiTokenEnv   = "SYSDIG_SECURE_API_TOKEN"
	SysdigIBMMonitorAPIKeyEnv = "SYSDIG_IBM_MONITOR_API_KEY"
	SysdigIBMSecureAPIKeyEnv  = "SYSDIG_IBM_SECURE_API_KEY"
)

func isAnyEnvSet(envs ...string) bool {
	for _, env := range envs {
		if value := os.Getenv(env); value != "" {
			return true
		}
	}
	return false
}

func preCheckAnyEnv(t *testing.T, envs ...string) func() {
	return func() {
		if !isAnyEnvSet(envs...) {
			t.Fatalf("%s must be set for acceptance tests", strings.Join(envs, " or "))
		}
	}
}

func sysdigOrIBMMonitorPreCheck(t *testing.T) func() {
	return func() {
		monitor := os.Getenv("SYSDIG_MONITOR_API_TOKEN")
		ibmMonitor := os.Getenv("SYSDIG_IBM_MONITOR_API_KEY")
		if monitor == "" && ibmMonitor == "" {
			t.Fatal("SYSDIG_MONITOR_API_TOKEN or SYSDIG_IBM_MONITOR_API_KEY must be set for acceptance tests")
		}
	}
}

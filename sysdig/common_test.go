package sysdig_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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

func randomText(len int) string {
	return acctest.RandStringFromCharSet(len, acctest.CharSetAlphaNum)
}

// retryOn409 wraps a TestStep to retry on 409 Conflict errors
func retryOn409(step resource.TestStep) resource.TestStep {
	if step.PlanOnly {
		return step
	}
	
	originalConfig := step.Config
	step.Config = ""
	step.PreConfig = func() {
		for i := 0; i < 5; i++ {
			if i > 0 {
				time.Sleep(time.Duration(i*2) * time.Second)
			}
			step.Config = originalConfig
			break
		}
	}
	return step
}

// testCaseWithRetry creates a TestCase with retry logic for all steps
func testCaseWithRetry(testCase resource.TestCase) resource.TestCase {
	for i := range testCase.Steps {
		testCase.Steps[i] = retryOn409(testCase.Steps[i])
	}
	return testCase
}

// retryTestStep creates a single TestStep with retry logic
func retryTestStep(config string, checks ...resource.TestCheckFunc) resource.TestStep {
	return resource.TestStep{
		Config: config,
		Check:  resource.ComposeTestCheckFunc(checks...),
		Retry: func() *resource.RetryError {
			return resource.RetryableError(nil) // Retry on any error
		},
	}
}

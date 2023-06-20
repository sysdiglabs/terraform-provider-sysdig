package sysdig_test

import (
	"errors"
	"github.com/draios/terraform-provider-sysdig/codeowner"
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

const (
	SysdigMonitorApiTokenEnv  = "SYSDIG_MONITOR_API_TOKEN"
	SysdigSecureApiTokenEnv   = "SYSDIG_SECURE_API_TOKEN"
	SysdigIBMMonitorAPIKeyEnv = "SYSDIG_IBM_MONITOR_API_KEY"
	SysdigIBMSecureAPIKeyEnv  = "SYSDIG_IBM_SECURE_API_KEY"
)

func findCaller(t *testing.T) (string, error) {
	// subtest can have / in the name,
	// we should handle that scenario
	name := strings.Split(t.Name(), "/")[0]
	r, err := regexp.Compile(name)
	if err != nil {
		t.Fatal(err)
	}

	// we are not going to handle more than 100 stack calls
	var pcs [100]uintptr
	n := runtime.Callers(0, pcs[:])
	iter := runtime.CallersFrames(pcs[:n])
	for {
		f, more := iter.Next()
		if r.MatchString(f.Func.Name()) {
			return f.File, nil
		}

		if !more {
			break
		}
	}

	return "", errors.New("failed to find function caller")
}

func handleReport(t *testing.T) {
	if t.Failed() {
		// get file path from which this function is called
		callerFile, err := findCaller(t)
		if err != nil {
			t.Fatalf("failed to recover caller information: %v", err)
		}

		owners, err := codeowner.LoadOwners(callerFile)
		if err != nil {
			t.Fatalf("failed to create report notification: %v", err)
		}

		t.Fatalf("report to %s", strings.Join(owners, ", "))
	}
}

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

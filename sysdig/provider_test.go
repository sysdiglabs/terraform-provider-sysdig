//go:build sysdig_monitor || sysdig_secure || ibm_monitor

package sysdig_test

import (
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestProvider(t *testing.T) {
	if err := sysdig.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

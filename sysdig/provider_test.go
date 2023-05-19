//go:build tf_acc_sysdig_monitor || tf_acc_sysdig_secure || tf_acc_ibm_monitor

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

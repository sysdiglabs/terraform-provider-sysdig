//go:build tf_acc_sysdig

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

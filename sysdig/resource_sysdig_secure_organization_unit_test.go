package sysdig

import (
	"context"
	"testing"

	"github.com/hashicorp/go-cty/cty"
)

// State written by released versions carries these collections as lists, so the V0 type has to
// keep that shape while the current schema reads them as sets.
func TestSecureOrganizationStateUpgradeFromV0(t *testing.T) {
	collections := []string{
		SchemaOrganizationalUnitIds,
		SchemaIncludedOrganizationalGroups,
		SchemaExcludedOrganizationalGroups,
		SchemaIncludedCloudAccounts,
		SchemaExcludedCloudAccounts,
	}

	res := resourceSysdigSecureOrganization()
	if res.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1", res.SchemaVersion)
	}
	if len(res.StateUpgraders) != 1 || res.StateUpgraders[0].Version != 0 {
		t.Fatalf("expected exactly one upgrader from version 0, got %+v", res.StateUpgraders)
	}

	v0 := resourceSysdigSecureOrganizationV0().CoreConfigSchema().ImpliedType()
	current := res.CoreConfigSchema().ImpliedType()
	for _, key := range collections {
		if got := v0.AttributeType(key); !got.Equals(cty.List(cty.String)) {
			t.Errorf("v0 %s = %s, want list of string", key, got.FriendlyName())
		}
		if got := current.AttributeType(key); !got.Equals(cty.Set(cty.String)) {
			t.Errorf("current %s = %s, want set of string", key, got.FriendlyName())
		}
	}

	// the collections carry no nested state, so the upgrade only re-reads them under the new type
	state := map[string]any{
		"id":                      "org-id",
		"excluded_cloud_accounts": []any{"second", "first"},
	}
	upgraded, err := res.StateUpgraders[0].Upgrade(context.Background(), state, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	accounts, ok := upgraded["excluded_cloud_accounts"].([]any)
	if !ok || len(accounts) != 2 || accounts[0] != "second" || accounts[1] != "first" {
		t.Errorf("upgrade altered the raw state: %+v", upgraded)
	}
}

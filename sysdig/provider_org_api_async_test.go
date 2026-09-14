//go:build unit

package sysdig_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

// Runs through the real diff layer: EnvDefaultFunc hands back the raw env string, and only
// the diff coerces it to TypeBool, so an explicit false could silently lose to the env var.
func TestProviderOrgAPIAsync(t *testing.T) {
	const key = "sysdig_secure_org_api_async"

	tests := []struct {
		name   string
		env    string
		config map[string]any
		want   bool
	}{
		{name: "unset everywhere", env: "", config: map[string]any{}, want: false},
		{name: "env var enables it", env: "true", config: map[string]any{}, want: true},
		{name: "explicit false beats env var", env: "true", config: map[string]any{key: false}, want: false},
		{name: "explicit true beats env var", env: "false", config: map[string]any{key: true}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SYSDIG_ORG_API_ASYNC", tt.env)

			sm := schema.InternalMap(sysdig.Provider().Schema)
			diff, err := sm.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(tt.config), nil, nil, true)
			if err != nil {
				t.Fatalf("Diff: %v", err)
			}
			d, err := sm.Data(nil, diff)
			if err != nil {
				t.Fatalf("Data: %v", err)
			}

			got, ok := d.Get(key).(bool)
			if !ok {
				t.Fatalf("%s is not a bool: %#v", key, d.Get(key))
			}
			if got != tt.want {
				t.Errorf("%s = %v, want %v", key, got, tt.want)
			}
		})
	}
}

//go:build unit

package sysdig_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

const orgAPIAsyncKey = "sysdig_secure_org_api_async"

func resolveOrgAPIAsync(t *testing.T, config map[string]any) any {
	t.Helper()

	sm := schema.InternalMap(sysdig.Provider().Schema)
	diff, err := sm.Diff(context.Background(), nil, terraform.NewResourceConfigRaw(config), nil, nil, true)
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	d, err := sm.Data(nil, diff)
	if err != nil {
		t.Fatalf("Data: %v", err)
	}
	return d.Get(orgAPIAsyncKey)
}

// Config has to beat the environment, and the diff layer only coerces to TypeBool after the
// default func runs, so an explicit false could otherwise lose to the variable.
func TestProviderOrgAPIAsync(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		config map[string]any
		want   bool
	}{
		{name: "unset everywhere", config: map[string]any{}, want: false},
		{name: "env var enables it", env: map[string]string{"SYSDIG_ORG_API_ASYNC": "true"}, config: map[string]any{}, want: true},
		{name: "prefixed env var enables it", env: map[string]string{"SYSDIG_SECURE_ORG_API_ASYNC": "true"}, config: map[string]any{}, want: true},
		{
			name:   "prefixed env var wins over the legacy one",
			env:    map[string]string{"SYSDIG_SECURE_ORG_API_ASYNC": "false", "SYSDIG_ORG_API_ASYNC": "true"},
			config: map[string]any{},
			want:   false,
		},
		{
			name:   "explicit false beats env var",
			env:    map[string]string{"SYSDIG_ORG_API_ASYNC": "true"},
			config: map[string]any{orgAPIAsyncKey: false},
			want:   false,
		},
		{name: "explicit true with no env var", config: map[string]any{orgAPIAsyncKey: true}, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SYSDIG_ORG_API_ASYNC", "")
			t.Setenv("SYSDIG_SECURE_ORG_API_ASYNC", "")
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			got, ok := resolveOrgAPIAsync(t, tt.config).(bool)
			if !ok {
				t.Fatalf("%s is not a bool: %#v", orgAPIAsyncKey, resolveOrgAPIAsync(t, tt.config))
			}
			if got != tt.want {
				t.Errorf("%s = %v, want %v", orgAPIAsyncKey, got, tt.want)
			}
		})
	}
}

// Only the exact string "true" has ever enabled this. Anything else must stay disabled rather
// than flip the flag or fail provider configuration on an unparseable value.
func TestProviderOrgAPIAsyncEnvValues(t *testing.T) {
	enabled := []string{"true"}
	disabled := []string{"", "false", "0", "1", "t", "TRUE", "True", "yes", "on", "enabled", "true "}

	for _, v := range enabled {
		t.Run("enabled/"+v, func(t *testing.T) {
			t.Setenv("SYSDIG_ORG_API_ASYNC", v)
			if got := resolveOrgAPIAsync(t, map[string]any{}); got != true {
				t.Errorf("env %q resolved to %#v, want true", v, got)
			}
		})
	}
	for _, v := range disabled {
		t.Run("disabled/"+v, func(t *testing.T) {
			t.Setenv("SYSDIG_ORG_API_ASYNC", v)
			if got := resolveOrgAPIAsync(t, map[string]any{}); got != false {
				t.Errorf("env %q resolved to %#v, want false", v, got)
			}
		})
	}
}

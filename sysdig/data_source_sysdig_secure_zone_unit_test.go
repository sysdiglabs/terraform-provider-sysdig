package sysdig

import (
	"testing"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
)

func TestGetZoneScopes_MatchByID(t *testing.T) {
	zoneV2 := &v2.ZoneV2{
		Scopes: []v2.ScopeV2{
			{
				Filters: []v2.FilterV2{
					{
						ID:           10,
						ResourceType: "kubernetes",
						Expressions: []v2.ExpressionV2{
							{Field: "cluster", Operator: "in", Values: []string{"a"}},
						},
					},
					{
						ID:           20,
						ResourceType: "kubernetes",
						Expressions: []v2.ExpressionV2{
							{Field: "cluster", Operator: "in", Values: []string{"b"}},
						},
					},
				},
			},
		},
	}

	result := getZoneScopes(zoneV2)
	if len(result) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(result))
	}

	// Verify each scope got its own expressions by checking the ID-expression pairing.
	for _, raw := range result {
		scope := raw.(map[string]any)
		id := scope[SchemaIDKey].(int)
		exprs, ok := scope[SchemaExpressionKey].([]any)
		if !ok || len(exprs) != 1 {
			t.Fatalf("scope ID=%d: expected 1 expression, got %v", id, scope[SchemaExpressionKey])
		}
		expr := exprs[0].(map[string]any)
		vals := expr["values"].([]string)
		if id == 10 && vals[0] != "a" {
			t.Errorf("scope ID=10: expected values [a], got %v", vals)
		}
		if id == 20 && vals[0] != "b" {
			t.Errorf("scope ID=20: expected values [b], got %v", vals)
		}
	}
}

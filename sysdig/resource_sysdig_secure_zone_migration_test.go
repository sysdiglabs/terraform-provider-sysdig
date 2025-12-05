package sysdig

import (
	"testing"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestExpandFlattenExpression_SingleValue(t *testing.T) {
	raw := map[string]interface{}{
		"field":    "agent.tag.key",
		"operator": "is_not",
		"value":    "test",
	}

	exp := expandExpressionV2(raw)
	flat := flattenExpressionV2(exp)

	require.Equal(t, "agent.tag.key", flat["field"])
	require.Equal(t, "is_not", flat["operator"])
	require.Equal(t, "test", flat["value"])
}

func TestExpandFlattenExpression_MultipleValues(t *testing.T) {
	raw := map[string]interface{}{
		"field":    "organization",
		"operator": "in",
		"values":   []interface{}{"o1", "o2"},
	}

	exp := expandExpressionV2(raw)
	flat := flattenExpressionV2(exp)

	require.Equal(t, "organization", flat["field"])
	require.Equal(t, "in", flat["operator"])

	values, ok := flat["values"].([]string)
	require.True(t, ok)
	require.ElementsMatch(t, []string{"o1", "o2"}, values)
}

func TestExpandExpression_ValuesWinOverValue(t *testing.T) {
	raw := map[string]interface{}{
		"field":    "organization",
		"operator": "in",
		"value":    "SHOULD_NOT_BE_USED",
		"values":   []interface{}{"o1", "o2"},
	}

	exp := expandExpressionV2(raw)

	require.Equal(t, []string{"o1", "o2"}, exp.Values)
	require.Empty(t, exp.Value)
}

func TestExpandFlattenFilter_MultipleExpressions(t *testing.T) {
	raw := map[string]interface{}{
		"id":          0,
		"target_type": "kubernetes",
		"expression": []interface{}{
			map[string]interface{}{
				"field":    "agent.tag.key",
				"operator": "is_not",
				"value":    "test",
			},
			map[string]interface{}{
				"field":    "agent.tag.key2",
				"operator": "not_contains",
				"value":    "value2",
			},
		},
	}

	filter := expandFilterV2(raw)
	flat := flattenFilterV21(filter, false)

	require.Equal(t, "kubernetes", flat["target_type"])
	require.Len(t, flat["expression"], 2)
}

func TestFlattenZoneV21_MultipleScopesAndFilters(t *testing.T) {
	zone := &v2.ZoneV2{
		Scopes: []v2.ScopeV2{
			{
				Filters: []v2.FilterV2{
					{ResourceType: "kubernetes", Expressions: []v2.ExpressionV2{
						{Field: "agent.tag.env", Operator: "in", Values: []string{"prod"}},
					}},
				},
			},
		},
	}
	d := schema.TestResourceDataRaw(t, resourceSysdigSecureZone().Schema, nil)

	scopes := flattenZoneV2(zone, false)
	err := d.Set(SchemaScopeKey, scopes)
	require.NoError(t, err)

	if rawScopes, ok := d.Get(SchemaScopeKey).(*schema.Set); ok {
		scopes := rawScopes.List()
		require.Len(t, scopes, 1)
		sc := scopes[0].(map[string]interface{})
		require.Equal(t, "kubernetes", sc["target_type"])

		exprs := sc["expression"].([]interface{})
		require.Len(t, exprs, 1)

	} else {
		t.Fatalf("expected *schema.Set for scopes, got %T", d.Get(SchemaScopeKey))
	}
}

func TestExpandFlattenZoneV21_RoundTrip(t *testing.T) {
	hclInput := map[string]interface{}{
		"name":        "example-zone-legacy",
		"description": "Migrated to expressions",
		"scope": []interface{}{
			map[string]interface{}{
				"target_type": "kubernetes",
				"expression": []interface{}{
					map[string]interface{}{
						"field":    "agent.tag.key",
						"operator": "is_not",
						"value":    "test",
					},
					map[string]interface{}{
						"field":    "agent.tag.key2",
						"operator": "not_contains",
						"value":    "value2",
					},
				},
			},
		},
	}

	d1 := schema.TestResourceDataRaw(t, resourceSysdigSecureZone().Schema, hclInput)

	zone := expandZoneV2(d1)

	d2 := schema.TestResourceDataRaw(t, resourceSysdigSecureZone().Schema, nil)
	scopes := flattenZoneV2(zone, false)
	err := d2.Set(SchemaScopeKey, scopes)
	require.NoError(t, err)

	require.ElementsMatch(t, d1.Get("scope").(*schema.Set).List(), d2.Get("scope").(*schema.Set).List())
}

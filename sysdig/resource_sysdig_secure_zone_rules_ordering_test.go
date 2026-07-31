//go:build unit

package sysdig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

// The value order inside an `in (...)` list carries no meaning. These tests pin
// down that the provider treats it that way without ever equating two rules that
// actually differ.

func TestCanonicalizeZoneRules(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "sorts values",
			in:   `clusterId in ("c", "a", "b")`,
			want: `clusterId in ("a", "b", "c")`,
		},
		{
			name: "already sorted is unchanged",
			in:   `clusterId in ("a", "b", "c")`,
			want: `clusterId in ("a", "b", "c")`,
		},
		{
			name: "normalizes missing spaces after commas",
			in:   `clusterId in ("b","a")`,
			want: `clusterId in ("a", "b")`,
		},
		{
			name: "handles every list of a compound expression",
			in:   `clusterId in ("b", "a") and label.team in ("z", "y")`,
			want: `clusterId in ("a", "b") and label.team in ("y", "z")`,
		},
		{
			name: "handles not in",
			in:   `clusterId not in ("b", "a")`,
			want: `clusterId not in ("a", "b")`,
		},
		{
			name: "operator casing is preserved",
			in:   `clusterId IN ("b", "a")`,
			want: `clusterId IN ("a", "b")`,
		},
		{
			name: "keeps a comma that belongs to a value",
			in:   `agentTags in ("project: x", "project: foo, bar")`,
			want: `agentTags in ("project: foo, bar", "project: x")`,
		},
		{
			name: "handles escaped quotes inside a value",
			in:   `name in ("h\"2", "h1")`,
			want: `name in ("h1", "h\"2")`,
		},
		{
			name: "single value",
			in:   `clusterId in ("a")`,
			want: `clusterId in ("a")`,
		},
		{
			name: "empty list",
			in:   `clusterId in ()`,
			want: `clusterId in ()`,
		},
		{
			name: "leaves contains untouched, even with a comma",
			in:   `name contains "foo, bar"`,
			want: `name contains "foo, bar"`,
		},
		{
			name: "leaves exists untouched",
			in:   `label.env exists`,
			want: `label.env exists`,
		},
		{
			name: "empty rules",
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := canonicalizeZoneRules(tt.in)
			require.Equal(t, tt.want, got)
			require.Equal(t, got, canonicalizeZoneRules(got), "canonicalization must be idempotent")
		})
	}
}

func TestCanonicalizeZoneRulesNeverEquatesDifferentRules(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
	}{
		{
			name: "different value",
			a:    `clusterId in ("a", "b")`,
			b:    `clusterId in ("a", "c")`,
		},
		{
			name: "extra value",
			a:    `clusterId in ("a", "b")`,
			b:    `clusterId in ("a", "b", "c")`,
		},
		{
			name: "duplicate instead of distinct values",
			a:    `clusterId in ("a", "b")`,
			b:    `clusterId in ("a", "a")`,
		},
		{
			name: "one value with a comma vs two values",
			a:    `agentTags in ("project: foo, bar")`,
			b:    `agentTags in ("project: foo", "bar")`,
		},
		{
			name: "different field",
			a:    `clusterId in ("a", "b")`,
			b:    `name in ("a", "b")`,
		},
		{
			name: "different operator",
			a:    `clusterId in ("a", "b")`,
			b:    `clusterId not in ("a", "b")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotEqual(t, canonicalizeZoneRules(tt.a), canonicalizeZoneRules(tt.b))
		})
	}
}

func TestSuppressZoneRulesValueOrder(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "reordering is suppressed",
			old:  `clusterId in ("b", "a")`,
			new:  `clusterId in ("a", "b")`,
			want: true,
		},
		{
			name: "a real change is not suppressed",
			old:  `clusterId in ("a", "b")`,
			new:  `clusterId in ("a", "c")`,
			want: false,
		},
		{
			// An expression-based scope carries an empty rules string. Suppressing
			// it would drop the attribute from the planned scope block.
			name: "empty on both sides is not suppressed",
			old:  "",
			new:  "",
			want: false,
		},
		{
			name: "rules being set is not suppressed",
			old:  "",
			new:  `clusterId in ("a")`,
			want: false,
		},
		{
			name: "rules being cleared is not suppressed",
			old:  `clusterId in ("a")`,
			new:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, suppressZoneRulesValueOrder("scope.0.rules", tt.old, tt.new, nil))
		})
	}
}

func rulesScope(id int, targetType, rules string) map[string]any {
	return map[string]any{
		SchemaIDKey:         id,
		SchemaTargetTypeKey: targetType,
		SchemaRulesKey:      rules,
		SchemaExpressionKey: []any{},
	}
}

func TestZoneScopeSetHashIgnoresRulesValueOrder(t *testing.T) {
	base := rulesScope(0, "kubernetes", `clusterId in ("b", "a", "c")`)
	reordered := rulesScope(0, "kubernetes", `clusterId in ("a", "c", "b")`)

	require.Equal(t, zoneScopeSetHash(base), zoneScopeSetHash(reordered),
		"reordering the values must not change the identity of a scope block")
}

func TestZoneScopeSetHashDistinguishesRealDifferences(t *testing.T) {
	base := rulesScope(0, "kubernetes", `clusterId in ("a", "b")`)

	tests := map[string]map[string]any{
		"different value":       rulesScope(0, "kubernetes", `clusterId in ("a", "c")`),
		"extra value":           rulesScope(0, "kubernetes", `clusterId in ("a", "b", "c")`),
		"different target type": rulesScope(0, "host", `clusterId in ("a", "b")`),
		"different field":       rulesScope(0, "kubernetes", `name in ("a", "b")`),
	}

	for name, other := range tests {
		t.Run(name, func(t *testing.T) {
			require.NotEqual(t, zoneScopeSetHash(base), zoneScopeSetHash(other))
		})
	}
}

func TestZoneScopeSetHashIgnoresComputedID(t *testing.T) {
	// The state holds a server-assigned id, the configuration does not. If the id
	// took part in the hash, a configured block could never match its own state.
	withID := rulesScope(5179, "kubernetes", `clusterId in ("a", "b")`)
	withoutID := rulesScope(0, "kubernetes", `clusterId in ("a", "b")`)

	require.Equal(t, zoneScopeSetHash(withID), zoneScopeSetHash(withoutID))
}

func TestZoneScopeSetHashKeepsExpressionScopesDistinct(t *testing.T) {
	// Expression scopes fall back to the stock hash on purpose: two blocks that
	// share a field and an operator must not collapse into one.
	expressionScope := func(values ...string) map[string]any {
		vals := make([]any, len(values))
		for i, v := range values {
			vals[i] = v
		}
		return map[string]any{
			SchemaIDKey:         0,
			SchemaTargetTypeKey: "kubernetes",
			SchemaRulesKey:      "",
			SchemaExpressionKey: []any{
				map[string]any{
					SchemaFieldKey:    "kubernetes.cluster.name",
					SchemaOperatorKey: "in",
					SchemaValueKey:    "",
					SchemaValuesKey:   vals,
				},
			},
		}
	}

	require.NotEqual(t,
		zoneScopeSetHash(expressionScope("a1", "a2")),
		zoneScopeSetHash(expressionScope("b1", "b2")),
		"distinct expression scopes must not collide")
}

// zoneDiff computes the diff the SDK would produce for a zone whose state holds
// stateScopes and whose configuration holds configScopes.
func zoneDiff(t *testing.T, stateScopes, configScopes []any) *terraform.InstanceDiff {
	t.Helper()

	r := resourceSysdigSecureZone()

	d := r.Data(&terraform.InstanceState{ID: "314"})
	require.NoError(t, d.Set(SchemaNameKey, "zone"))
	require.NoError(t, d.Set(SchemaScopeKey, stateScopes))
	state := d.State()
	state.ID = "314"

	config := terraform.NewResourceConfigRaw(map[string]any{
		SchemaNameKey:  "zone",
		SchemaScopeKey: configScopes,
	})

	diff, err := r.Diff(context.Background(), state, config, nil)
	require.NoError(t, err)
	return diff
}

func configScope(targetType, rules string) map[string]any {
	return map[string]any{
		SchemaTargetTypeKey: targetType,
		SchemaRulesKey:      rules,
	}
}

func TestZoneDiffIgnoresRulesValueReordering(t *testing.T) {
	diff := zoneDiff(t,
		[]any{rulesScope(5179, "kubernetes", `clusterId in ("c2", "c1", "c3")`)},
		[]any{configScope("kubernetes", `clusterId in ("c3", "c2", "c1")`)},
	)

	require.True(t, diff == nil || diff.Empty(), "reordering the values must not produce a diff, got %#v", diff)
}

func TestZoneDiffIgnoresReorderingInOneOfSeveralScopes(t *testing.T) {
	diff := zoneDiff(t,
		[]any{
			rulesScope(5179, "kubernetes", `clusterId in ("c2", "c1")`),
			rulesScope(5180, "host", `name in ("h1", "h2")`),
		},
		[]any{
			configScope("kubernetes", `clusterId in ("c1", "c2")`),
			configScope("host", `name in ("h1", "h2")`),
		},
	)

	require.True(t, diff == nil || diff.Empty(), "got %#v", diff)
}

func TestZoneDiffStillDetectsRealChanges(t *testing.T) {
	tests := map[string]string{
		"replaced value": `clusterId in ("c1", "brand-new")`,
		"added value":    `clusterId in ("c1", "c2", "c3")`,
		"removed value":  `clusterId in ("c1")`,
	}

	for name, configRules := range tests {
		t.Run(name, func(t *testing.T) {
			diff := zoneDiff(t,
				[]any{rulesScope(5179, "kubernetes", `clusterId in ("c1", "c2")`)},
				[]any{configScope("kubernetes", configRules)},
			)

			require.NotNil(t, diff)
			require.False(t, diff.Empty(), "a real change must still be planned")
		})
	}
}

func TestZoneDiffStillDetectsTargetTypeChange(t *testing.T) {
	diff := zoneDiff(t,
		[]any{rulesScope(5179, "kubernetes", `clusterId in ("c1", "c2")`)},
		[]any{configScope("host", `clusterId in ("c1", "c2")`)},
	)

	require.NotNil(t, diff)
	require.False(t, diff.Empty())
}

func TestZoneScopeSchemaUsesTheOrderInsensitiveHash(t *testing.T) {
	// Guards against the schema being rewritten without the custom hash, which
	// would silently bring the spurious diffs back.
	scope := resourceSysdigSecureZone().Schema[SchemaScopeKey]
	require.NotNil(t, scope.Set, "the scope set must use zoneScopeSetHash")

	elem, ok := scope.Elem.(*schema.Resource)
	require.True(t, ok)
	require.NotNil(t, elem.Schema[SchemaRulesKey].DiffSuppressFunc,
		"rules must keep its DiffSuppressFunc")
}

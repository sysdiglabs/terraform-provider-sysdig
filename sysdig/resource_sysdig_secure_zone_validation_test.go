package sysdig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveIdentifier(t *testing.T) {
	tests := []struct {
		name       string
		targetType string
		field      string
		wantIdent  ruleIdentifier
		wantOK     bool
	}{
		// label.<key> disambiguation
		{
			name:       "label on aws resolves to labels",
			targetType: "aws",
			field:      "label.team",
			wantIdent:  identLabels,
			wantOK:     true,
		},
		{
			name:       "label on kubernetes resolves to labelValues",
			targetType: "kubernetes",
			field:      "label.team",
			wantIdent:  identLabelValues,
			wantOK:     true,
		},
		{
			name:       "label on host resolves to labelValues",
			targetType: "host",
			field:      "label.env",
			wantIdent:  identLabelValues,
			wantOK:     true,
		},
		{
			name:       "label on ibm resolves to labels",
			targetType: "ibm",
			field:      "label.cost-center",
			wantIdent:  identLabels,
			wantOK:     true,
		},
		{
			name:       "label on oci resolves to labels",
			targetType: "oci",
			field:      "label.env",
			wantIdent:  identLabels,
			wantOK:     true,
		},
		// agent.tag.<key>
		{
			name:       "agent.tag resolves to agentTags",
			targetType: "kubernetes",
			field:      "agent.tag.cluster",
			wantIdent:  identAgentTags,
			wantOK:     true,
		},
		// Direct fields
		{
			name:       "account direct match",
			targetType: "aws",
			field:      "account",
			wantIdent:  identAccount,
			wantOK:     true,
		},
		{
			name:       "clusterId direct match",
			targetType: "kubernetes",
			field:      "clusterId",
			wantIdent:  identClusterId,
			wantOK:     true,
		},
		{
			name:       "namespace direct match",
			targetType: "kubernetes",
			field:      "namespace",
			wantIdent:  identNamespace,
			wantOK:     true,
		},
		{
			name:       "gitIntegrationId direct match",
			targetType: "git",
			field:      "gitIntegrationId",
			wantIdent:  identGitIntegrationId,
			wantOK:     true,
		},
		// Unrecognized fields
		{
			name:       "completely unknown field",
			targetType: "aws",
			field:      "foobar",
			wantIdent:  "",
			wantOK:     false,
		},
		{
			name:       "label. with no key",
			targetType: "aws",
			field:      "label.",
			wantIdent:  "",
			wantOK:     false,
		},
		{
			name:       "agent.tag. with no key",
			targetType: "kubernetes",
			field:      "agent.tag.",
			wantIdent:  "",
			wantOK:     false,
		},
		{
			name:       "agent.tag without trailing dot",
			targetType: "kubernetes",
			field:      "agent.tag",
			wantIdent:  "",
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ident, ok := resolveIdentifier(tt.targetType, tt.field)
			require.Equal(t, tt.wantOK, ok)
			require.Equal(t, tt.wantIdent, ident)
		})
	}
}

func TestValidateExpressionFields(t *testing.T) {
	tests := []struct {
		name       string
		targetType string
		exprs      []interface{}
		wantErr    string // empty = no error
	}{
		// --- Recognized but not allowed → ERROR ---
		{
			name:       "kubernetes + account => error",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "account", "operator": "in"},
			},
			wantErr: `field "account" is not allowed for target_type "kubernetes"`,
		},
		{
			name:       "aws + namespace => error",
			targetType: "aws",
			exprs: []interface{}{
				map[string]interface{}{"field": "namespace", "operator": "in"},
			},
			wantErr: `field "namespace" is not allowed for target_type "aws"`,
		},
		{
			name:       "aws + agent.tag.env => error",
			targetType: "aws",
			exprs: []interface{}{
				map[string]interface{}{"field": "agent.tag.env", "operator": "in"},
			},
			wantErr: `field "agent.tag.env" is not allowed for target_type "aws"`,
		},
		{
			name:       "host + label.team => error (host doesn't support label fields)",
			targetType: "host",
			exprs: []interface{}{
				map[string]interface{}{"field": "label.team", "operator": "in"},
			},
			wantErr: `field "label.team" is not allowed for target_type "host"`,
		},
		{
			name:       "multiple expressions, second invalid",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "clusterId", "operator": "in"},
				map[string]interface{}{"field": "account", "operator": "in"},
			},
			wantErr: `scope[0].expression[1]: field "account" is not allowed`,
		},

		// --- Recognized and allowed → OK ---
		{
			name:       "kubernetes + label.team => OK",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "label.team", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "aws + label.team => OK",
			targetType: "aws",
			exprs: []interface{}{
				map[string]interface{}{"field": "label.team", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "kubernetes + agent.tag.env => OK",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "agent.tag.env", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "image + registry => OK",
			targetType: "image",
			exprs: []interface{}{
				map[string]interface{}{"field": "registry", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "ibm + resourceGroupId => OK",
			targetType: "ibm",
			exprs: []interface{}{
				map[string]interface{}{"field": "resourceGroupId", "operator": "in"},
			},
			wantErr: "",
		},

		// --- Unrecognized field → silently allowed (forward compat) ---
		{
			name:       "kubernetes + unknown field => no error",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "someUnknownNewField", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "aws + future field => no error",
			targetType: "aws",
			exprs: []interface{}{
				map[string]interface{}{"field": "compartmentId", "operator": "in"},
			},
			wantErr: "",
		},
		{
			name:       "mixed: known-invalid + unknown => error on known-invalid",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "futureField", "operator": "in"},
				map[string]interface{}{"field": "account", "operator": "in"},
			},
			wantErr: `field "account" is not allowed for target_type "kubernetes"`,
		},
		{
			name:       "mixed: unknown + known-valid => OK",
			targetType: "kubernetes",
			exprs: []interface{}{
				map[string]interface{}{"field": "futureField", "operator": "in"},
				map[string]interface{}{"field": "clusterId", "operator": "in"},
			},
			wantErr: "",
		},

		// --- Unknown target_type → ERROR (fixed set) ---
		{
			name:       "unknown target_type => error",
			targetType: "lambda",
			exprs: []interface{}{
				map[string]interface{}{"field": "account", "operator": "in"},
			},
			wantErr: `unknown target_type "lambda"`,
		},

		// --- Edge cases ---
		{
			name:       "empty expressions => OK",
			targetType: "aws",
			exprs:      []interface{}{},
			wantErr:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExpressionFields(tt.targetType, tt.exprs, 0)
			if tt.wantErr == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}

func TestAllowedFieldFamilies(t *testing.T) {
	k8s := allowedFieldFamilies("kubernetes")
	require.Contains(t, k8s, "label.<key>")
	require.Contains(t, k8s, "agent.tag.<key>")
	require.Contains(t, k8s, "clusterId")
	require.Contains(t, k8s, "namespace")
	require.NotContains(t, k8s, "account")

	aws := allowedFieldFamilies("aws")
	require.Contains(t, aws, "label.<key>")
	require.NotContains(t, aws, "agent.tag.<key>")
	require.Contains(t, aws, "account")

	require.Nil(t, allowedFieldFamilies("unknown"))
}

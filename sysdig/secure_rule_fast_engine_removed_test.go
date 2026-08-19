package sysdig

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// emptyResourceData is a throwaway *schema.ResourceData for handlers that only
// ever clear the ID.
func emptyResourceData(t *testing.T) *schema.ResourceData {
	t.Helper()

	return schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]any{})
}

// Every handler is called with a nil meta. Real CRUD in this provider starts
// with meta.(SysdigClients), which panics on nil, so surviving this is what
// pins the promise that these paths never reach for a client -- and therefore
// never hit the network. It is the property most likely to be broken by a
// future refactor, and the cheapest one to assert.
var noClient any

func TestFastEngineResourceRemoved_ErrorsWithMigrationGuidance(t *testing.T) {
	diags := fastEngineResourceRemoved("sysdig_secure_rule_container", "CONTAINER")(
		context.Background(), emptyResourceData(t), noClient)

	require.Len(t, diags, 1)
	assert.True(t, diags.HasError(), "Create/Update must fail hard")
	assert.Equal(t, diag.Error, diags[0].Severity)
	assert.Equal(t, "sysdig_secure_rule_container is no longer supported", diags[0].Summary)
	assert.Contains(t, diags[0].Detail, "CONTAINER")
	assert.Contains(t, diags[0].Detail, fastEngineMigrationTarget)
}

func TestFastEngineDataSourceRemoved_UsesDataSourceWording(t *testing.T) {
	diags := fastEngineDataSourceRemoved("sysdig_secure_rule_process", "PROCESS")(
		context.Background(), emptyResourceData(t), noClient)

	require.Len(t, diags, 1)
	assert.True(t, diags.HasError())
	assert.Equal(t, "data source sysdig_secure_rule_process is no longer supported", diags[0].Summary)
	assert.Contains(t, diags[0].Detail, "PROCESS")
	assert.Contains(t, diags[0].Detail, fastEngineMigrationTarget+" data source")

	// A data source has no state, so resource-only guidance would be nonsense
	// here. This is the wording bug that shipping a single shared handler for
	// both cases originally introduced.
	assert.NotContains(t, diags[0].Detail, "state")
	assert.NotContains(t, diags[0].Detail, "Remove this block")
}

func TestFastEngineRuleGone_ClearsIDAndWarns(t *testing.T) {
	d := emptyResourceData(t)
	d.SetId("12345")

	diags := fastEngineRuleGone("sysdig_secure_rule_network", "NETWORK")(
		context.Background(), d, noClient)

	assert.Empty(t, d.Id(), "Read must report the object as gone so stale state self-heals")
	require.Len(t, diags, 1)
	assert.False(t, diags.HasError(), "Read must not error, or refresh wedges plan/apply/destroy")
	assert.Equal(t, diag.Warning, diags[0].Severity)
	assert.Contains(t, diags[0].Detail, "NETWORK")
}

func TestFastEngineRuleForget_WarnsWithoutError(t *testing.T) {
	d := emptyResourceData(t)
	d.SetId("12345")

	diags := fastEngineRuleForget("sysdig_secure_rule_syscall", "SYSCALL")(
		context.Background(), d, noClient)

	require.Len(t, diags, 1)
	assert.Equal(t, diag.Warning, diags[0].Severity)

	// The SDK's destroy path gates on diags.HasError() and only then clears the
	// ID itself ("Make sure the ID is gone", helper/schema/resource.go). A
	// warning must therefore stay a warning: promote it to an Error and the
	// entry is left in state forever. This is why the handler does not need to
	// call d.SetId("") of its own.
	assert.False(t, diags.HasError(), "an Error here would strand the entry in state")
	assert.Contains(t, diags[0].Detail, "no remote object was deleted")
}

func TestFastEngineDeprecation_NamesKindAndRuleType(t *testing.T) {
	for _, tc := range []struct{ kind, typeName, ruleType string }{
		{"resource", "sysdig_secure_rule_container", "CONTAINER"},
		{"data source", "sysdig_secure_rule_container", "CONTAINER"},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			msg := fastEngineDeprecation(tc.kind, tc.typeName, tc.ruleType)

			assert.Contains(t, msg, tc.kind)
			assert.Contains(t, msg, tc.typeName)
			assert.Contains(t, msg, tc.ruleType)
			assert.Contains(t, msg, fastEngineMigrationTarget)
		})
	}
}

// The five resources and five data sources keep their full argument schemas so
// that an existing configuration reaches the migration error above instead of a
// pile of "Unsupported argument" diagnostics that bury it.
func TestFastEngineResources_KeepSchemaAndAreFullyWired(t *testing.T) {
	for _, tc := range []struct {
		name      string
		resource  *schema.Resource
		dataSrc   *schema.Resource
		ruleType  string
		attribute string
	}{
		{"container", resourceSysdigSecureRuleContainer(), dataSourceSysdigSecureRuleContainer(), "CONTAINER", "containers"},
		{"filesystem", resourceSysdigSecureRuleFilesystem(), dataSourceSysdigSecureRuleFilesystem(), "FILESYSTEM", "read_only"},
		{"network", resourceSysdigSecureRuleNetwork(), dataSourceSysdigSecureRuleNetwork(), "NETWORK", "tcp"},
		{"process", resourceSysdigSecureRuleProcess(), dataSourceSysdigSecureRuleProcess(), "PROCESS", "processes"},
		{"syscall", resourceSysdigSecureRuleSyscall(), dataSourceSysdigSecureRuleSyscall(), "SYSCALL", "syscalls"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.resource
			assert.Contains(t, r.DeprecationMessage, tc.ruleType)
			assert.Contains(t, r.Schema, tc.attribute, "argument schema must be kept intact")
			assert.Contains(t, r.Schema, "name")

			require.NotNil(t, r.CreateContext)
			require.NotNil(t, r.ReadContext)
			require.NotNil(t, r.UpdateContext)
			require.NotNil(t, r.DeleteContext)

			// Import can only ever route into a Read that reports the object as
			// gone, so it is deliberately not offered.
			assert.Nil(t, r.Importer, "these types cannot be imported")

			assert.True(t, r.CreateContext(context.Background(), emptyResourceData(t), noClient).HasError())
			assert.True(t, r.UpdateContext(context.Background(), emptyResourceData(t), noClient).HasError())
			assert.False(t, r.ReadContext(context.Background(), emptyResourceData(t), noClient).HasError())
			assert.False(t, r.DeleteContext(context.Background(), emptyResourceData(t), noClient).HasError())

			ds := tc.dataSrc
			assert.Contains(t, ds.DeprecationMessage, tc.ruleType)
			assert.Contains(t, ds.Schema, tc.attribute)
			require.NotNil(t, ds.ReadContext)
			assert.True(t, ds.ReadContext(context.Background(), emptyResourceData(t), noClient).HasError())
		})
	}
}

package sysdig

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Support for the "fast engine" rule types — CONTAINER, FILESYSTEM, NETWORK,
// PROCESS and SYSCALL — was dropped from the backend along with the
// list-matching policy code they were built on, and /api/secure/rules now
// rejects all five with HTTP 400. The five resources and five data sources
// built on them cannot work against any current backend.
//
// The type names stay registered in provider.go on purpose: unregistering them
// fails with a bare "Invalid resource type" that gives no migration hint, and
// breaks anyone who still has one in a config or state file. Deleting them
// outright is left for the next major release. Until then the entry points
// below stand in, and none of them touch the network.

// fastEngineCRUD is the signature shared by schema.CreateContextFunc,
// ReadContextFunc, UpdateContextFunc and DeleteContextFunc. It is an alias, not
// a definition, so one value satisfies all four named types.
type fastEngineCRUD = func(context.Context, *schema.ResourceData, any) diag.Diagnostics

const fastEngineMigrationTarget = "sysdig_secure_rule_falco"

// fastEngineResourceRemoved backs resource Create and Update, which cannot
// succeed. It replaces the backend's opaque HTTP 400 with local guidance.
func fastEngineResourceRemoved(typeName, ruleType string) fastEngineCRUD {
	return func(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s is no longer supported", typeName),
			Detail: fmt.Sprintf(
				"The Sysdig backend removed the %s rule type and rejects it with HTTP 400.\n\n"+
					"Remove this block and express the same detection as a %s resource with an "+
					"equivalent Falco condition. Any leftover state entry is dropped on the next "+
					"refresh.",
				ruleType, fastEngineMigrationTarget),
		}}
	}
}

// fastEngineDataSourceRemoved backs data source Read. A data source has no
// state to reconcile, so the guidance differs from the resource case: there is
// nothing to look up and nothing to clean up.
func fastEngineDataSourceRemoved(typeName, ruleType string) fastEngineCRUD {
	return func(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("data source %s is no longer supported", typeName),
			Detail: fmt.Sprintf(
				"The Sysdig backend removed the %s rule type, so no rule of it can exist to "+
					"look up.\n\nUse the %s data source instead.",
				ruleType, fastEngineMigrationTarget),
		}}
	}
}

// fastEngineRuleGone backs resource Read, and deliberately does not error:
// Read reports what exists remotely, and the honest answer is "nothing".
// Clearing the ID is how a provider says that, and it lets state written by an
// older provider version self-heal.
//
// Erroring here would wedge the user instead. Refresh runs ahead of plan, apply
// and destroy alike, so one leftover address would make `terraform destroy`
// fail with `terraform state rm` as the only way out. Create and Update still
// refuse loudly, so a config that declares one of these cannot silently do
// nothing.
func fastEngineRuleGone(typeName, ruleType string) fastEngineCRUD {
	return func(_ context.Context, d *schema.ResourceData, _ any) diag.Diagnostics {
		d.SetId("")

		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("%s no longer exists and was dropped from state", typeName),
			Detail: fmt.Sprintf(
				"The Sysdig backend removed the %s rule type, so this resource cannot exist "+
					"remotely. Delete the block from your configuration and express the detection "+
					"as a %s resource instead.",
				ruleType, fastEngineMigrationTarget),
		}}
	}
}

// fastEngineRuleForget backs resource Delete. Returning without an error is
// enough for the SDK to drop the entry from state; the warning is here so a
// clean "Destroy complete" doesn't imply a remote delete that never happened.
func fastEngineRuleForget(typeName, ruleType string) fastEngineCRUD {
	return func(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("%s was removed from state only", typeName),
			Detail: fmt.Sprintf(
				"The %s rule type no longer exists in the Sysdig backend, so no remote object "+
					"was deleted. If a rule of this type predates the removal, delete it by hand "+
					"under Policies > Rules in Sysdig Secure.",
				ruleType),
		}}
	}
}

// fastEngineDeprecation is the plan-time warning, shown whether or not the
// configuration is ever applied. kind is "resource" or "data source".
func fastEngineDeprecation(kind, typeName, ruleType string) string {
	return fmt.Sprintf(
		"%s %s is no longer functional: the Sysdig backend removed the %s rule type and "+
			"rejects it with HTTP 400. Use %s instead. This %s will be deleted in the next "+
			"major version.",
		kind, typeName, ruleType, fastEngineMigrationTarget, kind)
}

package sysdig

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// falcoSnapshotOnce ensures the account-wide snapshot diagnostic (see
// snapshotDiagnostic) is emitted at most once per `terraform apply`
// invocation. Terraform launches a fresh provider process per command
// invocation, so process-lifetime state here is exactly apply-scoped — no
// cross-run leakage, no extra backend calls, no explicit reset needed.
var falcoSnapshotOnce sync.Once

// falcoWarningsToDiagnostics converts Falco static-analyzer warnings returned by
// the Sysdig backend on rule create/update into Terraform Diagnostics.
//
// Today the backend emits LOAD_NO_EVTTYPE (Falco's static analyzer warning
// that a rule's condition matches too many event types); the helper is
// code-agnostic and surfaces whatever the backend returns. Each diagnostic is
// a diag.Warning so `terraform apply` shows them inline without failing the
// apply.
//
// The backend returns the FULL current warning list for the customer's rules
// on every create/update call, not just warnings caused by this specific
// save. Two-tier output keeps that from becoming a wall of text that
// reprints every pre-existing warning on every future apply:
//
//   - Headline: warnings attributed to THIS resource's own rule (ruleName) —
//     one diag.Warning each, always shown. Actionable: this rule needs fixing.
//   - Snapshot: every OTHER warning elsewhere on the account is bundled into
//     a single diag.Warning, emitted at most once per apply via
//     falcoSnapshotOnce. Informational: shows blast radius
//     (EnabledInPolicies) without implying this save caused it.
//
// Wired on sysdig_secure_rule_falco only.
//
// Intentionally NOT wired:
//   - sysdig_secure_rule_container / _filesystem / _network / _process /
//     _syscall — deprecated and non-functional against current backends
//     (Create/Update return HTTP 400 before any validation runs; see their
//     DeprecationMessage). Wiring them would be dead code.
//   - sysdig_secure_rule_stateful — stateful-detection rules don't go through the
//     Falco static-analyzer path and the backend doesn't emit warnings for them.
//   - sysdig_secure_macro / sysdig_secure_list — the backend doesn't surface
//     warnings on those endpoints today. A macro/list's own noise still
//     surfaces via the snapshot on whichever rule save happens to trigger
//     validation next, attributed to the RULE that references it (warnings
//     are always rule-scoped — LOAD_NO_EVTTYPE never attaches to a macro/list
//     itself).
func falcoWarningsToDiagnostics(warnings []v2.FalcoWarning, ruleName string) diag.Diagnostics {
	if len(warnings) == 0 {
		return nil
	}

	var own, other []v2.FalcoWarning
	for _, w := range warnings {
		if w.Rule == ruleName {
			own = append(own, w)
		} else {
			other = append(other, w)
		}
	}

	diags := make(diag.Diagnostics, 0, len(own)+1)
	for _, w := range own {
		diags = append(diags, headlineDiagnostic(w))
	}

	if len(other) > 0 {
		falcoSnapshotOnce.Do(func() {
			diags = append(diags, snapshotDiagnostic(other))
		})
	}

	return diags
}

// headlineDiagnostic renders a single warning attributed to the rule being
// saved in this resource operation.
func headlineDiagnostic(w v2.FalcoWarning) diag.Diagnostic {
	detail := w.Message
	if len(w.EnabledInPolicies) > 0 {
		detail = fmt.Sprintf("%s\n\nEnabled in policies: %s", detail, strings.Join(sortedCopy(w.EnabledInPolicies), ", "))
	}
	if len(w.AgentVersions) > 0 {
		detail = fmt.Sprintf("%s\nAgent versions: %s", detail, strings.Join(sortedCopy(w.AgentVersions), ", "))
	}
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("Falco static-analyzer warning (%s) on rule %q", w.Code, w.Rule),
		Detail:   detail,
	}
}

// snapshotDiagnostic bundles every warning NOT attributed to the
// currently-saved rule into a single diagnostic. Sorted by blast radius
// (EnabledInPolicies count) descending, then by rule name for determinism —
// customers see their biggest-impact noisy rule first.
func snapshotDiagnostic(warnings []v2.FalcoWarning) diag.Diagnostic {
	sorted := make([]v2.FalcoWarning, len(warnings))
	copy(sorted, warnings)
	sort.SliceStable(sorted, func(i, j int) bool {
		if len(sorted[i].EnabledInPolicies) != len(sorted[j].EnabledInPolicies) {
			return len(sorted[i].EnabledInPolicies) > len(sorted[j].EnabledInPolicies)
		}
		return sorted[i].Rule < sorted[j].Rule
	})

	lines := make([]string, 0, len(sorted))
	for _, w := range sorted {
		line := fmt.Sprintf("- %s (%s)", w.Rule, w.Code)
		if len(w.EnabledInPolicies) > 0 {
			line = fmt.Sprintf("%s — enabled in policies: %s", line, strings.Join(sortedCopy(w.EnabledInPolicies), ", "))
		}
		lines = append(lines, line)
	}

	noun, verb := "other Falco rules", "emit"
	if len(sorted) == 1 {
		noun, verb = "other Falco rule", "emits"
	}
	return diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  fmt.Sprintf("%d %s currently %s static-analyzer warnings", len(sorted), noun, verb),
		Detail:   strings.Join(lines, "\n"),
	}
}

// sortedCopy returns a sorted copy of s so diagnostic text built from
// backend-ordered slices (EnabledInPolicies, AgentVersions) stays stable
// across applies regardless of the order the backend returns them in.
func sortedCopy(s []string) []string {
	out := make([]string, len(s))
	copy(out, s)
	sort.Strings(out)
	return out
}

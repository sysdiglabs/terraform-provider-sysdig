package sysdig

import (
	"sync"
	"testing"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFalcoWarningsToDiagnostics_HeadlineOnly(t *testing.T) {
	resetFalcoSnapshotOnce(t)

	warnings := []v2.FalcoWarning{
		{Code: "LOAD_NO_EVTTYPE", Rule: "my-rule", Message: "too broad"},
	}

	diags := falcoWarningsToDiagnostics(warnings, "my-rule")

	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, `on rule "my-rule"`)
}

func TestHeadlineDiagnostic_SortsEnabledInPoliciesAndAgentVersions(t *testing.T) {
	w := v2.FalcoWarning{
		Code:              "LOAD_NO_EVTTYPE",
		Rule:              "my-rule",
		Message:           "too broad",
		EnabledInPolicies: []string{"Zebra Policy", "Alpha Policy"},
		AgentVersions:     []string{"linux-14.7.0", "linux-13.8.0"},
	}

	d := headlineDiagnostic(w)

	// Backend order is Zebra, Alpha / 14.7.0, 13.8.0 — Detail must render them
	// sorted so the diagnostic text is stable across applies regardless of
	// the order the backend happens to return.
	assert.Contains(t, d.Detail, "Enabled in policies: Alpha Policy, Zebra Policy")
	assert.Contains(t, d.Detail, "Agent versions: linux-13.8.0, linux-14.7.0")
}

func TestFalcoWarningsToDiagnostics_HeadlinePlusSnapshot(t *testing.T) {
	resetFalcoSnapshotOnce(t)

	warnings := []v2.FalcoWarning{
		{Code: "LOAD_NO_EVTTYPE", Rule: "my-rule", Message: "too broad"},
		{Code: "LOAD_NO_EVTTYPE", Rule: "other-rule-1", Message: "also broad", EnabledInPolicies: []string{"Policy A"}},
		{Code: "LOAD_NO_EVTTYPE", Rule: "other-rule-2", Message: "also broad too"},
	}

	diags := falcoWarningsToDiagnostics(warnings, "my-rule")

	require.Len(t, diags, 2, "expect exactly 1 headline + 1 bundled snapshot, not 3 separate diagnostics")
	assert.Contains(t, diags[0].Summary, `on rule "my-rule"`)
	assert.Contains(t, diags[1].Summary, "2 other Falco rules")
	assert.Contains(t, diags[1].Detail, "other-rule-1")
	assert.Contains(t, diags[1].Detail, "other-rule-2")
}

func TestFalcoWarningsToDiagnostics_SnapshotEmittedOnceAcrossCalls(t *testing.T) {
	resetFalcoSnapshotOnce(t)

	// Simulates two resources being saved in the same apply: both calls see
	// the same account-wide "other" warning, but the snapshot must only
	// render once — on the second call it must be a headline-only result.
	warningsForFirstRule := []v2.FalcoWarning{
		{Code: "LOAD_NO_EVTTYPE", Rule: "rule-a", Message: "a is broad"},
		{Code: "LOAD_NO_EVTTYPE", Rule: "leftover-rule", Message: "pre-existing noise"},
	}
	warningsForSecondRule := []v2.FalcoWarning{
		{Code: "LOAD_NO_EVTTYPE", Rule: "rule-b", Message: "b is broad"},
		{Code: "LOAD_NO_EVTTYPE", Rule: "leftover-rule", Message: "pre-existing noise"},
	}

	first := falcoWarningsToDiagnostics(warningsForFirstRule, "rule-a")
	second := falcoWarningsToDiagnostics(warningsForSecondRule, "rule-b")

	require.Len(t, first, 2, "first call: headline for rule-a + snapshot for leftover-rule")
	require.Len(t, second, 1, "second call: headline for rule-b only — snapshot already emitted this apply")
	assert.Contains(t, second[0].Summary, `on rule "rule-b"`)
}

func TestFalcoWarningsToDiagnostics_NoWarnings(t *testing.T) {
	resetFalcoSnapshotOnce(t)
	assert.Nil(t, falcoWarningsToDiagnostics(nil, "my-rule"))
}

func TestSnapshotDiagnostic_SortedByPolicyCountDesc(t *testing.T) {
	warnings := []v2.FalcoWarning{
		{Rule: "z-rule-no-policies", Code: "LOAD_NO_EVTTYPE"},
		{Rule: "a-rule-many-policies", Code: "LOAD_NO_EVTTYPE", EnabledInPolicies: []string{"P1", "P2", "P3"}},
		{Rule: "m-rule-one-policy", Code: "LOAD_NO_EVTTYPE", EnabledInPolicies: []string{"P1"}},
	}

	diag := snapshotDiagnostic(warnings)

	lines := diag.Detail
	// a-rule-many-policies (3 policies) should appear before m-rule-one-policy (1),
	// which should appear before z-rule-no-policies (0).
	idxA := indexOf(lines, "a-rule-many-policies")
	idxM := indexOf(lines, "m-rule-one-policy")
	idxZ := indexOf(lines, "z-rule-no-policies")
	require.True(t, idxA >= 0 && idxM >= 0 && idxZ >= 0)
	assert.Less(t, idxA, idxM)
	assert.Less(t, idxM, idxZ)
}

// resetFalcoSnapshotOnce clears the process-lifetime dedup guard between
// test cases. In production this state is scoped to one `terraform apply`
// process; tests need to reset it since they share a process.
func resetFalcoSnapshotOnce(t *testing.T) {
	t.Helper()
	falcoSnapshotOnce = sync.Once{}
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

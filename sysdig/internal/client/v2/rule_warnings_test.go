package v2

import (
	"encoding/json"
	"testing"
)

// The payloads below are real HTTP response bodies captured from a live
// Sysdig Secure backend (POST/PUT /api/secure/rules), not hand-typed to match
// what we assume the JSON tags produce. This guards against a struct-tag
// typo silently deserializing a field to its zero value while every
// hand-built-struct test still passes.

const capturedCreateRuleResponse = `{
  "id": 10010392,
  "name": "ROUNDTRIP-capture-fixture-rule",
  "origin": "Secure UI",
  "versionId": "0.0.0",
  "filename": "falco_rules_local.yaml",
  "description": "captured for PR 733 review round-trip test",
  "details": {
    "ruleType": "FALCO",
    "condition": {
      "condition": "proc.name=nginx"
    },
    "description": "captured for PR 733 review round-trip test",
    "output": "nginx activity (proc.cmdline=%proc.cmdline)",
    "priority": "NOTICE",
    "source": "syscall",
    "append": false
  },
  "tags": [],
  "version": 1,
  "createdOn": 1783702635571,
  "modifiedOn": 1783702635571,
  "warnings": [
    {
      "code": "LOAD_NO_EVTTYPE",
      "rule": "ROUNDTRIP-capture-fixture-rule",
      "message": "Rule matches too many evt.type values. This has a significant performance penalty.",
      "enabled_in_policies": null,
      "agent_versions": [
        "linux-14.7.0"
      ]
    }
  ]
}`

// Same rule, after being attached to a policy and re-saved — captures the
// enabled_in_policies-populated case.
const capturedUpdateRuleResponseWithPolicy = `{
  "id": 10010392,
  "name": "ROUNDTRIP-capture-fixture-rule",
  "origin": "Secure UI",
  "versionId": "0.0.0",
  "filename": "falco_rules_local.yaml",
  "description": "captured for PR 733 review round-trip test (updated)",
  "details": {
    "ruleType": "FALCO",
    "condition": {
      "condition": "proc.name=nginx"
    },
    "description": "captured for PR 733 review round-trip test (updated)",
    "output": "nginx activity (proc.cmdline=%proc.cmdline)",
    "priority": "NOTICE",
    "source": "syscall",
    "append": false
  },
  "tags": [],
  "version": 2,
  "createdOn": 1783702635571,
  "modifiedOn": 1783702842138,
  "warnings": [
    {
      "code": "LOAD_NO_EVTTYPE",
      "rule": "ROUNDTRIP-capture-fixture-rule",
      "message": "Rule matches too many evt.type values. This has a significant performance penalty.",
      "enabled_in_policies": [
        "ROUNDTRIP-capture-policy"
      ],
      "agent_versions": [
        "linux-14.7.0"
      ]
    }
  ]
}`

func TestRule_UnmarshalWarnings_RealPayload_NoPolicyAttachment(t *testing.T) {
	var rule Rule
	if err := json.Unmarshal([]byte(capturedCreateRuleResponse), &rule); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(rule.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(rule.Warnings))
	}
	w := rule.Warnings[0]
	if w.Code != "LOAD_NO_EVTTYPE" {
		t.Errorf("Code = %q, want LOAD_NO_EVTTYPE", w.Code)
	}
	if w.Rule != "ROUNDTRIP-capture-fixture-rule" {
		t.Errorf("Rule = %q, want ROUNDTRIP-capture-fixture-rule", w.Rule)
	}
	if w.EnabledInPolicies != nil {
		t.Errorf("EnabledInPolicies = %v, want nil (rule not yet attached to any policy)", w.EnabledInPolicies)
	}
	if len(w.AgentVersions) != 1 || w.AgentVersions[0] != "linux-14.7.0" {
		t.Errorf("AgentVersions = %v, want [linux-14.7.0]", w.AgentVersions)
	}
}

func TestRule_UnmarshalWarnings_RealPayload_WithPolicyAttachment(t *testing.T) {
	var rule Rule
	if err := json.Unmarshal([]byte(capturedUpdateRuleResponseWithPolicy), &rule); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if len(rule.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(rule.Warnings))
	}
	w := rule.Warnings[0]
	if len(w.EnabledInPolicies) != 1 || w.EnabledInPolicies[0] != "ROUNDTRIP-capture-policy" {
		t.Errorf("EnabledInPolicies = %v, want [ROUNDTRIP-capture-policy]", w.EnabledInPolicies)
	}
}

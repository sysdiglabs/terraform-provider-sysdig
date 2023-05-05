package secure

import (
	"bytes"
	"encoding/json"
	"io"
)

// -------- Policies --------

type Policy struct {
	ID                     int      `json:"id,omitempty"`
	Name                   string   `json:"name"`
	Description            string   `json:"description"`
	Severity               int      `json:"severity"`
	Enabled                bool     `json:"enabled"`
	RuleNames              []string `json:"ruleNames"`
	Actions                []Action `json:"actions"`
	Scope                  string   `json:"scope,omitempty"`
	Version                int      `json:"version,omitempty"`
	NotificationChannelIds []int    `json:"notificationChannelIds"`
	Type                   string   `json:"type"`
	Runbook                string   `json:"runbook"`
}

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	Name                 string `json:"name,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer"`
	Type                 string `json:"type"`
}

func (policy *Policy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func PolicyFromJSON(body []byte) (result Policy) {
	_ = json.Unmarshal(body, &result)

	return result
}

// -------- VulnerabilityExceptionList --------

type VulnerabilityExceptionList struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

func (l *VulnerabilityExceptionList) ToJSON() io.Reader {
	payload, _ := json.Marshal(*l)
	return bytes.NewBuffer(payload)
}

func VulnerabilityExceptionListFromJSON(body []byte) *VulnerabilityExceptionList {
	var result VulnerabilityExceptionList
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- VulnerabilityException --------

type VulnerabilityException struct {
	ID             string `json:"id"`
	Gate           string `json:"gate"`
	TriggerID      string `json:"trigger_id"`
	Notes          string `json:"notes"`
	ExpirationDate *int   `json:"expiration_date,omitempty"`
	Enabled        bool   `json:"enabled"`
}

func (e *VulnerabilityException) ToJSON() io.Reader {
	payload, _ := json.Marshal(*e)
	return bytes.NewBuffer(payload)
}

func VulnerabilityExceptionFromJSON(body []byte) *VulnerabilityException {
	var result VulnerabilityException
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- CloudAccount --------

type CloudAccount struct {
	AccountID                    string `json:"accountId"`
	Provider                     string `json:"provider"`
	Alias                        string `json:"alias"`
	RoleAvailable                bool   `json:"roleAvailable"`
	RoleName                     string `json:"roleName"`
	ExternalID                   string `json:"externalId,omitempty"`
	WorkLoadIdentityAccountID    string `json:"workloadIdentityAccountId,omitempty"`
	WorkLoadIdentityAccountAlias string `json:"workLoadIdentityAccountAlias,omitempty"`
}

func (e *CloudAccount) ToJSON() io.Reader {
	payload, _ := json.Marshal(*e)
	return bytes.NewBuffer(payload)
}

func CloudAccountFromJSON(body []byte) *CloudAccount {
	var result CloudAccount
	_ = json.Unmarshal(body, &result)

	return &result
}

// -------- Scanning Policies --------
type ScanningPolicy struct {
	ID             string         `json:"id,omitempty"`
	Version        string         `json:"version,omitempty"`
	Name           string         `json:"name"`
	Comment        string         `json:"comment"`
	IsDefault      bool           `json:"isDefault,omitempty"`
	PolicyBundleId string         `json:"policyBundleId,omitempty"`
	Rules          []ScanningGate `json:"rules"`
}

type ScanningGate struct {
	ID      string              `json:"id,omitempty"`
	Gate    string              `json:"gate"`
	Trigger string              `json:"trigger"`
	Action  string              `json:"action"`
	Params  []ScanningGateParam `json:"params"`
}

type ScanningGateParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (policy *ScanningPolicy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func ScanningPolicyFromJSON(body []byte) (result ScanningPolicy) {
	_ = json.Unmarshal(body, &result)
	return result
}

// -------- Scanning Policy Assignments --------
type ScanningPolicyAssignmentList struct {
	Items          []ScanningPolicyAssignment `json:"items"`
	PolicyBundleId string                     `json:"policyBundleId"`
}

type ScanningPolicyAssignment struct {
	ID           string                        `json:"id,omitempty"`
	Name         string                        `json:"name"`
	Registry     string                        `json:"registry"`
	Repository   string                        `json:"repository"`
	Image        ScanningPolicyAssignmentImage `json:"image"`
	PolicyIDs    []string                      `json:"policy_ids"`
	WhitelistIDs []string                      `json:"whitelist_ids"`
}

type ScanningPolicyAssignmentImage struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (policy *ScanningPolicyAssignmentList) ToJSON() io.Reader {
	payload, _ := json.Marshal(policy)
	return bytes.NewBuffer(payload)
}

func ScanningPolicyAssignmentFromJSON(body []byte) (result ScanningPolicyAssignmentList) {
	_ = json.Unmarshal(body, &result)
	return result
}

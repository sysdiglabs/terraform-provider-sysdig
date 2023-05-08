package secure

import (
	"bytes"
	"encoding/json"
	"io"
)

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

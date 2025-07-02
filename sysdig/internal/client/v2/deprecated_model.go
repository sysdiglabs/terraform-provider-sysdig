package v2

type DeprecatedVulnerabilityExceptionList struct {
	ID      string `json:"id,omitempty"`
	Version string `json:"version"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type DeprecatedVulnerabilityException struct {
	ID             string `json:"id"`
	Gate           string `json:"gate"`
	TriggerID      string `json:"trigger_id"`
	Notes          string `json:"notes"`
	ExpirationDate *int   `json:"expiration_date,omitempty"`
	Enabled        bool   `json:"enabled"`
}

type DeprecatedScanningPolicy struct {
	ID             string                   `json:"id,omitempty"`
	Version        string                   `json:"version,omitempty"`
	Name           string                   `json:"name"`
	Comment        string                   `json:"comment"`
	IsDefault      bool                     `json:"isDefault,omitempty"`
	PolicyBundleId string                   `json:"policyBundleId,omitempty"`
	Rules          []DeprecatedScanningGate `json:"rules"`
}

type DeprecatedScanningGate struct {
	ID      string                        `json:"id,omitempty"`
	Gate    string                        `json:"gate"`
	Trigger string                        `json:"trigger"`
	Action  string                        `json:"action"`
	Params  []DeprecatedScanningGateParam `json:"params"`
}

type DeprecatedScanningGateParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeprecatedScanningPolicyAssignmentList struct {
	Items          []DeprecatedScanningPolicyAssignment `json:"items"`
	PolicyBundleId string                               `json:"policyBundleId"`
}

type DeprecatedScanningPolicyAssignment struct {
	ID           string                                  `json:"id,omitempty"`
	Name         string                                  `json:"name"`
	Registry     string                                  `json:"registry"`
	Repository   string                                  `json:"repository"`
	Image        DeprecatedScanningPolicyAssignmentImage `json:"image"`
	PolicyIDs    []string                                `json:"policy_ids"`
	WhitelistIDs []string                                `json:"whitelist_ids"`
}

type DeprecatedScanningPolicyAssignmentImage struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

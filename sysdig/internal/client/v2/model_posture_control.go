package v2

type SaveControlRequest struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	ResourceKind       string `json:"resourceKind"`
	Severity           string `json:"severity"`
	Rego               string `json:"rego"`
	RemediationDetails string `json:"remediationDetails"`
}

type SaveControlResponse struct {
	Data PostureControl `json:"data"`
}

type PostureControl struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	ResourceKind       string `json:"resourceKindDisplayName"`
	Severity           string `json:"severity"`
	Rego               string `json:"rego"`
	RemediationDetails string `json:"remediationDetails"`
}

type AccepetPostureRiskRequest struct {
	AcceptanceID string `json:"id"`
	ControlName  string `json:"controlName"`
	ZoneName     string `json:"zoneName"`
	Description  string `json:"description"`
	Filter       string `json:"filter"`
	Reason       string `json:"reason"`
	ExpiresAt    int64  `json:"expiresAt"`
}

type UpdateAccepetPostureRiskRequest struct {
	AcceptanceID string                        `json:"id"`
	Acceptance   UpdateAcceptPostureRiskFields `json:"riskAcceptance"`
}

type UpdateAccepetPostureResponse struct {
	Acceptance AcceptPostureRisk `json:"riskAcceptance"`
}

type AcceptPostureRisk struct {
	AcceptanceID    string `json:"id"`
	ControlName     string `json:"controlName"`
	ZoneName        string `json:"zoneName"`
	Description     string `json:"description"`
	Filter          string `json:"filter"`
	Reason          string `json:"reason"`
	ExpiresAt       string `json:"expiresAt"`
	AcceeptanceDate string `json:"acceptanceDate"`
	UserName        string `json:"username"`
	Type            string `json:"type"`
	IsExpired       bool   `json:"isExpired"`
	IsSystem        bool   `json:"isSystem"`
	AcceptPeriod    string `json:"acceptPeriod"`
}

type UpdateAcceptPostureRiskFields struct {
	Description  string `json:"description"`
	Reason       string `json:"reason"`
	ExpiresAt    string `json:"expiresAt"`
	AcceptPeriod string `json:"acceptPeriod"`
}

type AcceptPostureRiskResponse struct {
	Data AcceptPostureRisk `json:"data"`
}

type DeleteAcceptPostureRisk struct {
	AcceptanceID string `json:"id"`
	Filter       string `json:"filter"`
}

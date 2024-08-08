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

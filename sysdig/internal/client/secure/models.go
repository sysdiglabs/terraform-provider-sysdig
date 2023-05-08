package secure

import (
	"bytes"
	"encoding/json"
	"io"
)

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

package secure

import (
	"bytes"
	"encoding/json"
	"io"
)

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

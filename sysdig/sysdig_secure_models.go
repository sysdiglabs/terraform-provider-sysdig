package sysdig

import (
	"bytes"
	"encoding/json"
	"io"
)

type Action struct {
	AfterEventNs         int    `json:"afterEventNs,omitempty"`
	BeforeEventNs        int    `json:"beforeEventNs,omitempty"`
	IsLimitedToContainer bool   `json:"isLimitedToContainer,omitempty"`
	Type                 string `json:"type"`
}

type FalcoConfiguration struct {
	RuleNameRegEx string `json:"ruleNameRegEx"`
}

type Policy struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	Severity           int                `json:"severity"`
	ContainerScope     bool               `json:"containerScope"`
	HostScope          bool               `json:"hostScope"`
	Enabled            bool               `json:"enabled"`
	Actions            []Action           `json:"actions,omitempty"`
	Scope              string             `json:"scope,omitempty"`
	FalcoConfiguration FalcoConfiguration `json:"falcoConfiguration,omitempty"`
	Version            int                `json:"version,omitempty"`
}

type policyWrapper struct {
	Policy Policy `json:"policy"`
}

func (policy *Policy) ToJSON() io.Reader {
	payload, _ := json.Marshal(policyWrapper{*policy})
	return bytes.NewBuffer(payload)
}

func PolicyFromJSON(body []byte) Policy {
	var result policyWrapper
	json.Unmarshal(body, &result)

	return result.Policy
}

type UserRulesFile struct {
	Content string `json:"content"`
	Version int    `json:"version"`
}

type userRulesFileWrapper struct {
	UserRulesFile UserRulesFile `json:"userRulesFile"`
}

func (userRulesFile *UserRulesFile) ToJSON() io.Reader {
	payload, _ := json.Marshal(userRulesFileWrapper{*userRulesFile})
	return bytes.NewBuffer(payload)
}

func UserRulesFileFromJSON(body []byte) UserRulesFile {
	var result userRulesFileWrapper
	json.Unmarshal(body, &result)

	return result.UserRulesFile
}

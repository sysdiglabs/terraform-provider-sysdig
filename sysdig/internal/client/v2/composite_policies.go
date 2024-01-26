package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	skipPolicyV2MsgFlag       = "skipPolicyV2Msg"
	CreateCompositePolicyPath = "%s/api/v2/policies/batch?%s=%t"
	DeleteCompositePolicyPath = "%s/api/v2/policies/batch/%d?%s=%t"
	UpdateCompositePolicyPath = "%s/api/v2/policies/batch/%d?%s=%t"
	// TODO
	GetCompositePolicyPath   = "%s/api/v2/policies/%d"
	GetCompositePoliciesPath = "%s/api/v2/policies?policyType=malware&limit=200&filter=Test+Malware+Policy" // TODO: Implement pagination

	GetCompositePolicyRulesPath = "%s/api/policies/v3/rules/groups?%s"
)

type CompositePolicyInterface interface {
	Base
	CreateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	DeleteCompositePolicy(ctx context.Context, policyID int) error
	UpdateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	// TODO
	GetCompositePolicyByID(ctx context.Context, policyID int) (PolicyRulesComposite, int, error)
	GetCompositePolicies(ctx context.Context) ([]PolicyRulesComposite, int, error)
}

func (client *Client) CreateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error) {
	payload, err := Marshal(policy)
	if err != nil {
		return PolicyRulesComposite{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateCompositePolicyURL(), payload)
	if err != nil {
		return PolicyRulesComposite{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[PolicyRulesComposite](response.Body)
}

func (client *Client) UpdateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error) {
	payload, err := Marshal(policy)
	if err != nil {
		return PolicyRulesComposite{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateCompositePolicyURL(policy.Policy.ID), payload)
	if err != nil {
		return PolicyRulesComposite{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[PolicyRulesComposite](response.Body)
}

func (client *Client) DeleteCompositePolicy(ctx context.Context, policyID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteCompositePolicyURL(policyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) GetCompositePolicyByID(ctx context.Context, policyID int) (PolicyRulesComposite, int, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCompositePolicyURL(policyID), nil)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, response.StatusCode, client.ErrorFromResponse(response)
	}

	policy, err := Unmarshal[Policy](response.Body)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}

	names := []string{}
	for _, rule := range policy.Rules {
		names = append(names, rule.Name)
	}

	rules, _, err := client.GetCompositePolicyRulesByName(ctx, names)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}

	return PolicyRulesComposite{
		Policy: &policy,
		Rules:  rules,
	}, http.StatusOK, nil
}

func (client *Client) GetCompositePolicyRulesByName(ctx context.Context, names []string) ([]*RuntimePolicyRule, int, error) {
	if len(names) == 0 {
		return nil, http.StatusOK, errors.New("Please provide at least one rule name")
	}

	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCompositePolicyRulesURL(names), nil)
	if err != nil {
		return []*RuntimePolicyRule{}, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []*RuntimePolicyRule{}, response.StatusCode, client.ErrorFromResponse(response)
	}

	unmarshalled, err := Unmarshal[[][]*RuntimePolicyRule](response.Body)
	if err != nil {
		return []*RuntimePolicyRule{}, 0, err
	}

	rules := []*RuntimePolicyRule{}
	for _, arr := range unmarshalled {
		for _, item := range arr {
			rules = append(rules, item)
		}
	}

	if len(rules) == 0 {
		return nil, http.StatusOK, errors.New("Rules not found")
	}

	return rules, http.StatusOK, nil
}

// This method is used in a data source to retrieve a policy by name and type.
// We must retrieve and iterate over all policies, as there is no endpoint to get a policy by name.
func (client *Client) GetCompositePolicies(ctx context.Context) ([]PolicyRulesComposite, int, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCompositePoliciesURL(), nil)
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []PolicyRulesComposite{}, response.StatusCode, client.ErrorFromResponse(response)
	}

	policies, err := Unmarshal[[]Policy](response.Body) // TODO
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}

	compositePoliciesByPolicyID := map[int]PolicyRulesComposite{}
	policyIDByRuleName := map[string]int{}
	names := []string{}
	for _, policy := range policies {
		x := policy
		compositePoliciesByPolicyID[x.ID] = PolicyRulesComposite{
			Policy: &x,
			Rules:  []*RuntimePolicyRule{},
		}

		for _, rule := range x.Rules {
			policyIDByRuleName[rule.Name] = x.ID
			names = append(names, rule.Name)
		}
	}

	rules, _, err := client.GetCompositePolicyRulesByName(ctx, names)
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}

	if len(rules) != len(names) {
		return []PolicyRulesComposite{}, 0, fmt.Errorf("Some rules were not found: %d != %d", len(rules), len(names))
	}

	for _, rule := range rules {
		policyID := policyIDByRuleName[rule.Name]
		p := compositePoliciesByPolicyID[policyID]
		p.Rules = append(p.Rules, rule)
		compositePoliciesByPolicyID[policyID] = p
	}

	policiesFull := []PolicyRulesComposite{}
	for _, policy := range compositePoliciesByPolicyID {
		policiesFull = append(policiesFull, policy)
	}

	return policiesFull, http.StatusOK, nil
}

func (client *Client) CreateCompositePolicyURL() string {
	return fmt.Sprintf(CreateCompositePolicyPath, client.config.url, skipPolicyV2MsgFlag, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) DeleteCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(DeleteCompositePolicyPath, client.config.url, policyID, skipPolicyV2MsgFlag, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) UpdateCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(UpdateCompositePolicyPath, client.config.url, policyID, skipPolicyV2MsgFlag, client.config.secureSkipPolicyV2Msg)
}

// TODO
func (client *Client) GetCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(GetCompositePolicyPath, client.config.url, policyID)
}

// TODO
func (client *Client) GetCompositePoliciesURL() string {
	return fmt.Sprintf(GetCompositePoliciesPath, client.config.url)
}

func (client *Client) GetCompositePolicyRulesURL(names []string) string {
	items := []string{}
	for _, name := range names {
		items = append(items, fmt.Sprintf("names=%s", name))
	}
	return fmt.Sprintf(GetCompositePolicyRulesPath, client.config.url, strings.Join(items[:], "&"))
}

package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	skipPolicyV2MsgFlag       = "skipPolicyV2Msg"
	CreateCompositePolicyPath = "%s/api/v2/policies/batch?%s=%t"
	DeleteCompositePolicyPath = "%s/api/v2/policies/batch/%d?%s=%t"
	UpdateCompositePolicyPath = "%s/api/v2/policies/batch/%d"

	GetCompositePolicyPath   = "%s/api/v2/policies/%d" // TODO: Add skip query param
	GetCompositePoliciesPath = "%s/api/v2/policies?%s" // TODO: Implement pagination otherwise up to getPoliciesLimit number of policies will be returned

	GetCompositePolicyRulesPath = "%s/api/policies/v3/rules/groups?%s"
	getPoliciesLimit            = 1000 // TODO: What is a good limit?
)

type CompositePolicyInterface interface {
	Base
	CreateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	DeleteCompositePolicy(ctx context.Context, policyID int) error
	UpdateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	GetCompositePolicyByID(ctx context.Context, policyID int) (PolicyRulesComposite, int, error)
	FilterCompositePoliciesByNameAndType(ctx context.Context, policyType string, policyName string) ([]PolicyRulesComposite, int, error)
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
		rules = append(rules, arr...)
	}

	if len(rules) == 0 {
		return nil, http.StatusOK, errors.New("Rules not found")
	}

	return rules, http.StatusOK, nil
}

// This method is used in a data source to retrieve a policy by name and type.
// We must retrieve and iterate over all policies, as there is no endpoint to get a policy by name.
func (client *Client) FilterCompositePoliciesByNameAndType(ctx context.Context, policyType string, policyName string) ([]PolicyRulesComposite, int, error) {
	// TODO: Implement pagination in order to get all policies?
	q := GetPoliciesQueryParams{policyType, policyName, getPoliciesLimit}
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetCompositePoliciesURL(q), nil)
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
	return fmt.Sprintf(UpdateCompositePolicyPath, client.config.url, policyID)
}

func (client *Client) GetCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(GetCompositePolicyPath, client.config.url, policyID)
}

func (client *Client) GetCompositePoliciesURL(queryParams GetPoliciesQueryParams) string {
	return fmt.Sprintf(GetCompositePoliciesPath, client.config.url, queryParams.Encode())
}

func (client *Client) GetCompositePolicyRulesURL(names []string) string {
	items := []string{}
	for _, name := range names {
		items = append(items, fmt.Sprintf("names=%s", name))
	}
	return fmt.Sprintf(GetCompositePolicyRulesPath, client.config.url, strings.Join(items[:], "&"))
}

type GetPoliciesQueryParams struct {
	PolicyType string
	Filter     string
	Limit      int
}

func (q *GetPoliciesQueryParams) Encode() string {
	return fmt.Sprintf("policyType=%s&filter=%s&limit=%d", q.PolicyType, strings.Replace(q.Filter, " ", "+", -1), q.Limit)
}

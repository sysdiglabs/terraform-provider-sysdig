package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	skipPolicyV2MsgFlag = "skipPolicyV2Msg"
	getPoliciesLimit    = 1000 // TODO: What is a good limit?
)

const (
	createCompositePolicyPath = "%s/api/v2/policies/batch?%s=%t"
	deleteCompositePolicyPath = "%s/api/v2/policies/batch/%d?%s=%t"
	updateCompositePolicyPath = "%s/api/v2/policies/batch/%d"

	getCompositePolicyPath   = "%s/api/v2/policies/%d" // TODO: Add skip query param
	getCompositePoliciesPath = "%s/api/v2/policies?%s" // TODO: Implement pagination otherwise up to getPoliciesLimit number of policies will be returned

	getCompositePolicyRulesPath = "%s/api/policies/v3/rules/groups?%s"
)

type CompositePolicyInterface interface {
	Base
	CreateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	DeleteCompositePolicy(ctx context.Context, policyID int) error
	UpdateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (PolicyRulesComposite, error)
	GetCompositePolicyByID(ctx context.Context, policyID int) (PolicyRulesComposite, int, error)
	ListCompositePoliciesByNameAndType(ctx context.Context, policyType string, policyName string) ([]PolicyRulesComposite, int, error)
}

func (c *Client) CreateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (policyComposite PolicyRulesComposite, err error) {
	payload, err := Marshal(policy)
	if err != nil {
		return PolicyRulesComposite{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createCompositePolicyURL(), payload)
	if err != nil {
		return PolicyRulesComposite{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[PolicyRulesComposite](response.Body)
}

func (c *Client) UpdateCompositePolicy(ctx context.Context, policy PolicyRulesComposite) (policyComposite PolicyRulesComposite, err error) {
	payload, err := Marshal(policy)
	if err != nil {
		return PolicyRulesComposite{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateCompositePolicyURL(policy.Policy.ID), payload)
	if err != nil {
		return PolicyRulesComposite{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[PolicyRulesComposite](response.Body)
}

func (c *Client) DeleteCompositePolicy(ctx context.Context, policyID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteCompositePolicyURL(policyID), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return c.ErrorFromResponse(response)
	}

	return err
}

func (c *Client) GetCompositePolicyByID(ctx context.Context, policyID int) (policyComposite PolicyRulesComposite, statusCode int, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCompositePolicyURL(policyID), nil)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return PolicyRulesComposite{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	policy, err := Unmarshal[Policy](response.Body)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}

	names := []string{}
	for _, rule := range policy.Rules {
		names = append(names, rule.Name)
	}

	rules, _, err := c.getCompositePolicyRulesByName(ctx, names)
	if err != nil {
		return PolicyRulesComposite{}, 0, err
	}

	return PolicyRulesComposite{
		Policy: &policy,
		Rules:  rules,
	}, http.StatusOK, nil
}

func (c *Client) getCompositePolicyRulesByName(ctx context.Context, names []string) (policies []*RuntimePolicyRule, statusCode int, err error) {
	if len(names) == 0 {
		return nil, http.StatusOK, errors.New("please provide at least one rule name")
	}

	response, err := c.requester.Request(ctx, http.MethodGet, c.getCompositePolicyRulesURL(names), nil)
	if err != nil {
		return []*RuntimePolicyRule{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return []*RuntimePolicyRule{}, response.StatusCode, c.ErrorFromResponse(response)
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
		return nil, http.StatusOK, errors.New("rules not found")
	}

	return rules, http.StatusOK, nil
}

// ListCompositePoliciesByNameAndType is used in a data source to retrieve a policy by name and type.
// We must retrieve and iterate over all policies, as there is no endpoint to get a policy by name.
func (c *Client) ListCompositePoliciesByNameAndType(ctx context.Context, policyType string, policyName string) (list []PolicyRulesComposite, statusCode int, err error) {
	// TODO: Implement pagination in order to get all policies?
	q := getPoliciesQueryParams{policyType, policyName, getPoliciesLimit}
	response, err := c.requester.Request(ctx, http.MethodGet, c.getCompositePoliciesURL(q), nil)
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return []PolicyRulesComposite{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	policies, err := Unmarshal[[]Policy](response.Body) // TODO
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}

	if len(policies) == 0 {
		return []PolicyRulesComposite{}, 0, fmt.Errorf("Policy was not found: %s %s", policyType, policyName)
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

	rules, _, err := c.getCompositePolicyRulesByName(ctx, names)
	if err != nil {
		return []PolicyRulesComposite{}, 0, err
	}

	if len(rules) != len(names) {
		return []PolicyRulesComposite{}, 0, fmt.Errorf("some rules were not found: %d != %d", len(rules), len(names))
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

func (c *Client) createCompositePolicyURL() string {
	return fmt.Sprintf(createCompositePolicyPath, c.config.url, skipPolicyV2MsgFlag, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) deleteCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(deleteCompositePolicyPath, c.config.url, policyID, skipPolicyV2MsgFlag, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) updateCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(updateCompositePolicyPath, c.config.url, policyID)
}

func (c *Client) getCompositePolicyURL(policyID int) string {
	return fmt.Sprintf(getCompositePolicyPath, c.config.url, policyID)
}

func (c *Client) getCompositePoliciesURL(queryParams getPoliciesQueryParams) string {
	return fmt.Sprintf(getCompositePoliciesPath, c.config.url, queryParams.Encode())
}

func (c *Client) getCompositePolicyRulesURL(names []string) string {
	items := []string{}
	for _, name := range names {
		items = append(items, fmt.Sprintf("names=%s", name))
	}
	return fmt.Sprintf(getCompositePolicyRulesPath, c.config.url, strings.Join(items[:], "&"))
}

type getPoliciesQueryParams struct {
	PolicyType string
	Filter     string
	Limit      int
}

func (q *getPoliciesQueryParams) Encode() string {
	return fmt.Sprintf("policyType=%s&filter=%s&limit=%d", q.PolicyType, url.QueryEscape(q.Filter), q.Limit)
}

package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	createRulePath           = "%s/api/secure/rules?skipPolicyV2Msg=%t"
	getRuleByIDPath          = "%s/api/secure/rules/%d"
	updateRulePath           = "%s/api/secure/rules/%d?skipPolicyV2Msg=%t"
	deleteURLPath            = "%s/api/secure/rules/%d?skipPolicyV2Msg=%t"
	getRuleGroupPath         = "%s/api/secure/rules/groups?name=%s&type=%s"
	createStatefulRulePath   = "%s/api/policies/v3/statefulRules"
	updateStatefulRulePath   = "%s/api/policies/v3/statefulRules/%d"
	deleteStatefulRulePath   = "%s/api/policies/v3/statefulRules/%d"
	getStatefulRuleGroupPath = "%s/api/policies/v3/statefulRules/groups?name=%s&type=%s"
)

type RuleInterface interface {
	Base
	CreateRule(ctx context.Context, rule Rule) (Rule, error)
	GetRuleByID(ctx context.Context, ruleID int) (Rule, int, error)
	UpdateRule(ctx context.Context, rule Rule) (Rule, error)
	DeleteRule(ctx context.Context, ruleID int) error
	GetRuleGroup(ctx context.Context, ruleName string, ruleType string) ([]Rule, error)
	CreateStatefulRule(ctx context.Context, rule Rule) (Rule, error)
	UpdateStatefulRule(ctx context.Context, rule Rule) (Rule, error)
	DeleteStatefulRule(ctx context.Context, ruleID int) error
	GetStatefulRuleGroup(ctx context.Context, ruleName string, ruleType string) ([]Rule, error)
}

func (c *Client) CreateRule(ctx context.Context, rule Rule) (createdRule Rule, err error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createRuleURL(), payload)
	if err != nil {
		return Rule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (c *Client) GetRuleByID(ctx context.Context, ruleID int) (rule Rule, statusCode int, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getRuleByIDURL(ruleID), nil)
	if err != nil {
		return Rule{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Rule{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	rule, err = Unmarshal[Rule](response.Body)
	return rule, 0, err
}

func (c *Client) UpdateRule(ctx context.Context, rule Rule) (updatedRule Rule, err error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateRuleURL(rule.ID), payload)
	if err != nil {
		return Rule{}, err
	}

	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (c *Client) DeleteRule(ctx context.Context, ruleID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteRuleURL(ruleID), nil)
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

func (c *Client) GetRuleGroup(ctx context.Context, ruleName string, ruleType string) (rules []Rule, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getRuleGroupURL(ruleName, ruleType), nil)
	if err != nil {
		return []Rule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return []Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[[]Rule](response.Body)
}

func (c *Client) CreateStatefulRule(ctx context.Context, rule Rule) (createdRule Rule, err error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.createStatefulRuleURL(), payload)
	if err != nil {
		return Rule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (c *Client) UpdateStatefulRule(ctx context.Context, rule Rule) (updatedRule Rule, err error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateStatefulRuleURL(rule.ID), payload)
	if err != nil {
		return Rule{}, err
	}

	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (c *Client) DeleteStatefulRule(ctx context.Context, ruleID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteStatefulRuleURL(ruleID), nil)
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

func (c *Client) GetStatefulRuleGroup(ctx context.Context, ruleName string, ruleType string) (rules []Rule, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getStatefulRuleGroupURL(ruleName, ruleType), nil)
	if err != nil {
		return []Rule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return []Rule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[[]Rule](response.Body)
}

func (c *Client) createRuleURL() string {
	return fmt.Sprintf(createRulePath, c.config.url, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) getRuleByIDURL(ruleID int) string {
	return fmt.Sprintf(getRuleByIDPath, c.config.url, ruleID)
}

func (c *Client) updateRuleURL(ruleID int) string {
	return fmt.Sprintf(updateRulePath, c.config.url, ruleID, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) deleteRuleURL(ruleID int) string {
	return fmt.Sprintf(deleteURLPath, c.config.url, ruleID, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) getRuleGroupURL(ruleName string, ruleType string) string {
	return fmt.Sprintf(getRuleGroupPath, c.config.url, url.QueryEscape(ruleName), url.QueryEscape(ruleType))
}

func (c *Client) createStatefulRuleURL() string {
	return fmt.Sprintf(createStatefulRulePath, c.config.url)
}

func (c *Client) updateStatefulRuleURL(ruleID int) string {
	return fmt.Sprintf(updateStatefulRulePath, c.config.url, ruleID)
}

func (c *Client) deleteStatefulRuleURL(ruleID int) string {
	return fmt.Sprintf(deleteStatefulRulePath, c.config.url, ruleID)
}

func (c *Client) getStatefulRuleGroupURL(ruleName string, ruleType string) string {
	return fmt.Sprintf(getStatefulRuleGroupPath, c.config.url, url.QueryEscape(ruleName), url.QueryEscape(ruleType))
}

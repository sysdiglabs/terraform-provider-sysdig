package v2

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	CreateRulePath   = "%s/api/secure/rules"
	GetRuleByIDPath  = "%s/api/secure/rules/%d"
	UpdateRulePath   = "%s/api/secure/rules/%d"
	DeleteURLPath    = "%s/api/secure/rules/%d"
	GetRuleGroupPath = "%s/api/secure/rules/groups?name=%s&type=%s"
)

type RuleInterface interface {
	Base
	CreateRule(ctx context.Context, rule Rule) (Rule, error)
	GetRuleByID(ctx context.Context, ruleID int) (Rule, error, int)
	UpdateRule(ctx context.Context, rule Rule) (Rule, error)
	DeleteRule(ctx context.Context, ruleID int) error
	GetRuleGroup(ctx context.Context, ruleName string, ruleType string) ([]Rule, error)
}

func (client *Client) CreateRule(ctx context.Context, rule Rule) (Rule, error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateRuleURL(), payload)
	if err != nil {
		return Rule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (client *Client) GetRuleByID(ctx context.Context, ruleID int) (Rule, error, int) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetRuleByIDURL(ruleID), nil)
	if err != nil {
		return Rule{}, err, 0
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, client.ErrorFromResponse(response), response.StatusCode
	}

	rule, err := Unmarshal[Rule](response.Body)
	return rule, err, 0
}

func (client *Client) UpdateRule(ctx context.Context, rule Rule) (Rule, error) {
	payload, err := Marshal(rule)
	if err != nil {
		return Rule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateRuleURL(rule.ID), payload)
	if err != nil {
		return Rule{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[Rule](response.Body)
}

func (client *Client) DeleteRule(ctx context.Context, ruleID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteRuleURL(ruleID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	return err
}

func (client *Client) GetRuleGroup(ctx context.Context, ruleName string, ruleType string) ([]Rule, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetRuleGroupURL(ruleName, ruleType), nil)
	if err != nil {
		return []Rule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Rule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[[]Rule](response.Body)

}

func (client *Client) CreateRuleURL() string {
	return fmt.Sprintf(CreateRulePath, client.config.url)
}

func (client *Client) GetRuleByIDURL(ruleID int) string {
	return fmt.Sprintf(GetRuleByIDPath, client.config.url, ruleID)
}

func (client *Client) UpdateRuleURL(ruleID int) string {
	return fmt.Sprintf(UpdateRulePath, client.config.url, ruleID)
}

func (client *Client) DeleteRuleURL(ruleID int) string {
	return fmt.Sprintf(DeleteURLPath, client.config.url, ruleID)
}

func (client *Client) GetRuleGroupURL(ruleName string, ruleType string) string {
	return fmt.Sprintf(GetRuleGroupPath, client.config.url, url.QueryEscape(ruleName), url.QueryEscape(ruleType))
}

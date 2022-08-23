package secure

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (client *sysdigSecureClient) rulesURL() string {
	return fmt.Sprintf("%s/api/secure/rules", client.URL)
}

func (client *sysdigSecureClient) ruleURL(ruleID int) string {
	return fmt.Sprintf("%s/api/secure/rules/%d", client.URL, ruleID)
}

func (client *sysdigSecureClient) CreateRule(ctx context.Context, rule Rule) (result Rule, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.rulesURL(), rule.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) GetRuleByID(ctx context.Context, ruleID int) (result Rule, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.ruleURL(ruleID), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateRule(ctx context.Context, rule Rule) (result Rule, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.ruleURL(rule.ID), rule.ToJSON())
	if err != nil {
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Rule{}, errorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteRule(ctx context.Context, ruleID int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.ruleURL(ruleID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}

	return err
}

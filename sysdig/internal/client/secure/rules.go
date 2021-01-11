package secure

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
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
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode != http.StatusOK {
		return Rule{}, errors.New(string(body))
	}

	defer response.Body.Close()

	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) GetRuleByID(ctx context.Context, ruleID int) (result Rule, err error) {
	response, _ := client.doSysdigSecureRequest(ctx, http.MethodGet, client.ruleURL(ruleID), nil)
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return Rule{}, errors.New(string(body))
	}

	defer response.Body.Close()

	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) UpdateRule(ctx context.Context, rule Rule) (result Rule, err error) {
	response, _ := client.doSysdigSecureRequest(ctx, http.MethodPut, client.ruleURL(rule.ID), rule.ToJSON())
	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		return Rule{}, errors.New(string(body))
	}

	defer response.Body.Close()

	result, err = RuleFromJSON(body)
	return
}

func (client *sysdigSecureClient) DeleteRule(ctx context.Context, ruleID int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.ruleURL(ruleID), nil)

	defer response.Body.Close()

	return err
}

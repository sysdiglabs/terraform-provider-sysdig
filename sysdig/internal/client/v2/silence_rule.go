package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	silenceRulesPath = "%s/api/v1/silencingRules"
	silenceRulePath  = "%s/api/v1/silencingRules/%d"
)

var SilenceRuleNotFound = errors.New("silence rule not found")

type SilenceRuleInterface interface {
	Base
	GetSilenceRule(ctx context.Context, id int) (SilenceRule, error)
	CreateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error)
	UpdateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error)
	DeleteSilenceRule(ctx context.Context, id int) error
}

func (client *Client) GetSilenceRule(ctx context.Context, id int) (SilenceRule, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getSilenceRuleURL(id), nil)
	if err != nil {
		return SilenceRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return SilenceRule{}, SilenceRuleNotFound
	}
	if response.StatusCode != http.StatusOK {
		return SilenceRule{}, client.ErrorFromResponse(response)
	}

	silenceRule, err := Unmarshal[SilenceRule](response.Body)
	if err != nil {
		return SilenceRule{}, err
	}

	return silenceRule, nil
}

func (client *Client) CreateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error) {
	payload, err := Marshal(silenceRule)
	if err != nil {
		return SilenceRule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.getSilenceRulesURL(), payload)
	if err != nil {
		return SilenceRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return SilenceRule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[SilenceRule](response.Body)
}

func (client *Client) UpdateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error) {
	payload, err := Marshal(silenceRule)
	if err != nil {
		return SilenceRule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.getSilenceRuleURL(silenceRule.ID), payload)
	if err != nil {
		return SilenceRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return SilenceRule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[SilenceRule](response.Body)
}

func (client *Client) DeleteSilenceRule(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getSilenceRuleURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getSilenceRulesURL() string {
	return fmt.Sprintf(silenceRulesPath, client.config.url)
}

func (client *Client) getSilenceRuleURL(id int) string {
	return fmt.Sprintf(silenceRulePath, client.config.url, id)
}

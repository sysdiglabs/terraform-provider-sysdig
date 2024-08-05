package v2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	inhibitionRulesPath = "%s/monitor/alerts/v1/inhibition-rules"
	inhibitionRulePath  = "%s/monitor/alerts/v1/inhibition-rules/%d"
)

var ErrInhibitionRuleNotFound = errors.New("inhibition rule not found")

type InhibitionRuleInterface interface {
	Base
	GetInhibitionRule(ctx context.Context, id int) (InhibitionRule, error)
	CreateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error)
	UpdateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error)
	DeleteInhibitionRule(ctx context.Context, id int) error
}

func (client *Client) GetInhibitionRule(ctx context.Context, id int) (InhibitionRule, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.getInhibitionRuleURL(id), nil)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return InhibitionRule{}, ErrInhibitionRuleNotFound
	}
	if response.StatusCode != http.StatusOK {
		return InhibitionRule{}, client.ErrorFromResponse(response)
	}

	inhibitionRule, err := Unmarshal[InhibitionRule](response.Body)
	if err != nil {
		return InhibitionRule{}, err
	}

	return inhibitionRule, nil
}

func (client *Client) CreateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error) {
	payload, err := Marshal(inhibitionRule)
	if err != nil {
		return InhibitionRule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.getInhibitionRulesURL(), payload)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return InhibitionRule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[InhibitionRule](response.Body)
}

func (client *Client) UpdateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error) {
	payload, err := Marshal(inhibitionRule)
	if err != nil {
		return InhibitionRule{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.getInhibitionRuleURL(inhibitionRule.ID), payload)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return InhibitionRule{}, client.ErrorFromResponse(response)
	}

	return Unmarshal[InhibitionRule](response.Body)
}

func (client *Client) DeleteInhibitionRule(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.getInhibitionRuleURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) getInhibitionRulesURL() string {
	return fmt.Sprintf(inhibitionRulesPath, client.config.url)
}

func (client *Client) getInhibitionRuleURL(id int) string {
	return fmt.Sprintf(inhibitionRulePath, client.config.url, id)
}

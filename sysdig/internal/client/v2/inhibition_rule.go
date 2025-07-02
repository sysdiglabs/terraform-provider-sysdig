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
	GetInhibitionRuleByID(ctx context.Context, id int) (InhibitionRule, error)
	CreateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error)
	UpdateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (InhibitionRule, error)
	DeleteInhibitionRule(ctx context.Context, id int) error
}

func (c *Client) GetInhibitionRuleByID(ctx context.Context, id int) (rule InhibitionRule, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getInhibitionRuleURL(id), nil)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return InhibitionRule{}, ErrInhibitionRuleNotFound
	}
	if response.StatusCode != http.StatusOK {
		return InhibitionRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[InhibitionRule](response.Body)
}

func (c *Client) CreateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (rule InhibitionRule, err error) {
	payload, err := Marshal(inhibitionRule)
	if err != nil {
		return InhibitionRule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getInhibitionRulesURL(), payload)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return InhibitionRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[InhibitionRule](response.Body)
}

func (c *Client) UpdateInhibitionRule(ctx context.Context, inhibitionRule InhibitionRule) (rule InhibitionRule, err error) {
	payload, err := Marshal(inhibitionRule)
	if err != nil {
		return InhibitionRule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getInhibitionRuleURL(inhibitionRule.ID), payload)
	if err != nil {
		return InhibitionRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return InhibitionRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[InhibitionRule](response.Body)
}

func (c *Client) DeleteInhibitionRule(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getInhibitionRuleURL(id), nil)
	if err != nil {
		return err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return c.ErrorFromResponse(response)
	}

	return nil
}

func (c *Client) getInhibitionRulesURL() string {
	return fmt.Sprintf(inhibitionRulesPath, c.config.url)
}

func (c *Client) getInhibitionRuleURL(id int) string {
	return fmt.Sprintf(inhibitionRulePath, c.config.url, id)
}

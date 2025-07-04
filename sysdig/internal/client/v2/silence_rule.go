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

var ErrSilenceRuleNotFound = errors.New("silence rule not found")

type SilenceRuleInterface interface {
	Base
	GetSilenceRule(ctx context.Context, id int) (SilenceRule, error)
	CreateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error)
	UpdateSilenceRule(ctx context.Context, silenceRule SilenceRule) (SilenceRule, error)
	DeleteSilenceRule(ctx context.Context, id int) error
}

func (c *Client) GetSilenceRule(ctx context.Context, id int) (rule SilenceRule, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getSilenceRuleURL(id), nil)
	if err != nil {
		return SilenceRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode == http.StatusNotFound {
		return SilenceRule{}, ErrSilenceRuleNotFound
	}
	if response.StatusCode != http.StatusOK {
		return SilenceRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[SilenceRule](response.Body)
}

func (c *Client) CreateSilenceRule(ctx context.Context, silenceRule SilenceRule) (rule SilenceRule, err error) {
	payload, err := Marshal(silenceRule)
	if err != nil {
		return SilenceRule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.getSilenceRulesURL(), payload)
	if err != nil {
		return SilenceRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return SilenceRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[SilenceRule](response.Body)
}

func (c *Client) UpdateSilenceRule(ctx context.Context, silenceRule SilenceRule) (rule SilenceRule, err error) {
	payload, err := Marshal(silenceRule)
	if err != nil {
		return SilenceRule{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.getSilenceRuleURL(silenceRule.ID), payload)
	if err != nil {
		return SilenceRule{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return SilenceRule{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[SilenceRule](response.Body)
}

func (c *Client) DeleteSilenceRule(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.getSilenceRuleURL(id), nil)
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

func (c *Client) getSilenceRulesURL() string {
	return fmt.Sprintf(silenceRulesPath, c.config.url)
}

func (c *Client) getSilenceRuleURL(id int) string {
	return fmt.Sprintf(silenceRulePath, c.config.url, id)
}

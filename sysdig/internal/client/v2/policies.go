package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	createPolicyPath         = "%s/api/v2/policies?skipPolicyV2Msg=%t"
	deletePolicyPath         = "%s/api/v2/policies/%d?skipPolicyV2Msg=%t"
	updatePolicyPath         = "%s/api/v2/policies/%d?skipPolicyV2Msg=%t"
	getPolicyPath            = "%s/api/v2/policies/%d"
	getPoliciesPath          = "%s/api/v2/policies"
	sendPoliciesToAgentsPath = "%s/api/v2/policies/actions?action=forwardPolicyV2Msg"
)

type PolicyInterface interface {
	Base
	CreatePolicy(ctx context.Context, policy Policy) (Policy, error)
	DeletePolicy(ctx context.Context, policyID int) error
	UpdatePolicy(ctx context.Context, policy Policy) (Policy, error)
	GetPolicyByID(ctx context.Context, policyID int) (Policy, int, error)
	GetPolicies(ctx context.Context) ([]Policy, int, error)
	SendPoliciesToAgents(ctx context.Context) error
}

func (c *Client) CreatePolicy(ctx context.Context, policy Policy) (createdPolicy Policy, err error) {
	payload, err := Marshal(policy)
	if err != nil {
		return Policy{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createPolicyURL(), payload)
	if err != nil {
		return Policy{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Policy{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Policy](response.Body)
}

func (c *Client) DeletePolicy(ctx context.Context, policyID int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deletePolicyURL(policyID), nil)
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

func (c *Client) UpdatePolicy(ctx context.Context, policy Policy) (updatedPolicy Policy, err error) {
	payload, err := Marshal(policy)
	if err != nil {
		return Policy{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updatePolicyURL(policy.ID), payload)
	if err != nil {
		return Policy{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Policy{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Policy](response.Body)
}

func (c *Client) GetPolicyByID(ctx context.Context, policyID int) (policy Policy, statusCode int, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getPolicyURL(policyID), nil)
	if err != nil {
		return Policy{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Policy{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	policy, err = Unmarshal[Policy](response.Body)
	if err != nil {
		return Policy{}, 0, err
	}

	return policy, http.StatusOK, nil
}

func (c *Client) GetPolicies(ctx context.Context) (policies []Policy, statusCode int, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getPoliciesURL(), nil)
	if err != nil {
		return []Policy{}, 0, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return []Policy{}, response.StatusCode, c.ErrorFromResponse(response)
	}

	policies, err = Unmarshal[[]Policy](response.Body)
	if err != nil {
		return []Policy{}, 0, err
	}

	return policies, http.StatusOK, nil
}

func (c *Client) SendPoliciesToAgents(ctx context.Context) (err error) {
	if c.config.secureSkipPolicyV2Msg {
		// We only need to send policies if we've been configured to skip sending them during updates
		response, err := c.requester.Request(ctx, http.MethodPost, c.sendPoliciesToAgentsURL(), nil)
		if err != nil {
			return err
		}
		defer func() {
			if dErr := response.Body.Close(); dErr != nil {
				err = fmt.Errorf("unable to close response body: %w", dErr)
			}
		}()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected response when sending policies to agents: %s", response.Status)
		}
	}
	return nil
}

func (c *Client) createPolicyURL() string {
	return fmt.Sprintf(createPolicyPath, c.config.url, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) deletePolicyURL(policyID int) string {
	return fmt.Sprintf(deletePolicyPath, c.config.url, policyID, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) updatePolicyURL(policyID int) string {
	return fmt.Sprintf(updatePolicyPath, c.config.url, policyID, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) getPolicyURL(policyID int) string {
	return fmt.Sprintf(getPolicyPath, c.config.url, policyID)
}

func (c *Client) getPoliciesURL() string {
	return fmt.Sprintf(getPoliciesPath, c.config.url)
}

func (c *Client) sendPoliciesToAgentsURL() string {
	return fmt.Sprintf(sendPoliciesToAgentsPath, c.config.url)
}

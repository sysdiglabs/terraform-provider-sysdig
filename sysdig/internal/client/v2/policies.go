package v2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	CreatePolicyPath         = "%s/api/v2/policies?skipPolicyV2Msg=%t"
	DeletePolicyPath         = "%s/api/v2/policies/%d?skipPolicyV2Msg=%t"
	UpdatePolicyPath         = "%s/api/v2/policies/%d?skipPolicyV2Msg=%t"
	GetPolicyPath            = "%s/api/v2/policies/%d"
	GetPoliciesPath          = "%s/api/v2/policies"
	SendPoliciesToAgentsPath = "%s/api/v2/policies/actions?action=forwardPolicyV2Msg"
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

func (client *Client) CreatePolicy(ctx context.Context, policy Policy) (Policy, error) {
	payload, err := Marshal(policy)
	if err != nil {
		return Policy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreatePolicyURL(), payload)
	if err != nil {
		return Policy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Policy{}, client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return Unmarshal[Policy](response.Body)
}

func (client *Client) DeletePolicy(ctx context.Context, policyID int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeletePolicyURL(policyID), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return err
}

func (client *Client) UpdatePolicy(ctx context.Context, policy Policy) (Policy, error) {
	payload, err := Marshal(policy)
	if err != nil {
		return Policy{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdatePolicyURL(policy.ID), payload)
	if err != nil {
		return Policy{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Policy{}, client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return Unmarshal[Policy](response.Body)
}

func (client *Client) GetPolicyByID(ctx context.Context, policyID int) (Policy, int, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetPolicyURL(policyID), nil)
	if err != nil {
		return Policy{}, 0, err

	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Policy{}, response.StatusCode, client.ErrorFromResponse(response)
	}

	policy, err := Unmarshal[Policy](response.Body)
	if err != nil {
		return Policy{}, 0, err
	}

	return policy, http.StatusOK, nil
}

func (client *Client) GetPolicies(ctx context.Context) ([]Policy, int, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetPoliciesURL(), nil)
	if err != nil {
		return []Policy{}, 0, err

	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return []Policy{}, response.StatusCode, client.ErrorFromResponse(response)
	}

	policies, err := Unmarshal[[]Policy](response.Body)
	if err != nil {
		return []Policy{}, 0, err
	}

	return policies, http.StatusOK, nil
}

func (client *Client) SendPoliciesToAgents(ctx context.Context) error {
	if client.config.secureSkipPolicyV2Msg && client.policiesChanged {
		// If we have been skipping sending the policy v2 message and there have been changes
		// we need to tell the Policies API to go ahead and send the policies to the agents now.

		tflog.Warn(ctx, "Policies have changed - Sending policies to agents")
		response, err := client.requester.Request(ctx, http.MethodPost, client.SendPoliciesToAgentsURL(), nil)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error while sending policies to agents: %s", err.Error()))
			return err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			tflog.Warn(ctx, fmt.Sprintf("Unexpected response when sending policies to agents: %s", response.Status))
		}
	}

	return nil
}

func (client *Client) CreatePolicyURL() string {
	return fmt.Sprintf(CreatePolicyPath, client.config.url, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) DeletePolicyURL(policyID int) string {
	return fmt.Sprintf(DeletePolicyPath, client.config.url, policyID, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) UpdatePolicyURL(policyID int) string {
	return fmt.Sprintf(UpdatePolicyPath, client.config.url, policyID, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) GetPolicyURL(policyID int) string {
	return fmt.Sprintf(GetPolicyPath, client.config.url, policyID)
}

func (client *Client) GetPoliciesURL() string {
	return fmt.Sprintf(GetPoliciesPath, client.config.url)
}

func (client *Client) SendPoliciesToAgentsURL() string {
	return fmt.Sprintf(SendPoliciesToAgentsPath, client.config.url)
}

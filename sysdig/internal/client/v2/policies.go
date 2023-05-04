package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	CreatePolicyPath = "%s/api/v2/policies"
	DeletePolicyPath = "%s/api/v2/policies/%d"
	UpdatePolicyPath = "%s/api/v2/policies/%d"
	GetPolicyPath    = "%s/api/v2/policies/%d"
)

type PolicyInterface interface {
	CreatePolicy(ctx context.Context, policy Policy) (Policy, error)
	DeletePolicy(ctx context.Context, policyID int) error
	UpdatePolicy(ctx context.Context, policy Policy) (Policy, error)
	GetPolicyByID(ctx context.Context, policyID int) (Policy, int, error)
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

func (client *Client) CreatePolicyURL() string {
	return fmt.Sprintf(CreatePolicyPath, client.config.url)
}

func (client *Client) DeletePolicyURL(policyID int) string {
	return fmt.Sprintf(DeletePolicyPath, client.config.url, policyID)
}

func (client *Client) UpdatePolicyURL(policyID int) string {
	return fmt.Sprintf(UpdatePolicyPath, client.config.url, policyID)
}

func (client *Client) GetPolicyURL(policyID int) string {
	return fmt.Sprintf(GetPolicyPath, client.config.url, policyID)
}

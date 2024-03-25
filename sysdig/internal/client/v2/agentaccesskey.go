package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetAgentAccessKeyByIdPath = "%s/platform/v1/access-keys/%s"
	CreateAgentAccessKeyPath  = "%s/platform/v1/access-keys"
	DeleteAgentAccessKeyPath  = "%s/platform/v1/access-keys/%s"
	PutAgentAccessKeyPath     = "%s/platform/v1/access-keys/%s"
)

type AgentAccessKeyInterface interface {
	Base
	GetAgentAccessKeyById(ctx context.Context, id string) (*AgentAccessKey, error)
	CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error)
	DeleteAgentAccessKey(ctx context.Context, id string) error
	UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey, id string) (*AgentAccessKey, error)
}

func (client *Client) GetAgentAccessKeyById(ctx context.Context, id string) (*AgentAccessKey, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetAgentAccessKeyByIdUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	agentAccessKey, err := Unmarshal[AgentAccessKey](response.Body)
	if err != nil {
		return nil, err
	}

	return &agentAccessKey, nil
}

func (client *Client) CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error) {
	payload, err := Marshal(agentAccessKey)
	if err != nil {
		return nil, err
	}
	response, err := client.requester.Request(ctx, http.MethodPost, client.PostAgentAccessKeyUrl(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	createdAgentAccessKey, err := Unmarshal[AgentAccessKey](response.Body)

	if err != nil {
		return nil, err
	}

	return &createdAgentAccessKey, nil
}

func (client *Client) UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey, id string) (*AgentAccessKey, error) {

	payload, err := Marshal(agentAccessKey)
	if err != nil {
		return nil, err
	}
	response, err := client.requester.Request(ctx, http.MethodPut, client.PutAgentAccessKeyUrl(id), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	updatedAgentAccessKey, err := Unmarshal[AgentAccessKey](response.Body)
	if err != nil {
		return nil, err
	}

	return &updatedAgentAccessKey, nil
}

func (client *Client) DeleteAgentAccessKey(ctx context.Context, id string) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteAgentAccessKeyUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	return nil
}

func (client *Client) GetAgentAccessKeyByIdUrl(id string) string {
	return fmt.Sprintf(GetAgentAccessKeyByIdPath, client.config.url, id)
}

func (client *Client) PostAgentAccessKeyUrl() string {
	return fmt.Sprintf(CreateAgentAccessKeyPath, client.config.url)
}

func (client *Client) PutAgentAccessKeyUrl(id string) string {
	return fmt.Sprintf(PutAgentAccessKeyPath, client.config.url, id)
}

func (client *Client) DeleteAgentAccessKeyUrl(id string) string {
	return fmt.Sprintf(DeleteAgentAccessKeyPath, client.config.url, id)
}

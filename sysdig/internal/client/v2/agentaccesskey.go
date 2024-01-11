package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetAgentAccessKeyPath = "%s/api/customer/accessKeys"
	//CreateAgentAccessKeyPath = "%s/api/customer/accessKeys"
)

type AgentAccessKeyInterface interface {
	Base
	GetAgentAccessKeyById(ctx context.Context, id string) (*AgentAccessKey, error)
	//CreateAgentAccessKey(ctx context.Context, user *AgentAccessKey) (*AgentAccessKey, error)
}

func (client *Client) GetAgentAccessKeyById(ctx context.Context, id string) (*AgentAccessKey, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetAgentAccessKeyUrl(), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	wrapper, err := Unmarshal[AgentAccessKeyReadWrapper](response.Body)
	if err != nil {
		return nil, err
	}
	for _, key := range wrapper.CustomerAccessKey {
		if key.AgentAccessKeyId == id {
			return &key, nil // Found the key, return it
		}
	}

	return nil, fmt.Errorf("no AgentAccessKey found with ID %s", id)
}

// Function to find an AgentAccessKey by ID
func findAgentAccessKeyByID(keys []AgentAccessKey, accessKey string) (*AgentAccessKey, bool) {
	for _, key := range keys {
		if key.AgentAccessKeyId == accessKey {
			return &key, true
		}
	}
	return nil, false
}

func (client *Client) GetAgentAccessKeyUrl() string {
	return fmt.Sprintf(GetAgentAccessKeyPath, client.config.url)
}

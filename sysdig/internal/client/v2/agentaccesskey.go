package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	GetAgentAccessKeyPath     = "%s/api/customer/accessKeys"
	CreateAgentAccessKeyPath  = "%s/api/customer/accessKeys"
	DeleteAgentAccessKeyPath  = "%s/api/customer/accessKeys/%s"
	DisableAgentAccessKeyPath = "%s/api/customer/accessKeys/%s/disable"
	EnableAgentAccessKeyPath  = "%s/api/customer/accessKeys/%s/enable"
	PutAgentAccessKeyPath     = "%s/api/customer/accessKeys/%s"
)

type AgentAccessKeyInterface interface {
	Base
	GetAgentAccessKeyById(ctx context.Context, id string) (*AgentAccessKey, error)
	CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error)
	DeleteAgentAccessKey(ctx context.Context, id string) error
	UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error)
	EnableOrDisableAgentAccessKey(ctx context.Context, id string, enable bool) error
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
	fmt.Println("Trying to get agent access keys with id: ", id)
	return nil, fmt.Errorf("no AgentAccessKey found with ID %s", id)
}

func (client *Client) CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error) {
	agentAccessKeyWriteWrapper := AgentAccessKeyWriteWrapper{CustomerAccessKey: *agentAccessKey}
	payload, err := Marshal(agentAccessKeyWriteWrapper)
	response, err := client.requester.Request(ctx, http.MethodPost, client.PostAgentAccessKeyUrl(), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	created, err := Unmarshal[AgentAccessKeyWriteWrapper](response.Body)

	if err != nil {
		return nil, err
	}

	return &created.CustomerAccessKey, nil
}

func (client *Client) UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error) {
	agentAccessKeyWriteWrapper := AgentAccessKeyWriteWrapper{CustomerAccessKey: *agentAccessKey}
	fmt.Println("agent config: ", agentAccessKeyWriteWrapper.CustomerAccessKey)
	agentAccessKeyId := agentAccessKeyWriteWrapper.CustomerAccessKey.AgentAccessKeyId
	fmt.Println("ID: ", agentAccessKeyId)
	payload, err := Marshal(agentAccessKeyWriteWrapper)
	response, err := client.requester.Request(ctx, http.MethodPut, client.PutAgentAccessKeyUrl(agentAccessKeyId), payload)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return nil, err
	}

	updated, err := Unmarshal[AgentAccessKeyWriteWrapper](response.Body)
	if err != nil {
		return nil, err
	}

	return &updated.CustomerAccessKey, nil
}

func (client *Client) EnableOrDisableAgentAccessKey(ctx context.Context, id string, enable bool) error {
	var url string
	if enable == true {
		url = client.EnableAgentAccessKeyUrl(id)
	} else {
		url = client.DisableAgentAccessKeyUrl(id)
	}
	response, err := client.requester.Request(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = client.ErrorFromResponse(response)
		return err
	}

	_, err = Unmarshal[AgentAccessKeyWriteWrapper](response.Body)
	if err != nil {
		return err
	}

	return nil
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

func (client *Client) GetAgentAccessKeyUrl() string {
	return fmt.Sprintf(GetAgentAccessKeyPath, client.config.url)
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

func (client *Client) DisableAgentAccessKeyUrl(id string) string {
	return fmt.Sprintf(DisableAgentAccessKeyPath, client.config.url, id)
}

func (client *Client) EnableAgentAccessKeyUrl(id string) string {
	return fmt.Sprintf(EnableAgentAccessKeyPath, client.config.url, id)
}

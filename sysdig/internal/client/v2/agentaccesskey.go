package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	getAgentAccessKeyByIDPath = "%s/platform/v1/access-keys/%s"
	createAgentAccessKeyPath  = "%s/platform/v1/access-keys"
	deleteAgentAccessKeyPath  = "%s/platform/v1/access-keys/%s"
	putAgentAccessKeyPath     = "%s/platform/v1/access-keys/%s"
)

type AgentAccessKeyInterface interface {
	Base
	GetAgentAccessKeyByID(ctx context.Context, id string) (*AgentAccessKey, error)
	CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (*AgentAccessKey, error)
	DeleteAgentAccessKey(ctx context.Context, id string) error
	UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey, id string) (*AgentAccessKey, error)
}

func (c *Client) GetAgentAccessKeyByID(ctx context.Context, id string) (accessKey *AgentAccessKey, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getAgentAccessKeyByIDUrl(id), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*AgentAccessKey](response.Body)
}

func (c *Client) CreateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey) (createdAccessKey *AgentAccessKey, err error) {
	payload, err := Marshal(agentAccessKey)
	if err != nil {
		return nil, err
	}
	response, err := c.requester.Request(ctx, http.MethodPost, c.postAgentAccessKeyURL(), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusCreated {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*AgentAccessKey](response.Body)
}

func (c *Client) UpdateAgentAccessKey(ctx context.Context, agentAccessKey *AgentAccessKey, id string) (updatedAccessKey *AgentAccessKey, err error) {
	payload, err := Marshal(agentAccessKey)
	if err != nil {
		return nil, err
	}
	response, err := c.requester.Request(ctx, http.MethodPut, c.putAgentAccessKeyURL(id), payload)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		err = c.ErrorFromResponse(response)
		return nil, err
	}

	return Unmarshal[*AgentAccessKey](response.Body)
}

func (c *Client) DeleteAgentAccessKey(ctx context.Context, id string) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.DeleteAgentAccessKeyURL(id), nil)
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

func (c *Client) getAgentAccessKeyByIDUrl(id string) string {
	return fmt.Sprintf(getAgentAccessKeyByIDPath, c.config.url, id)
}

func (c *Client) postAgentAccessKeyURL() string {
	return fmt.Sprintf(createAgentAccessKeyPath, c.config.url)
}

func (c *Client) putAgentAccessKeyURL(id string) string {
	return fmt.Sprintf(putAgentAccessKeyPath, c.config.url, id)
}

func (c *Client) DeleteAgentAccessKeyURL(id string) string {
	return fmt.Sprintf(deleteAgentAccessKeyPath, c.config.url, id)
}

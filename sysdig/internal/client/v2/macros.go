package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	CreateMacroPath  = "%s/api/secure/falco/macros?skipPolicyV2Msg=%t"
	GetMacroByIDPath = "%s/api/secure/falco/macros/%d"
	UpdateMacroPath  = "%s/api/secure/falco/macros/%d?skipPolicyV2Msg=%t"
	DeleteMacroPath  = "%s/api/secure/falco/macros/%d?skipPolicyV2Msg=%t"
)

type MacroInterface interface {
	Base
	CreateMacro(ctx context.Context, macro Macro) (Macro, error)
	GetMacroByID(ctx context.Context, id int) (Macro, error)
	UpdateMacro(ctx context.Context, macro Macro) (Macro, error)
	DeleteMacro(ctx context.Context, id int) error
}

func (client *Client) CreateMacro(ctx context.Context, macro Macro) (Macro, error) {
	payload, err := Marshal(macro)
	if err != nil {
		return Macro{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPost, client.CreateMacroURL(), payload)
	if err != nil {
		return Macro{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return Macro{}, client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return Unmarshal[Macro](response.Body)
}

func (client *Client) GetMacroByID(ctx context.Context, id int) (Macro, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, client.GetMacroByIDURL(id), nil)
	if err != nil {
		return Macro{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Macro{}, client.ErrorFromResponse(response)
	}

	macro, err := Unmarshal[Macro](response.Body)
	if err != nil {
		return Macro{}, err
	}

	if macro.Version == 0 {
		return Macro{}, fmt.Errorf("macro with ID: %d does not exists", id)
	}

	return macro, nil
}

func (client *Client) UpdateMacro(ctx context.Context, macro Macro) (Macro, error) {
	payload, err := Marshal(macro)
	if err != nil {
		return Macro{}, err
	}

	response, err := client.requester.Request(ctx, http.MethodPut, client.UpdateMacroURL(macro.ID), payload)
	if err != nil {
		return Macro{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Macro{}, client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return Unmarshal[Macro](response.Body)
}

func (client *Client) DeleteMacro(ctx context.Context, id int) error {
	response, err := client.requester.Request(ctx, http.MethodDelete, client.DeleteMacroURL(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNotFound {
		return client.ErrorFromResponse(response)
	}

	client.policiesChanged = true

	return nil
}

func (client *Client) CreateMacroURL() string {
	return fmt.Sprintf(CreateMacroPath, client.config.url, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) GetMacroByIDURL(id int) string {
	return fmt.Sprintf(GetMacroByIDPath, client.config.url, id)
}

func (client *Client) UpdateMacroURL(id int) string {
	return fmt.Sprintf(UpdateMacroPath, client.config.url, id, client.config.secureSkipPolicyV2Msg)
}

func (client *Client) DeleteMacroURL(id int) string {
	return fmt.Sprintf(DeleteMacroPath, client.config.url, id, client.config.secureSkipPolicyV2Msg)
}

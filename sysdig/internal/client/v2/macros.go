package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	createMacroPath  = "%s/api/secure/falco/macros?skipPolicyV2Msg=%t"
	getMacroByIDPath = "%s/api/secure/falco/macros/%d"
	updateMacroPath  = "%s/api/secure/falco/macros/%d?skipPolicyV2Msg=%t"
	deleteMacroPath  = "%s/api/secure/falco/macros/%d?skipPolicyV2Msg=%t"
)

type MacroInterface interface {
	Base
	CreateMacro(ctx context.Context, macro Macro) (Macro, error)
	GetMacroByID(ctx context.Context, id int) (Macro, error)
	UpdateMacro(ctx context.Context, macro Macro) (Macro, error)
	DeleteMacro(ctx context.Context, id int) error
}

func (c *Client) CreateMacro(ctx context.Context, macro Macro) (createdMacro Macro, err error) {
	payload, err := Marshal(macro)
	if err != nil {
		return Macro{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPost, c.createMacroURL(), payload)
	if err != nil {
		return Macro{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return Macro{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Macro](response.Body)
}

func (c *Client) GetMacroByID(ctx context.Context, id int) (macro Macro, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, c.getMacroByIDURL(id), nil)
	if err != nil {
		return Macro{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Macro{}, c.ErrorFromResponse(response)
	}

	macro, err = Unmarshal[Macro](response.Body)
	if err != nil {
		return Macro{}, err
	}

	if macro.Version == 0 {
		return Macro{}, fmt.Errorf("macro with ID: %d does not exists", id)
	}

	return macro, nil
}

func (c *Client) UpdateMacro(ctx context.Context, macro Macro) (updateMacro Macro, err error) {
	payload, err := Marshal(macro)
	if err != nil {
		return Macro{}, err
	}

	response, err := c.requester.Request(ctx, http.MethodPut, c.updateMacroURL(macro.ID), payload)
	if err != nil {
		return Macro{}, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return Macro{}, c.ErrorFromResponse(response)
	}

	return Unmarshal[Macro](response.Body)
}

func (c *Client) DeleteMacro(ctx context.Context, id int) (err error) {
	response, err := c.requester.Request(ctx, http.MethodDelete, c.deleteMacroURL(id), nil)
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

func (c *Client) createMacroURL() string {
	return fmt.Sprintf(createMacroPath, c.config.url, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) getMacroByIDURL(id int) string {
	return fmt.Sprintf(getMacroByIDPath, c.config.url, id)
}

func (c *Client) updateMacroURL(id int) string {
	return fmt.Sprintf(updateMacroPath, c.config.url, id, c.config.secureSkipPolicyV2Msg)
}

func (c *Client) deleteMacroURL(id int) string {
	return fmt.Sprintf(deleteMacroPath, c.config.url, id, c.config.secureSkipPolicyV2Msg)
}

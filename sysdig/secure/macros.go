package secure

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) CreateMacro(ctx context.Context, macroRequest Macro) (macro Macro, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.GetMacrosUrl(), macroRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errors.New(response.Status)
		return
	}

	macro, err = MacroFromJSON(body)
	return
}

func (client *sysdigSecureClient) GetMacroById(ctx context.Context, id int) (macro Macro, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.GetMacroUrl(id), nil)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	macro, err = MacroFromJSON(body)
	if err != nil {
		return
	}

	if macro.Version == 0 {
		err = fmt.Errorf("Macro with ID: %d does not exists", id)
		return
	}
	return
}

func (client *sysdigSecureClient) UpdateMacro(ctx context.Context, macroRequest Macro) (macro Macro, err error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.GetMacroUrl(macroRequest.ID), macroRequest.ToJSON())
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK {
		err = errors.New(response.Status)
		return
	}

	return MacroFromJSON(body)
}

func (client *sysdigSecureClient) DeleteMacro(ctx context.Context, id int) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.GetMacroUrl(id), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}
	return nil
}

func (client *sysdigSecureClient) GetMacrosUrl() string {
	return fmt.Sprintf("%s/api/secure/falco/macros", client.URL)
}

func (client *sysdigSecureClient) GetMacroUrl(id int) string {
	return fmt.Sprintf("%s/api/secure/falco/macros/%d", client.URL, id)
}

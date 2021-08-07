package secure

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (client *sysdigSecureClient) cloudAccountURL(includeExternalID bool) string {
	if includeExternalID {
		return fmt.Sprintf("%s/api/cloud/v2/accounts?includeExternalID=true", client.URL)
	}
	return fmt.Sprintf("%s/api/cloud/v2/accounts", client.URL)
}

func (client *sysdigSecureClient) cloudAccountByIdURL(accountID string, includeExternalID bool) string {
	if includeExternalID {
		return fmt.Sprintf("%s/api/cloud/v2/accounts/%s?includeExternalID=true", client.URL, accountID)
	}
	return fmt.Sprintf("%s/api/cloud/v2/accounts/%s", client.URL, accountID)
}

func (client *sysdigSecureClient) trustedUserURL() string {
	return fmt.Sprintf("%s/api/cloud/v2/trustedUser", client.URL)
}

func (client *sysdigSecureClient) CreateCloudAccount(ctx context.Context, cloudAccount *CloudAccount) (*CloudAccount, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPost, client.cloudAccountURL(true), cloudAccount.ToJSON())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		err = errorFromResponse(response)
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return CloudAccountFromJSON(bodyBytes), nil
}

func (client *sysdigSecureClient) GetCloudAccountById(ctx context.Context, accountID string) (*CloudAccount, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.cloudAccountByIdURL(accountID, true), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errorFromResponse(response)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return CloudAccountFromJSON(bodyBytes), nil
}

func (client *sysdigSecureClient) DeleteCloudAccount(ctx context.Context, accountID string) error {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodDelete, client.cloudAccountByIdURL(accountID, false), nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent && response.StatusCode != http.StatusOK {
		return errorFromResponse(response)
	}
	return nil
}

func (client *sysdigSecureClient) UpdateCloudAccount(ctx context.Context, accountID string, cloudAccount *CloudAccount) (*CloudAccount, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodPut, client.cloudAccountByIdURL(accountID, true), cloudAccount.ToJSON())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errorFromResponse(response)
		return nil, err
	}

	bodyBytes, _ := ioutil.ReadAll(response.Body)
	return CloudAccountFromJSON(bodyBytes), nil
}

func (client *sysdigSecureClient) GetTrustedCloudUser(ctx context.Context) (string, error) {
	response, err := client.doSysdigSecureRequest(ctx, http.MethodGet, client.trustedUserURL(), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", errorFromResponse(response)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

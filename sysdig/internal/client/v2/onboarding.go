package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	onboardingTrustedIdentityPath         = "%s/api/secure/onboarding/v2/trustedIdentity?provider=%s"
	onboardingTenantExternaIDPath         = "%s/api/secure/onboarding/v2/externalID"
	onboardingAgentlessScanningAssetsPath = "%s/api/secure/onboarding/v2/agentlessScanningAssets"
)

type OnboardingSecureInterface interface {
	Base
	GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error)
	GetTenantExternalIDSecure(ctx context.Context) (string, error)
	GetAgentlessScanningAssetsSecure(ctx context.Context) (map[string]interface{}, error)
}

func (client *Client) GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedIdentityPath, client.config.url, provider), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", client.ErrorFromResponse(response)
	}

	return Unmarshal[string](response.Body)
}

func (client *Client) GetTenantExternalIDSecure(ctx context.Context) (string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTenantExternaIDPath, client.config.url), nil)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", client.ErrorFromResponse(response)
	}

	return Unmarshal[string](response.Body)
}

func (client *Client) GetAgentlessScanningAssetsSecure(ctx context.Context) (map[string]interface{}, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingAgentlessScanningAssetsPath, client.config.url), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]interface{}](response.Body)
}

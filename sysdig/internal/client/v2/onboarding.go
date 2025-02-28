package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	onboardingTrustedIdentityPath         = "%s/api/secure/onboarding/v2/trustedIdentity?provider=%s"
	onboardingTrustedAzureAppPath         = "%s/api/secure/onboarding/v2/trustedAzureApp?app=%s"
	onboardingTenantExternaIDPath         = "%s/api/secure/onboarding/v2/externalID"
	onboardingAgentlessScanningAssetsPath = "%s/api/secure/onboarding/v2/agentlessScanningAssets"
	onboardingCloudIngestionAssetsPath    = "%s/api/secure/onboarding/v2/cloudIngestionAssets?provider=%s&providerID=%s&componentType=%s"
	onboardingTrustedRegulationAssetsPath = "%s/api/secure/onboarding/v2/trustedRegulationAssets?provider=%s"
	onboardingTrustedOracleAppPath        = "%s/api/secure/onboarding/v2/trustedOracleApp?app=%s"
)

type OnboardingSecureInterface interface {
	Base
	GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error)
	GetTrustedAzureAppSecure(ctx context.Context, app string) (map[string]string, error)
	GetTenantExternalIDSecure(ctx context.Context) (string, error)
	GetAgentlessScanningAssetsSecure(ctx context.Context) (map[string]any, error)
	GetCloudIngestionAssetsSecure(ctx context.Context, provider, providerID, componentType string) (map[string]any, error)
	GetTrustedCloudRegulationAssetsSecure(ctx context.Context, provider string) (map[string]string, error)
	GetTrustedOracleAppSecure(ctx context.Context, app string) (map[string]string, error)
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

func (client *Client) GetTrustedAzureAppSecure(ctx context.Context, app string) (map[string]string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedAzureAppPath, client.config.url, app), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
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

func (client *Client) GetCloudIngestionAssetsSecure(ctx context.Context, provider, providerID, componentType string) (map[string]interface{}, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingCloudIngestionAssetsPath, client.config.url, provider, providerID, componentType), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]interface{}](response.Body)
}

func (client *Client) GetTrustedCloudRegulationAssetsSecure(ctx context.Context, provider string) (map[string]string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedRegulationAssetsPath, client.config.url, provider), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
}

func (client *Client) GetTrustedOracleAppSecure(ctx context.Context, app string) (map[string]string, error) {
	response, err := client.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedOracleAppPath, client.config.url, app), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, client.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
}

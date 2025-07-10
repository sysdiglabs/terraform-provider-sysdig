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

func (c *Client) GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (identity string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedIdentityPath, c.config.url, provider), nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", c.ErrorFromResponse(response)
	}

	return Unmarshal[string](response.Body)
}

func (c *Client) GetTrustedAzureAppSecure(ctx context.Context, app string) (trusted map[string]string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedAzureAppPath, c.config.url, app), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
}

func (c *Client) GetTenantExternalIDSecure(ctx context.Context) (tenant string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTenantExternaIDPath, c.config.url), nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return "", c.ErrorFromResponse(response)
	}

	return Unmarshal[string](response.Body)
}

func (c *Client) GetAgentlessScanningAssetsSecure(ctx context.Context) (assets map[string]any, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingAgentlessScanningAssetsPath, c.config.url), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]any](response.Body)
}

func (c *Client) GetCloudIngestionAssetsSecure(ctx context.Context, provider, providerID, componentType string) (assets map[string]any, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingCloudIngestionAssetsPath, c.config.url, provider, providerID, componentType), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]any](response.Body)
}

func (c *Client) GetTrustedCloudRegulationAssetsSecure(ctx context.Context, provider string) (assets map[string]string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedRegulationAssetsPath, c.config.url, provider), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
}

func (c *Client) GetTrustedOracleAppSecure(ctx context.Context, app string) (trustedApp map[string]string, err error) {
	response, err := c.requester.Request(ctx, http.MethodGet, fmt.Sprintf(onboardingTrustedOracleAppPath, c.config.url, app), nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if dErr := response.Body.Close(); dErr != nil {
			err = fmt.Errorf("unable to close response body: %w", dErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, c.ErrorFromResponse(response)
	}

	return Unmarshal[map[string]string](response.Body)
}

package v2

import (
	"context"
	"fmt"
	"net/http"
)

const (
	onboardingTrustedIdentityPath = "%s/api/secure/onboarding/v2/trustedIdentity?provider=%s"
	onboardingTenantExternaIDPath = "%s/api/secure/onboarding/v2/externalID"
)

type OnboardingSecureInterface interface {
	Base
	GetTrustedCloudIdentitySecure(ctx context.Context, provider string) (string, error)
	GetTenantExternalIDSecure(ctx context.Context) (string, error)
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

//go:build unit

package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIBMClient_DoIBMRequest(t *testing.T) {
	testTable := []struct {
		TokenExpiration  int64
		Iterations       int
		ExpectedIAMCalls int
	}{
		{
			TokenExpiration:  time.Now().Add(-time.Hour).Unix(),
			Iterations:       3,
			ExpectedIAMCalls: 3,
		},
		{
			TokenExpiration:  time.Now().Add(time.Hour).Unix(),
			Iterations:       3,
			ExpectedIAMCalls: 1,
		},
	}
	for _, testCase := range testTable {
		instanceID := "instance ID"
		apiKey := "api key"
		token := "token"
		iamEndpointCalled := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == IBMIAMPath {
				iamEndpointCalled++
				data, err := json.Marshal(IAMTokenResponse{
					AccessToken: token,
					Expiration:  testCase.TokenExpiration,
				})
				if err != nil {
					t.Errorf("failed to create IAM response, err: %v", err)
				}

				_, err = w.Write(data)
				if err != nil {
					t.Errorf("failed to create IAM response, err: %v", err)
				}
				return
			}

			if value := r.Header.Get(AuthorizationHeader); value != fmt.Sprintf("Bearer %s", token) {
				t.Errorf("invalid authorization header, %v", value)
			}
			if value := r.Header.Get(IBMInstanceIDHeader); value != instanceID {
				t.Errorf("expected instance id %v, got %v", instanceID, value)
			}
		}))

		c := newIBMClient(
			WithIBMInstanceID(instanceID),
			WithIBMAPIKey(apiKey),
			WithIBMIamURL(server.URL),
			WithURL(server.URL),
		)

		url := fmt.Sprintf("%s/foo/bar", server.URL)
		for i := 0; i < testCase.Iterations; i++ {
			_, err := c.requester.Request(context.Background(), http.MethodGet, url, nil)
			if err != nil {
				t.Errorf("got error while sending request, err: %v", err)
			}
		}

		if iamEndpointCalled != testCase.ExpectedIAMCalls {
			t.Errorf("expected IAM calls %d, got %d", testCase.ExpectedIAMCalls, iamEndpointCalled)
		}

		server.Close()
	}
}

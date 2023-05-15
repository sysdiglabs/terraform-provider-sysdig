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
			if value := r.Header.Get(SysdigProviderHeader); value != SysdigProviderHeaderValue {
				t.Errorf("expected sysdig provider %v, got %v", SysdigProviderHeaderValue, value)
			}
		}))

		var teamID int
		c := newIBMClient(
			WithIBMInstanceID(instanceID),
			WithIBMAPIKey(apiKey),
			WithIBMIamURL(server.URL),
			WithURL(server.URL),
			WithSysdigTeamID(&teamID),
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

func TestIBMClient_CurrentTeamID(t *testing.T) {
	teamID1 := 1
	teamID2 := 2
	teamID3 := 3
	teamName := "team"
	instanceID := "instance ID"
	apiKey := "api key"
	token := "token"

	testTable := []struct {
		name           string
		opts           []ClientOption
		expectedTeamID int
	}{
		{
			name: "use current team id from user",
			opts: []ClientOption{
				WithSysdigTeamID(nil),
			},
			expectedTeamID: teamID1,
		},
		{
			name: "use specified team id",
			opts: []ClientOption{
				WithSysdigTeamID(&teamID2),
			},
			expectedTeamID: teamID2,
		},
		{
			name: "get team id from team name",
			opts: []ClientOption{
				WithSysdigTeamName(teamName),
			},
			expectedTeamID: teamID3,
		},
		{
			name: "team id has priority over team name",
			opts: []ClientOption{
				WithSysdigTeamID(&teamID2),
				WithSysdigTeamName(teamName),
			},
			expectedTeamID: teamID2,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var data []byte

		switch r.URL.Path {
		case fmt.Sprintf("%steam", GetTeamByNamePath):
			data, err = json.Marshal(teamWrapper{Team: Team{
				ID: teamID3,
			}})
			if err != nil {
				t.Errorf("failed to create team response, err: %v", err)
			}
		case IBMIAMPath:
			data, err = json.Marshal(IAMTokenResponse{
				AccessToken: token,
				Expiration:  time.Now().Add(time.Hour).Unix(),
			})
			if err != nil {
				t.Errorf("failed to create IAM response, err: %v", err)
			}
		case GetMePath:
			data, err = json.Marshal(userWrapper{
				User: User{
					CurrentTeam: &teamID1,
				},
			})
			if err != nil {
				t.Errorf("failed to create user response, err: %v", err)
			}
		}

		_, err = w.Write(data)
		if err != nil {
			t.Errorf("failed to create response, err: %v", err)
		}
		return
	}))

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			opts := []ClientOption{
				WithIBMInstanceID(instanceID),
				WithIBMAPIKey(apiKey),
				WithIBMIamURL(server.URL),
				WithURL(server.URL),
			}
			opts = append(opts, testCase.opts...)
			c := newIBMClient(opts...)

			id, err := c.CurrentTeamID(context.Background())
			if err != nil {
				t.Errorf("got error while getting current team ID: %v", err)
			}

			if id != testCase.expectedTeamID {
				t.Errorf("expected team ID %d, got %d", testCase.expectedTeamID, id)
			}
		})
	}
}

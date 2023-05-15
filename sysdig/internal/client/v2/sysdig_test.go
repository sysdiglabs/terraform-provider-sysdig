//go:build unit

package v2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSysdigRequest(t *testing.T) {
	t.Parallel()
	type foo struct {
		Number int `json:"number"`
	}
	token := "token"
	given := foo{Number: 15}
	extraHeader := "extra-header"
	extraHeaderValue := "extra-header-value"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value := r.Header.Get(AuthorizationHeader); value != fmt.Sprintf("Bearer %s", token) {
			t.Errorf("invalid authorization header, %v", value)
		}
		if value := r.Header.Get(extraHeader); value != extraHeaderValue {
			t.Errorf("invalid extra header %v", value)
		}
		if value := r.Header.Get(SysdigProviderHeader); value != SysdigProviderHeaderValue {
			t.Errorf("expected sysdig provider %v, got %v", SysdigProviderHeaderValue, value)
		}
		agent := r.Header.Get(UserAgentHeader)
		agentParts := strings.Split(agent, "/")
		if len(agentParts) != 2 || agentParts[0] != SysdigProviderHeaderValue || agentParts[1] == "" {
			t.Errorf("invalid user agent: %v", agent)
		}
		unmarshalled, err := Unmarshal[foo](r.Body)
		if err != nil {
			t.Errorf("failed to unmarshal payload, err: %v", err)
		}
		if given != unmarshalled {
			t.Errorf("expected %v, got %v", given, unmarshalled)
		}
	}))
	defer server.Close()

	c := newSysdigClient(
		WithURL(server.URL),
		WithInsecure(true),
		WithToken(token),
		WithExtraHeaders(map[string]string{
			extraHeader: extraHeaderValue,
		}),
	)

	payload, err := Marshal(given)
	if err != nil {
		t.Errorf("failed to marshal payload, err: %v", err)
	}
	_, err = c.requester.Request(context.Background(), http.MethodPost, server.URL, payload)
	if err != nil {
		t.Errorf("got error while sending request, err: %v", err)
	}
}

func TestSysdigClient_CurrentTeamID(t *testing.T) {
	token := "token"
	teamID := 1

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(userWrapper{
			User: User{
				CurrentTeam: &teamID,
			},
		})
		if err != nil {
			t.Errorf("failed to create user response, err: %v", err)
		}

		_, err = w.Write(data)
		if err != nil {
			t.Errorf("failed to create response, err: %v", err)
		}
	}))

	c := newSysdigClient(
		WithURL(server.URL),
		WithToken(token),
	)

	id, err := c.CurrentTeamID(context.Background())
	if err != nil {
		t.Errorf("failed to get current team id, %v", err)
	}

	if id != teamID {
		t.Errorf("expecting team id %d, got %d", teamID, id)
	}
}

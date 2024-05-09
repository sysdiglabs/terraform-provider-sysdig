//go:build unit

package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	t.Parallel()
	type foo struct {
		Number int `json:"number"`
	}

	given := &foo{Number: 15}
	expected := `{"number":15}`

	data, err := Marshal(given)
	if err != nil {
		t.Errorf("failed to marshal %v", err)
	}

	buf := &strings.Builder{}
	_, err = io.Copy(buf, data)
	if err != nil {
		t.Errorf("failed to populate buffer, err: %v", err)
	}

	marshaled := buf.String()
	if marshaled != expected {
		t.Errorf("expected %v, got %v", expected, marshaled)
	}
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	type foo struct {
		Number int `json:"number"`
	}
	given := `{"number":15}`
	expected := foo{Number: 15}

	unmarshalled, err := Unmarshal[foo](io.NopCloser(strings.NewReader(given)))
	if err != nil {
		t.Errorf("got error while unmarshaling, err: %v", err)
	}

	if expected != unmarshalled {
		t.Errorf("expected %v, got %v", expected, unmarshalled)
	}
}

func TestClient_ErrorFromResponse_non_json(t *testing.T) {

	givenPayload := "non json body"
	expected := "401 Unauthorized"
	c := Client{}

	resp := &http.Response{
		Status: "401 Unauthorized",
		Body:   io.NopCloser(strings.NewReader(givenPayload)),
	}
	err := c.ErrorFromResponse(resp)
	if err.Error() != expected {
		t.Errorf("expected err %v, got %v", expected, err)
	}
}

func TestClient_ErrorFromResponse_standard_error_format(t *testing.T) {
	type Error struct {
		Reason  string `json:"reason"`
		Message string `json:"message"`
	}

	type Errors struct {
		Errors []Error `json:"errors"`
	}

	given := Errors{
		Errors: []Error{
			{
				Reason:  "error1",
				Message: "message1",
			},
			{
				Reason:  "error2",
				Message: "message2",
			},
		},
	}
	expected := "error1, message1, error2, message2"
	c := Client{}
	payload, err := json.Marshal(given)
	if err != nil {
		t.Errorf("failed to marshal errors, %v", err)
	}

	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(string(payload))),
	}
	err = c.ErrorFromResponse(resp)
	if err.Error() != expected {
		t.Errorf("expected err %v, got %v", expected, err)
	}
}

func TestClient_ErrorFromResponse_standard_error_format_2(t *testing.T) {

	givenPayload := `
	{
		"timestamp" : 1715255725613,
		"status" : 401,
		"error" : "Unauthorized",
		"path" : "/api/v2/alerts/46667521"
	}
	`
	expected := "Unauthorized"
	c := Client{}

	resp := &http.Response{
		Status: "401 Unauthorized",
		Body:   io.NopCloser(strings.NewReader(givenPayload)),
	}
	err := c.ErrorFromResponse(resp)
	if err.Error() != expected {
		t.Errorf("expected err %v, got %v", expected, err)
	}
}

func TestClient_ErrorFromResponse_json_nonStandard_error_format(t *testing.T) {
	type Error struct {
		Reason  string `json:"nonStandardFieldNameReason"`
		Message string `json:"nonStandardFieldNameMessage"`
	}

	type Errors struct {
		Errors []Error `json:"errors"`
	}

	given := Errors{
		Errors: []Error{
			{
				Reason:  "error1",
				Message: "message1",
			},
			{
				Reason:  "error2",
				Message: "message2",
			},
		},
	}
	expected := "401 Unauthorized"
	c := Client{}
	payload, err := json.Marshal(given)
	if err != nil {
		t.Errorf("failed to marshal errors, %v", err)
	}

	resp := &http.Response{
		Status: "401 Unauthorized",
		Body:   io.NopCloser(strings.NewReader(string(payload))),
	}
	err = c.ErrorFromResponse(resp)
	if err.Error() != expected {
		t.Errorf("expected err %v, got %v", expected, err)
	}
}

func TestRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agent := r.Header.Get(UserAgentHeader)
		agentParts := strings.Split(agent, "/")
		if len(agentParts) != 2 || agentParts[0] != SysdigUserAgentHeaderValue || agentParts[1] == "" {
			t.Errorf("invalid user agent: %v", agent)
		}
	}))

	cfg := &config{
		url: server.URL,
	}
	client := newHTTPClient(cfg)

	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/", cfg.url), nil)
	if err != nil {
		t.Errorf("failed to create request, %v", err)
	}

	_, err = request(client, cfg, r)
	if err != nil {
		t.Errorf("failed to send request, %v", err)
	}
}

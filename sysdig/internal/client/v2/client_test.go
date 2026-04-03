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

func TestClient_APIErrorFromResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		statusCode     int
		status         string
		body           string
		wantMessage    string
		wantStatusCode int
	}{
		{
			name:           "message field",
			statusCode:     400,
			status:         "400 Bad Request",
			body:           `{"message":"invalid zone name"}`,
			wantMessage:    "invalid zone name",
			wantStatusCode: 400,
		},
		{
			name:           "error field",
			statusCode:     401,
			status:         "401 Unauthorized",
			body:           `{"timestamp":1715255725613,"status":401,"error":"Unauthorized","path":"/api/v2/zones/123"}`,
			wantMessage:    "Unauthorized",
			wantStatusCode: 401,
		},
		{
			name:       "errors array with reason and message",
			statusCode: 422,
			status:     "422 Unprocessable Entity",
			body: `{"errors":[
				{"reason":"validation_error","message":"name is required"},
				{"reason":"validation_error","message":"scope is required"}
			]}`,
			wantMessage:    "validation_error, name is required, validation_error, scope is required",
			wantStatusCode: 422,
		},
		{
			name:           "details array",
			statusCode:     400,
			status:         "400 Bad Request",
			body:           `{"details":["field 'name' is required","field 'scope' is required"]}`,
			wantMessage:    "field 'name' is required, field 'scope' is required",
			wantStatusCode: 400,
		},
		{
			name:           "non-json body falls back to status",
			statusCode:     502,
			status:         "502 Bad Gateway",
			body:           "not json",
			wantMessage:    "502 Bad Gateway",
			wantStatusCode: 502,
		},
		{
			name:           "empty json falls back to status",
			statusCode:     500,
			status:         "500 Internal Server Error",
			body:           `{}`,
			wantMessage:    "500 Internal Server Error",
			wantStatusCode: 500,
		},
		{
			name:           "404 not found",
			statusCode:     404,
			status:         "404 Not Found",
			body:           `{"message":"zone not found"}`,
			wantMessage:    "zone not found",
			wantStatusCode: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Client{}
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Status:     tt.status,
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			err := c.APIErrorFromResponse(resp)

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.StatusCode != tt.wantStatusCode {
				t.Errorf("StatusCode: want %d, got %d", tt.wantStatusCode, apiErr.StatusCode)
			}
			if apiErr.Error() != tt.wantMessage {
				t.Errorf("Message: want %q, got %q", tt.wantMessage, apiErr.Error())
			}
		})
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

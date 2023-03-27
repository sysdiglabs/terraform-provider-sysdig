//go:build unit

package v2

import (
	"encoding/json"
	"io"
	"net/http"
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

func TestClient_ErrorFromResponse(t *testing.T) {
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

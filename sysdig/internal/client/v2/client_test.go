package v2

import (
	"io"
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

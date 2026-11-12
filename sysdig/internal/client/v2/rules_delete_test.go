//go:build unit

package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteRuleStatusHandling(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"ok", http.StatusOK, false},
		{"no content", http.StatusNoContent, false},
		{"already deleted", http.StatusNotFound, false},
		{"bad request", http.StatusBadRequest, true},
		{"server error", http.StatusInternalServerError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(`{"message":"Failed to find a rule with id: 1"}`))
			}))
			defer server.Close()

			c := newSysdigClient(WithURL(server.URL), WithToken("test-token"))

			// Without this the transport retries 5xx with ~31s of backoff, which would test
			// retry exhaustion rather than DeleteRule's status handling.
			err := c.DeleteRule(WithoutTransportRetries(context.Background()), 1)
			if (err != nil) != tt.wantErr {
				t.Fatalf("DeleteRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

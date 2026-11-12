//go:build unit

package v2

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
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

// A DELETE that outlives the gateway timeout returns 504 but still completes server-side, so the
// transport's retry sees 404. Every delete helper must treat that as success.
func TestDeleteSucceedsWhenRetryAfterGatewayTimeoutGets404(t *testing.T) {
	t.Parallel()

	deletes := map[string]func(*Client) error{
		"DeleteRule":            func(c *Client) error { return c.DeleteRule(context.Background(), 1) },
		"DeleteStatefulRule":    func(c *Client) error { return c.DeleteStatefulRule(context.Background(), 1) },
		"DeletePolicy":          func(c *Client) error { return c.DeletePolicy(context.Background(), 1) },
		"DeleteCompositePolicy": func(c *Client) error { return c.DeleteCompositePolicy(context.Background(), 1) },
		"DeleteList":            func(c *Client) error { return c.DeleteList(context.Background(), 1) },
	}

	for name, del := range deletes {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if calls.Add(1) == 1 {
					w.WriteHeader(http.StatusGatewayTimeout)
					return
				}
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"message":"not found"}`))
			}))
			defer server.Close()

			c := newSysdigClient(WithURL(server.URL), WithToken("test-token"), withRetryWaitMin(time.Millisecond))

			if err := del(c); err != nil {
				t.Fatalf("expected success, got %v", err)
			}
			if got := calls.Load(); got != 2 {
				t.Fatalf("expected 2 requests (504 then 404), got %d", got)
			}
		})
	}
}

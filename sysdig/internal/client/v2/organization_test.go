//go:build unit

package v2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

func TestMarshalOrg(t *testing.T) {
	t.Parallel()
	c := Client{}

	given := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "58ca66a5-ac87-497b-a501-7a4c934b3017",
			Provider:            cloudauth.Provider_PROVIDER_GCP,
		},
	}
	expected := `{"managementAccountId":"58ca66a5-ac87-497b-a501-7a4c934b3017", "provider":"PROVIDER_GCP"}`

	payload, err := c.marshalCloudauthProto(given)
	if err != nil {
		t.Errorf("failed to marshal payload, err: %v", err)
	}

	buf := &strings.Builder{}
	_, err = io.Copy(buf, payload)
	if err != nil {
		t.Errorf("failed to populate buffer, err: %v", err)
	}
	marshaled := buf.String()

	if strings.ReplaceAll(marshaled, " ", "") != strings.ReplaceAll(expected, " ", "") {
		t.Errorf("expected %v, got %v", expected, marshaled)
	}
}

func TestUpdateOrganizationSecureUsesCamelCase(t *testing.T) {
	t.Parallel()

	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		// Return a valid protobuf JSON response
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"managementAccountId":"31ebd166-82ef-4ca4-baf9-ce760afa46fb","organizationRootId":"r-8llw"}`))
	}))
	defer server.Close()

	c := newSysdigClient(
		WithURL(server.URL),
		WithToken("test-token"),
	)

	org := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "31ebd166-82ef-4ca4-baf9-ce760afa46fb",
			OrganizationRootId:  "r-8llw",
		},
	}

	_, _, err := c.UpdateOrganizationSecure(context.Background(), "e69f12fd-934d-43cc-8b8c-2964aba20003", org)
	if err != nil {
		t.Fatalf("UpdateOrganizationSecure failed: %v", err)
	}

	if strings.Contains(receivedBody, "management_account_id") {
		t.Errorf("request body uses snake_case (json.Marshal), expected camelCase (protojson.Marshal): %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "managementAccountId") {
		t.Errorf("request body missing camelCase field 'managementAccountId': %s", receivedBody)
	}
	if strings.Contains(receivedBody, "organization_root_id") {
		t.Errorf("request body uses snake_case 'organization_root_id', expected camelCase 'organizationRootId': %s", receivedBody)
	}
	if !strings.Contains(receivedBody, "organizationRootId") {
		t.Errorf("request body missing camelCase field 'organizationRootId': %s", receivedBody)
	}
}

func TestUnmarshalOrg(t *testing.T) {
	t.Parallel()
	c := Client{}
	given := `{"managementAccountId":"58ca66a5-ac87-497b-a501-7a4c934b3017","provider":"PROVIDER_GCP"}`
	expected := &OrganizationSecure{
		cloudauth.CloudOrganization{
			ManagementAccountId: "58ca66a5-ac87-497b-a501-7a4c934b3017",
			Provider:            cloudauth.Provider_PROVIDER_GCP,
		},
	}

	unmarshalled := &OrganizationSecure{}
	err := c.unmarshalCloudauthProto(io.NopCloser(strings.NewReader(given)), unmarshalled)
	if err != nil {
		t.Errorf("got error while unmarshaling, err: %v", err)
	}

	if expected.String() != unmarshalled.String() {
		t.Errorf("expected %v, got %v", expected, unmarshalled)
	}
}

func TestOrganizationURLsAsyncFlag(t *testing.T) {
	t.Parallel()

	const orgID = "4c53102d-6846-447b-bfd1-4c0d5002cddf"
	const base = "http://localhost/api/cloudauth/v1/organizations"

	for _, async := range []bool{false, true} {
		t.Run(map[bool]string{false: "disabled", true: "enabled"}[async], func(t *testing.T) {
			t.Parallel()
			c := newSysdigClient(WithURL("http://localhost"), WithOrgAPIAsync(async))

			suffix := ""
			if async {
				suffix = "?async=true"
			}
			if got, want := c.withAsync(c.organizationsURL()), base+suffix; got != want {
				t.Errorf("collection URL = %q, want %q", got, want)
			}
			if got, want := c.withAsync(c.organizationURL(orgID)), base+"/"+orgID+suffix; got != want {
				t.Errorf("mutation URL = %q, want %q", got, want)
			}
			// the read is never wrapped: GetOrganizationSecure accepts only 200
			if got, want := c.organizationURL(orgID), base+"/"+orgID; got != want {
				t.Errorf("read URL = %q, want %q", got, want)
			}
		})
	}
}

// cloudauth answers an async mutation with 202; update discards the body, so an empty one is
// tolerated there, while create needs an id and so keeps decoding strictly.
func TestOrganizationAsyncAckWithoutBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"), WithOrgAPIAsync(true))
	org := &OrganizationSecure{}

	updated, _, err := c.UpdateOrganizationSecure(context.Background(), "oid", org)
	if err != nil {
		t.Errorf("UpdateOrganizationSecure on a bodyless 202: %v", err)
	}
	if updated != nil {
		t.Errorf("UpdateOrganizationSecure returned %+v, want no organization for an acknowledgement", updated)
	}
	// the client only reports transport and decode failures; the missing id is the resource's call
	created, _, err := c.CreateOrganizationSecure(context.Background(), org)
	if err != nil {
		t.Errorf("CreateOrganizationSecure on a bodyless 202: %v", err)
	} else if created.GetId() != "" {
		t.Errorf("CreateOrganizationSecure returned id %q, want nothing invented for an empty ack", created.GetId())
	}
	// the read keeps rejecting anything but 200, so an unexpected 202 there stays an error
	if _, _, err := c.GetOrganizationSecure(context.Background(), "oid"); err == nil {
		t.Error("GetOrganizationSecure on a 202: expected an error")
	}

	// without async the same answer is not an acknowledgement, and tolerating it would report a
	// synchronous update that returned nothing as a success
	sync := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"))
	if _, _, err := sync.UpdateOrganizationSecure(context.Background(), "oid", org); err == nil {
		t.Error("UpdateOrganizationSecure on a bodyless 202 without async: expected an error")
	}
	if _, _, err := sync.CreateOrganizationSecure(context.Background(), org); err == nil {
		t.Error("CreateOrganizationSecure on a bodyless 202 without async: expected an error")
	}
}

// a requester answering with a fixed status and body, whose Close fails after the read
type closeFailingRequester struct {
	status int
	body   string
}

func (r closeFailingRequester) CurrentTeamID(_ context.Context) (int, error) { return 0, nil }

func (r closeFailingRequester) Request(_ context.Context, _ string, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: r.status,
		Status:     fmt.Sprintf("%d %s", r.status, http.StatusText(r.status)),
		Body:       closeFailingBody{Reader: strings.NewReader(r.body)},
	}, nil
}

type closeFailingBody struct{ io.Reader }

func (closeFailingBody) Close() error { return errors.New("connection reset by peer") }

// a body that fails midway through, as a truncated response does
type unreadableRequester struct{}

func (unreadableRequester) CurrentTeamID(_ context.Context) (int, error) { return 0, nil }

func (unreadableRequester) Request(_ context.Context, _ string, _ string, _ io.Reader) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       closeFailingBody{Reader: iotest.TimeoutReader(strings.NewReader(`{"id":"4c53102d"}`))},
	}, nil
}

// a malformed 202 payload must not be mistaken for an empty acknowledgement
func TestOrganizationAsyncAckMalformedBody(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"error":"quota exceeded"`))
	}))
	defer srv.Close()

	c := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"), WithOrgAPIAsync(true))
	if _, _, err := c.CreateOrganizationSecure(context.Background(), &OrganizationSecure{}); err == nil {
		t.Error("CreateOrganizationSecure on a malformed 202: expected an error")
	}
	if _, _, err := c.UpdateOrganizationSecure(context.Background(), "oid", &OrganizationSecure{}); err == nil {
		t.Error("UpdateOrganizationSecure on a malformed 202: expected an error")
	}
}

// G2: whether an empty body is an acknowledgement follows from the call being async, not from
// which 2xx carried it; a read is never an acknowledgement and still needs the organization.
func TestOrganizationEmptyBodyDependsOnTheCallNotTheStatus(t *testing.T) {
	t.Parallel()

	for _, status := range []int{http.StatusOK, http.StatusCreated, http.StatusAccepted} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			c := newSysdigClient(WithURL("http://localhost"), WithToken("fake-token"), WithOrgAPIAsync(true))
			c.requester = closeFailingRequester{status: status}

			if _, _, err := c.CreateOrganizationSecure(context.Background(), &OrganizationSecure{}); err != nil {
				t.Errorf("CreateOrganizationSecure on an empty %d: %v", status, err)
			}
		})
	}

	// the read has to report it: an empty answer would otherwise become empty state
	c := newSysdigClient(WithURL("http://localhost"), WithToken("fake-token"), WithOrgAPIAsync(true))
	c.requester = closeFailingRequester{status: http.StatusOK}
	if _, _, err := c.GetOrganizationSecure(context.Background(), "oid"); err == nil {
		t.Error("GetOrganizationSecure on an empty 200: expected an error")
	}
}

// The probe answers the one question the wait has, and it must not pay for the transport's retries
// or for the backend listing every member account.
func TestOrganizationProbeIsThinAndUnretried(t *testing.T) {
	t.Parallel()

	t.Run("identity verbosity and no retries", func(t *testing.T) {
		var requests int
		var query string
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests++
			query = r.URL.RawQuery
			w.WriteHeader(http.StatusServiceUnavailable)
		}))
		defer srv.Close()

		c := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"))
		exists, errStatus, err := c.OrganizationExistsSecure(context.Background(), "oid")
		if exists || err == nil {
			t.Fatalf("exists=%v err=%v, want the failure reported", exists, err)
		}
		if requests != 1 {
			t.Errorf("issued %d requests, want one: the transport must not retry a probe", requests)
		}
		if query != "verbosity=VERBOSITY_IDENT" {
			t.Errorf("query = %q, want identity verbosity", query)
		}
		// the status has to survive, or the caller cannot tell a transient 503 from a permanent 403
		if errStatus == "" || errStatus[:3] != "503" {
			t.Errorf("errStatus = %q, want the 503 the server returned", errStatus)
		}
	})

	t.Run("gone statuses", func(t *testing.T) {
		for _, status := range []int{http.StatusNotFound, http.StatusGone} {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
			}))
			c := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"))
			exists, errStatus, err := c.OrganizationExistsSecure(context.Background(), "oid")
			if exists || errStatus != "" || err != nil {
				t.Errorf("%d: exists=%v errStatus=%q err=%v, want a clean answer that it is gone", status, exists, errStatus, err)
			}
			srv.Close()
		}
	})

	t.Run("a body that does not parse still proves it exists", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("<html><body>maintenance</body></html>"))
		}))
		defer srv.Close()

		c := newSysdigClient(WithURL(srv.URL), WithToken("fake-token"))
		exists, _, err := c.OrganizationExistsSecure(context.Background(), "oid")
		if err != nil || !exists {
			t.Errorf("exists=%v err=%v, want a 200 taken as proof the organization is there", exists, err)
		}
	})
}

// The async delete answers 204 with the cascade still running, which is the response the wait is
// built around; nothing else is accepted.
func TestOrganizationDeleteAcceptsWhatTheBackendSends(t *testing.T) {
	t.Parallel()

	for status, wantErr := range map[int]bool{
		http.StatusNoContent: false,
		http.StatusOK:        false,
		http.StatusAccepted:  true,
	} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(status)
		}))
		for _, async := range []bool{false, true} {
			opts := []ClientOption{WithURL(srv.URL), WithToken("fake-token")}
			if async {
				opts = append(opts, WithOrgAPIAsync(true))
			}
			_, err := newSysdigClient(opts...).DeleteOrganizationSecure(context.Background(), "oid")
			if got := err != nil; got != wantErr {
				t.Errorf("status %d async=%v: error = %v, want %v (%v)", status, async, got, wantErr, err)
			}
		}
		srv.Close()
	}
}

// A close failure arrives after the organization has already been created server side. Reporting
// it would stop the resource from recording the id, leaving the organization untracked.
func TestOrganizationCloseFailureDoesNotDiscardCreate(t *testing.T) {
	t.Parallel()

	c := newSysdigClient(WithURL("http://localhost"), WithToken("fake-token"))
	c.requester = closeFailingRequester{status: http.StatusOK, body: `{"id":"4c53102d"}`}

	created, _, err := c.CreateOrganizationSecure(context.Background(), &OrganizationSecure{})
	if err != nil {
		t.Fatalf("CreateOrganizationSecure: %v", err)
	}
	if created.GetId() != "4c53102d" {
		t.Errorf("id = %q, want the created organization to come back so the resource can track it", created.GetId())
	}
}

package sysdig

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureOrganization() *schema.Resource {
	timeout := 5 * time.Minute
	// an async delete waits out the member-account cascade, which outlasts the other operations
	deleteTimeout := 30 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureOrganizationCreate,
		DeleteContext: resourceSysdigSecureOrganizationDelete,
		ReadContext:   resourceSysdigSecureOrganizationRead,
		UpdateContext: resourceSysdigSecureOrganizationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Version: 0,
				Type:    resourceSysdigSecureOrganizationV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceSysdigSecureOrganizationUpgradeV0,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(deleteTimeout),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			SchemaManagementAccountID: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaOrganizationalUnitIds: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaIncludedOrganizationalGroups: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaExcludedOrganizationalGroups: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaIncludedCloudAccounts: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaExcludedCloudAccounts: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaOrganizationRootID: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaAutomaticOnboarding: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func getSecureOrganizationClient(c SysdigClients) (v2.OrganizationSecureInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigSecureOrganizationCreate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)

	orgCreated, errStatus, err := client.CreateOrganizationSecure(ctx, org)
	if err != nil {
		return diag.Errorf("Error creating resource: %s %s", errStatus, err)
	}
	// an empty id would drop the organization out of state while it exists server side
	if orgCreated.GetId() == "" {
		return diag.Errorf("Error creating resource: the organization was accepted but no id was returned")
	}

	data.SetId(orgCreated.Id)

	return nil
}

func resourceSysdigSecureOrganizationDelete(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	errStatus, err := client.DeleteOrganizationSecure(ctx, data.Id())
	if err != nil {
		if organizationIsGone(errStatus) {
			return nil
		}
		return diag.Errorf("Error deleting resource: %s %s", errStatus, err)
	}

	// only an async delete is merely acknowledged; a synchronous one is done when it returns
	if i.(SysdigClients).orgAPIAsyncEnabled() {
		if err := waitForOrganizationDeletion(ctx, client, data.Id(), data.Timeout(schema.TimeoutDelete)); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// a confirmation read carries the transport's own retries, which back off as far as 30 seconds
// each; without a cap one unlucky probe spends that inside a single tick of the wait, burning the
// delete budget on backoff rather than on checking whether the cascade finished
const organizationDeletionProbeTimeout = 10 * time.Second

// the timeout is passed in rather than read from the resource so it can be exercised in tests;
// the SDK has already put the same deadline on ctx
func waitForOrganizationDeletion(ctx context.Context, client v2.OrganizationSecureInterface, orgID string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		probeCtx, cancel := context.WithTimeout(ctx, organizationDeletionProbeTimeout)
		defer cancel()

		_, errStatus, err := client.GetOrganizationSecure(probeCtx, orgID)
		if err == nil {
			return retry.RetryableError(fmt.Errorf("organization %s has not been deleted yet; the deletion may also have failed server side", orgID))
		}

		return classifyDeletionProbe(orgID, errStatus, err)
	})
}

// a gone organization is reported as either of these, and both mean there is nothing left to act on
func organizationIsGone(errStatus string) bool {
	code := statusCodeFromStatus(errStatus)
	return code == http.StatusNotFound || code == http.StatusGone
}

// nil means the organization is gone; anything retryable is worth another poll, and everything
// else has to be reported instead of waiting out the whole delete timeout.
func classifyDeletionProbe(orgID, errStatus string, err error) *retry.RetryError {
	retryable := func() *retry.RetryError {
		return retry.RetryableError(fmt.Errorf("confirming deletion of organization %s: %s %w", orgID, errStatus, err))
	}
	fatal := func() *retry.RetryError {
		return retry.NonRetryableError(fmt.Errorf("confirming deletion of organization %s: %s %w", orgID, errStatus, err))
	}

	if organizationIsGone(errStatus) {
		return nil
	}
	if code := statusCodeFromStatus(errStatus); code != 0 {
		switch {
		// the same set the transport itself retries: 409, 429 and 5xx other than 501
		case code == http.StatusConflict, code == http.StatusTooManyRequests:
			return retryable()
		case code >= 500:
			if code == http.StatusNotImplemented {
				return fatal()
			}
			return retryable()
		case code >= 400:
			return fatal()
		}
		// an unexpected 2xx or 3xx still means the organization answered, so keep waiting for it
		return retryable()
	}

	// without a status the failure is transport level or a body problem; only the deterministic
	// ones are hopeless: a body that does not parse will not parse on the next attempt either,
	// and for the transport the same classification the request layer uses already knows which
	var malformed *v2.MalformedBodyError
	if errors.As(err, &malformed) {
		return fatal()
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		// asked with a fresh context on purpose: the wait's own deadline, once expired, makes the
		// policy report any error as permanent and the failure would blame the network instead
		if transient, _ := retryablehttp.DefaultRetryPolicy(context.Background(), nil, urlErr); !transient {
			return fatal()
		}
	}
	return retryable()
}

// the organization endpoints report a failure as an HTTP status line plus a plain error, so the
// code has to be taken off that line
func statusCodeFromStatus(status string) int {
	fields := strings.Fields(status)
	if len(fields) == 0 {
		return 0
	}
	code, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0
	}
	return code
}

func resourceSysdigSecureOrganizationRead(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org, errStatus, err := client.GetOrganizationSecure(ctx, data.Id())
	if err != nil {
		if organizationIsGone(errStatus) {
			data.SetId("")
			return nil
		}
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	err = secureOrganizationToResourceData(data, org)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureOrganizationUpdate(ctx context.Context, data *schema.ResourceData, i any) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)
	orgID := data.Id()

	_, errStatus, err := client.UpdateOrganizationSecure(ctx, orgID, org)
	if err != nil {
		// clearing the id here instead would hand Terraform an empty state for an update it
		// planned, which it rejects as an inconsistent result
		if organizationIsGone(errStatus) {
			return diag.Errorf("Error updating resource: organization %s no longer exists", orgID)
		}
		return diag.Errorf("Error updating resource: %s %s", errStatus, err)
	}

	// the answer to an update is not read into state: whatever it carries, mapping it would let a
	// body that omits fields overwrite what was just applied. The plan holds what the organization
	// was set to, and the next refresh reads the server's own view through Read.
	return nil
}

func secureOrganizationFromResourceData(data *schema.ResourceData) *v2.OrganizationSecure {
	secureOrganization := &v2.OrganizationSecure{CloudOrganization: cloudauth.CloudOrganization{}}
	secureOrganization.ManagementAccountId = data.Get(SchemaManagementAccountID).(string)
	secureOrganization.OrganizationRootId = data.Get(SchemaOrganizationRootID).(string)
	secureOrganization.AutomaticOnboarding = data.Get(SchemaAutomaticOnboarding).(bool)
	secureOrganization.OrganizationalUnitIds = schemaSetToList(data.Get(SchemaOrganizationalUnitIds))
	secureOrganization.IncludedOrganizationalGroups = schemaSetToList(data.Get(SchemaIncludedOrganizationalGroups))
	secureOrganization.ExcludedOrganizationalGroups = schemaSetToList(data.Get(SchemaExcludedOrganizationalGroups))
	secureOrganization.IncludedCloudAccounts = schemaSetToList(data.Get(SchemaIncludedCloudAccounts))
	secureOrganization.ExcludedCloudAccounts = schemaSetToList(data.Get(SchemaExcludedCloudAccounts))

	return secureOrganization
}

func secureOrganizationToResourceData(data *schema.ResourceData, org *v2.OrganizationSecure) error {
	err := data.Set(SchemaManagementAccountID, org.ManagementAccountId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaOrganizationalUnitIds, org.OrganizationalUnitIds)
	if err != nil {
		return err
	}

	err = data.Set(SchemaIncludedOrganizationalGroups, org.IncludedOrganizationalGroups)
	if err != nil {
		return err
	}

	err = data.Set(SchemaExcludedOrganizationalGroups, org.ExcludedOrganizationalGroups)
	if err != nil {
		return err
	}

	err = data.Set(SchemaIncludedCloudAccounts, org.IncludedCloudAccounts)
	if err != nil {
		return err
	}

	err = data.Set(SchemaExcludedCloudAccounts, org.ExcludedCloudAccounts)
	if err != nil {
		return err
	}

	err = data.Set(SchemaOrganizationRootID, org.OrganizationRootId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaAutomaticOnboarding, org.AutomaticOnboarding)
	if err != nil {
		return err
	}

	err = data.Set(SchemaIDKey, org.Id)
	if err != nil {
		return err
	}

	return nil
}

// resourceSysdigSecureOrganizationV0 is the schema as released before the include/exclude
// collections became sets, so state written by those versions is decoded with the types it was
// written with. The collections carry no nested state, hence the no-op upgrade.
func resourceSysdigSecureOrganizationV0() *schema.Resource {
	stringList := func() *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		}
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SchemaIDKey:                        {Type: schema.TypeString},
			SchemaManagementAccountID:          {Type: schema.TypeString},
			SchemaOrganizationalUnitIds:        stringList(),
			SchemaIncludedOrganizationalGroups: stringList(),
			SchemaExcludedOrganizationalGroups: stringList(),
			SchemaIncludedCloudAccounts:        stringList(),
			SchemaExcludedCloudAccounts:        stringList(),
			SchemaOrganizationRootID:           {Type: schema.TypeString},
			SchemaAutomaticOnboarding:          {Type: schema.TypeBool},
		},
	}
}

func resourceSysdigSecureOrganizationUpgradeV0(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	return rawState, nil
}

package sysdig

import (
	"context"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureOrganization() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureOrganizationCreate,
		DeleteContext: resourceSysdigSecureOrganizationDelete,
		ReadContext:   resourceSysdigSecureOrganizationRead,
		UpdateContext: resourceSysdigSecureOrganizationUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			SchemaManagementAccountId: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaOrganizationalUnitIds: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaIncludedOrganizationalGroups: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaExcludedOrganizationalGroups: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaIncludedCloudAccounts: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaExcludedCloudAccounts: {
				Type:     schema.TypeList,
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

func resourceSysdigSecureOrganizationCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)

	orgCreated, errStatus, err := client.CreateOrganizationSecure(ctx, org)
	if err != nil {
		return diag.Errorf("Error creating resource: %s %s", errStatus, err)
	}

	data.SetId(orgCreated.Id)

	return nil
}

func resourceSysdigSecureOrganizationDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	errStatus, err := client.DeleteOrganizationSecure(ctx, data.Id())
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error deleting resource: %s %s", errStatus, err)
	}

	return nil
}

func resourceSysdigSecureOrganizationRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org, errStatus, err := client.GetOrganizationSecure(ctx, data.Id())
	if err != nil {
		if strings.Contains(errStatus, "404") {
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

func resourceSysdigSecureOrganizationUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getSecureOrganizationClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	org := secureOrganizationFromResourceData(data)

	_, errStatus, err := client.UpdateOrganizationSecure(ctx, data.Id(), org)
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error updating resource: %s %s", errStatus, err)
	}

	return nil
}

func secureOrganizationFromResourceData(data *schema.ResourceData) *v2.OrganizationSecure {
	secureOrganization := &v2.OrganizationSecure{CloudOrganization: cloudauth.CloudOrganization{}}
	secureOrganization.CloudOrganization.ManagementAccountId = data.Get(SchemaManagementAccountId).(string)
	secureOrganization.CloudOrganization.OrganizationRootId = data.Get(SchemaOrganizationRootID).(string)
	secureOrganization.CloudOrganization.AutomaticOnboarding = data.Get(SchemaAutomaticOnboarding).(bool)
	organizationalUnitIdsData := data.Get(SchemaOrganizationalUnitIds).([]interface{})
	for _, organizationalUnitIdData := range organizationalUnitIdsData {
		secureOrganization.CloudOrganization.OrganizationalUnitIds = append(
			secureOrganization.CloudOrganization.OrganizationalUnitIds,
			organizationalUnitIdData.(string),
		)
	}

	includedOrganizationalGroups := data.Get(SchemaIncludedOrganizationalGroups).([]interface{})
	for _, includedOrganizationalGroup := range includedOrganizationalGroups {
		secureOrganization.CloudOrganization.IncludedOrganizationalGroups = append(
			secureOrganization.CloudOrganization.IncludedOrganizationalGroups,
			includedOrganizationalGroup.(string),
		)
	}

	excludedOrganizationalGroups := data.Get(SchemaExcludedOrganizationalGroups).([]interface{})
	for _, excludedOrganizationalGroup := range excludedOrganizationalGroups {
		secureOrganization.CloudOrganization.ExcludedOrganizationalGroups = append(
			secureOrganization.CloudOrganization.ExcludedOrganizationalGroups,
			excludedOrganizationalGroup.(string),
		)
	}

	includedCloudAccounts := data.Get(SchemaIncludedCloudAccounts).([]interface{})
	for _, includedCloudAccount := range includedCloudAccounts {
		secureOrganization.CloudOrganization.IncludedCloudAccounts = append(
			secureOrganization.CloudOrganization.IncludedCloudAccounts,
			includedCloudAccount.(string),
		)
	}

	excludedCloudAccounts := data.Get(SchemaExcludedCloudAccounts).([]interface{})
	for _, excludedCloudAccount := range excludedCloudAccounts {
		secureOrganization.CloudOrganization.ExcludedCloudAccounts = append(
			secureOrganization.CloudOrganization.ExcludedCloudAccounts,
			excludedCloudAccount.(string),
		)
	}
	return secureOrganization
}

func secureOrganizationToResourceData(data *schema.ResourceData, org *v2.OrganizationSecure) error {
	err := data.Set(SchemaManagementAccountId, org.ManagementAccountId)
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

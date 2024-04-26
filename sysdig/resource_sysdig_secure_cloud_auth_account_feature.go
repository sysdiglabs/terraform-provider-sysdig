package sysdig

import (
	"context"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureCloudauthAccountFeature() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureCloudauthAccountFeatureCreate,
		UpdateContext: resourceSysdigSecureCloudauthAccountFeatureUpdate,
		ReadContext:   resourceSysdigSecureCloudauthAccountFeatureRead,
		DeleteContext: resourceSysdigSecureCloudauthAccountFeatureDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		Schema: getAccountFeatureSchema(),
	}
}

func getAccountFeatureSchema() map[string]*schema.Schema {
	// for AccountFeature resource, account_id & featureType are needed additionally
	featureSchema := map[string]*schema.Schema{
		SchemaAccountId: {
			Type:     schema.TypeString,
			Required: true,
		},
		SchemaFeatureType: {
			Type:     schema.TypeString,
			Required: true,
		},
		SchemaFeatureEnabled: {
			Type:     schema.TypeBool,
			Required: true,
		},
		SchemaFeatureFlags: {
			Type:     schema.TypeMap,
			Optional: true,
		},
		SchemaFeatureComponents: {
			Type:     schema.TypeMap,
			Required: true,
		},
	}

	for field, schema := range accountFeature.Schema {
		featureSchema[field] = schema
	}
	return featureSchema
}

func getSecureCloudauthAccountFeatureClient(client SysdigClients) (v2.CloudauthAccountFeatureSecureInterface, error) {
	return client.sysdigSecureClientV2()
}

func resourceSysdigSecureCloudauthAccountFeatureCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountFeatureClient((meta.(SysdigClients)))
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := data.Get(SchemaAccountId).(string)
	cloudauthAccountFeature, errStatus, err := client.CreateCloudauthAccountFeatureSecure(ctx, accountId, cloudauthAccountFeatureFromResourceData(data))
	if err != nil {
		return diag.Errorf("Error creating resource: %s %s", errStatus, err)
	}

	// using tuple 'accountId/featureType' as TF resource identifier
	data.SetId(accountId + "/" + cloudauthAccountFeature.GetType().String())
	err = data.Set(SchemaAccountId, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountFeatureRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountFeatureClient((meta.(SysdigClients)))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccountFeature, errStatus, err := client.GetCloudauthAccountFeatureSecure(
		ctx, data.Get(SchemaAccountId).(string), data.Get(SchemaType).(string))

	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	err = cloudauthAccountFeatureToResourceData(data, cloudauthAccountFeature)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountFeatureUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountFeatureClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := data.Get(SchemaAccountId).(string)
	existingCloudAccountFeature, errStatus, err := client.GetCloudauthAccountFeatureSecure(
		ctx, accountId, data.Get(SchemaType).(string))
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	newCloudAccountFeature := cloudauthAccountFeaturetFromResourceData(data)

	// validate and reject non-updatable resource schema fields upfront
	err = validateCloudauthAccountFeatureUpdate(existingCloudAccountFeature, newCloudAccountFeature)
	if err != nil {
		return diag.Errorf("Error updating resource: %s", err)
	}

	_, errStatus, err = client.UpdateCloudauthAccountFeatureSecure(
		ctx, accountId, data.Get(SchemaType).(string), newCloudAccountFeature)
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error updating resource: %s %s", errStatus, err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountFeatureDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountFeatureClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	errStatus, err := client.DeleteCloudauthAccountFeatureSecure(
		ctx, data.Get(SchemaAccountId).(string), data.Get(SchemaType).(string))
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error deleting resource: %s %s", errStatus, err)
	}

	return nil
}

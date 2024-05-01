package sysdig

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
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
	// though the schema fields are already defined in cloud_auth_account resource, for AccountFeature
	// calls they are required fields. Also, account_id & flags are needed additionally.
	featureSchema := map[string]*schema.Schema{
		SchemaAccountId: {
			Type:     schema.TypeString,
			Required: true,
		},
		SchemaType: {
			Type:     schema.TypeString,
			Required: true,
		},
		SchemaEnabled: {
			Type:     schema.TypeBool,
			Required: true,
		},
		SchemaComponents: {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		SchemaFeatureFlags: {
			Type:     schema.TypeMap,
			Optional: true,
		},
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
	cloudauthAccountFeature, errStatus, err := client.CreateOrUpdateCloudauthAccountFeatureSecure(
		ctx, accountId, data.Get(SchemaType).(string), cloudauthAccountFeatureFromResourceData(data))
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

	newCloudAccountFeature := cloudauthAccountFeatureFromResourceData(data)

	// validate and reject non-updatable resource schema fields upfront
	err = validateCloudauthAccountFeatureUpdate(existingCloudAccountFeature, newCloudAccountFeature)
	if err != nil {
		return diag.Errorf("Error updating resource: %s", err)
	}

	_, errStatus, err = client.CreateOrUpdateCloudauthAccountFeatureSecure(
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

/*
This function validates and restricts any fields not allowed to be updated during resource updates.
*/
func validateCloudauthAccountFeatureUpdate(existingFeature *v2.CloudauthAccountFeatureSecure, newFeature *v2.CloudauthAccountFeatureSecure) error {
	if existingFeature.Type != newFeature.Type {
		errorInvalidResourceUpdate := fmt.Sprintf("Bad Request. Updating restricted fields not allowed: %s", []string{"type"})
		return errors.New(errorInvalidResourceUpdate)
	}

	return nil
}

func getFeatureComponentsList(data *schema.ResourceData) []string {
	componentsList := []string{}
	componentsResourceList := data.Get(SchemaComponents).([]interface{})
	for _, componentID := range componentsResourceList {
		componentsList = append(componentsList, componentID.(string))
	}
	return componentsList
}

func getFeatureFlags(data *schema.ResourceData) map[string]string {
	featureFlags := map[string]string{}
	flagsResource := data.Get(SchemaFeatureFlags).(map[string]interface{})
	for name, value := range flagsResource {
		featureFlags[name] = value.(string)
	}
	return featureFlags
}

func cloudauthAccountFeatureFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountFeatureSecure {
	cloudAccountFeature := &v2.CloudauthAccountFeatureSecure{
		AccountFeature: cloudauth.AccountFeature{
			Type:       cloudauth.Feature(cloudauth.Feature_value[data.Get(SchemaType).(string)]),
			Enabled:    data.Get(SchemaEnabled).(bool),
			Components: getFeatureComponentsList(data),
			Flags:      getFeatureFlags(data),
		},
	}

	return cloudAccountFeature
}

func cloudauthAccountFeatureToResourceData(data *schema.ResourceData, cloudAccountFeature *v2.CloudauthAccountFeatureSecure) error {

	accountId := data.Get(SchemaAccountId).(string)
	data.SetId(accountId + "/" + cloudAccountFeature.GetType().String())

	err := data.Set(SchemaAccountId, accountId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaType, cloudAccountFeature.GetType().String())
	if err != nil {
		return err
	}

	err = data.Set(SchemaEnabled, cloudAccountFeature.GetEnabled())
	if err != nil {
		return err
	}

	err = data.Set(SchemaComponents, cloudAccountFeature.GetComponents())
	if err != nil {
		return err
	}

	err = data.Set(SchemaFeatureFlags, cloudAccountFeature.GetFlags())
	if err != nil {
		return err
	}

	return nil
}

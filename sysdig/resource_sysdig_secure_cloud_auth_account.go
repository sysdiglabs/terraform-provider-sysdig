package sysdig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

/*
declare common schemas used across resources here
*/
var (
	accountComponent = &schema.Resource{
		Schema: map[string]*schema.Schema{
			SchemaType: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaInstance: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaVersion: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaCloudConnectorMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaTrustedRoleMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaEventBridgeMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaServicePrincipalMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaWebhookDatasourceMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaCryptoKeyMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaCloudLogsMetadata: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
	}
)

func resourceSysdigSecureCloudauthAccount() *schema.Resource {
	timeout := 5 * time.Minute

	accountFeature := &schema.Resource{
		Schema: map[string]*schema.Schema{
			SchemaType: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaEnabled: {
				Type:     schema.TypeBool,
				Required: true,
			},
			SchemaComponents: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	accountFeatures := &schema.Resource{
		Schema: map[string]*schema.Schema{
			SchemaSecureConfigPosture: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			SchemaSecureIdentityEntitlement: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			SchemaSecureThreatDetection: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			SchemaSecureAgentlessScanning: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			SchemaMonitorCloudMetrics: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
		},
	}

	return &schema.Resource{
		CreateContext: resourceSysdigSecureCloudauthAccountCreate,
		UpdateContext: resourceSysdigSecureCloudauthAccountUpdate,
		ReadContext:   resourceSysdigSecureCloudauthAccountRead,
		DeleteContext: resourceSysdigSecureCloudauthAccountDelete,
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
			SchemaCloudProviderId: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaCloudProviderType: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String(), cloudauth.Provider_PROVIDER_ORACLECLOUD.String()}, false),
			},
			SchemaEnabled: {
				Type:     schema.TypeBool,
				Required: true,
			},
			SchemaFeature: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeatures,
			},
			SchemaComponent: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountComponent,
			},
			SchemaOrganizationIDKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaCloudProviderTenantId: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaCloudProviderAlias: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaProviderPartition: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func getSecureCloudauthAccountClient(client SysdigClients) (v2.CloudauthAccountSecureInterface, error) {
	return client.sysdigSecureClientV2()
}

func resourceSysdigSecureCloudauthAccountCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient((meta.(SysdigClients)))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccount, errStatus, err := client.CreateCloudauthAccountSecure(ctx, cloudauthAccountFromResourceData(data))
	if err != nil {
		return diag.Errorf("Error creating resource: %s %s", errStatus, err)
	}

	data.SetId(cloudauthAccount.Id)
	err = data.Set(SchemaOrganizationIDKey, cloudauthAccount.OrganizationId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccount, errStatus, err := client.GetCloudauthAccountSecure(ctx, data.Id())
	if err != nil {
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	err = cloudauthAccountToResourceData(data, cloudauthAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	existingCloudAccount, errStatus, err := client.GetCloudauthAccountSecure(ctx, data.Id())
	if err != nil {
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	newCloudAccount := cloudauthAccountFromResourceData(data)

	// validate and reject non-updatable resource schema fields upfront
	err = validateCloudauthAccountUpdate(existingCloudAccount, newCloudAccount)
	if err != nil {
		return diag.Errorf("Error updating resource: %s", err)
	}

	_, errStatus, err = client.UpdateCloudauthAccountSecure(ctx, data.Id(), newCloudAccount)
	if err != nil {
		return diag.Errorf("Error updating resource: %s %s", errStatus, err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	errStatus, err := client.DeleteCloudauthAccountSecure(ctx, data.Id())
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
func validateCloudauthAccountUpdate(existingCloudAccount *v2.CloudauthAccountSecure, newCloudAccount *v2.CloudauthAccountSecure) error {
	if existingCloudAccount.Enabled != newCloudAccount.Enabled || existingCloudAccount.Provider != newCloudAccount.Provider ||
		existingCloudAccount.ProviderId != newCloudAccount.ProviderId || existingCloudAccount.OrganizationId != newCloudAccount.OrganizationId {
		errorInvalidResourceUpdate := fmt.Sprintf("Bad Request. Updating restricted fields not allowed: %s", []string{"enabled", "provider_type", "provider_id", "organization_id"})
		return errors.New(errorInvalidResourceUpdate)
	}

	return nil
}

/*
This function converts a schema set to map for iteration.
*/
func convertSchemaSetToMap(set *schema.Set) map[string]interface{} {
	result := make(map[string]interface{})

	for _, element := range set.List() {
		if entry, ok := element.(map[string]interface{}); ok {
			for key, value := range entry {
				result[key] = value
			}
		}
	}

	return result
}

/*
This helper function dynamically populates the account features object for account creation
*/
func setAccountFeature(accountFeatures *cloudauth.AccountFeatures, fieldName string, featureType cloudauth.Feature, valueMap map[string]interface{}) {
	target := reflect.ValueOf(accountFeatures).Elem().FieldByName(fieldName)
	target.Elem().FieldByName("Type").Set(reflect.ValueOf(cloudauth.Feature(featureType)))

	for name, value := range valueMap {
		switch name {
		case SchemaEnabled:
			target.Elem().FieldByName("Enabled").SetBool(value.(bool))
		case SchemaComponents:
			for _, componentID := range value.([]interface{}) {
				target.Elem().FieldByName("Components").Set(reflect.Append(target.Elem().FieldByName("Components"), reflect.ValueOf(componentID.(string))))
			}
		}
	}
}

/*
This helper function aggregates the account features object that will be used in the
cloudauthAccountFromResourceData() function
*/
func constructAccountFeatures(data *schema.ResourceData) *cloudauth.AccountFeatures {
	accountFeatures := &cloudauth.AccountFeatures{}
	featureMap := convertSchemaSetToMap(data.Get(SchemaFeature).(*schema.Set))

	for name, value := range featureMap {
		if featureValues := value.(*schema.Set).List(); len(featureValues) > 0 {
			valueMap := convertSchemaSetToMap(value.(*schema.Set))
			switch name {
			case SchemaSecureConfigPosture:
				accountFeatures.SecureConfigPosture = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureConfigPosture", cloudauth.Feature_FEATURE_SECURE_CONFIG_POSTURE, valueMap)
			case SchemaSecureIdentityEntitlement:
				accountFeatures.SecureIdentityEntitlement = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureIdentityEntitlement", cloudauth.Feature_FEATURE_SECURE_IDENTITY_ENTITLEMENT, valueMap)
			case SchemaSecureThreatDetection:
				accountFeatures.SecureThreatDetection = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureThreatDetection", cloudauth.Feature_FEATURE_SECURE_THREAT_DETECTION, valueMap)
			case SchemaSecureAgentlessScanning:
				accountFeatures.SecureAgentlessScanning = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureAgentlessScanning", cloudauth.Feature_FEATURE_SECURE_AGENTLESS_SCANNING, valueMap)
			case SchemaMonitorCloudMetrics:
				accountFeatures.MonitorCloudMetrics = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "MonitorCloudMetrics", cloudauth.Feature_FEATURE_MONITOR_CLOUD_METRICS, valueMap)
			}
		}
	}

	return accountFeatures
}

/*
This helper function aggregates the account components list that will be used in the
cloudauthAccountFromResourceData() function
*/
func constructAccountComponents(data *schema.ResourceData) []*cloudauth.AccountComponent {
	accountComponents := []*cloudauth.AccountComponent{}

	for _, rc := range data.Get(SchemaComponent).(*schema.Set).List() {
		resourceComponent := rc.(map[string]interface{})
		component := &cloudauth.AccountComponent{}
		var err error

		for key, value := range resourceComponent {
			if value != nil && value.(string) != "" {
				switch key {
				case SchemaType:
					component.Type = cloudauth.Component(cloudauth.Component_value[value.(string)])
				case SchemaInstance:
					component.Instance = value.(string)
				case SchemaVersion:
					component.Version = value.(string)
				case SchemaCloudConnectorMetadata:
					component.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{CloudConnectorMetadata: &cloudauth.CloudConnectorMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetCloudConnectorMetadata())
				case SchemaTrustedRoleMetadata:
					component.Metadata = &cloudauth.AccountComponent_TrustedRoleMetadata{TrustedRoleMetadata: &cloudauth.TrustedRoleMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetTrustedRoleMetadata())
				case SchemaEventBridgeMetadata:
					component.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{EventBridgeMetadata: &cloudauth.EventBridgeMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetEventBridgeMetadata())
				case SchemaServicePrincipalMetadata:
					component.Metadata = &cloudauth.AccountComponent_ServicePrincipalMetadata{ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetServicePrincipalMetadata())
					// special handling for GCP service principal if it has keys which are base64 encoded
					if data.Get(SchemaCloudProviderType).(string) == cloudauth.Provider_PROVIDER_GCP.String() {
						component.Metadata = constructGcpServicePrincipalMetadata(value.(string))
					}
				case SchemaWebhookDatasourceMetadata:
					component.Metadata = &cloudauth.AccountComponent_WebhookDatasourceMetadata{WebhookDatasourceMetadata: &cloudauth.WebhookDatasourceMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetWebhookDatasourceMetadata())
				case SchemaCryptoKeyMetadata:
					component.Metadata = &cloudauth.AccountComponent_CryptoKeyMetadata{CryptoKeyMetadata: &cloudauth.CryptoKeyMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetCryptoKeyMetadata())
				case SchemaCloudLogsMetadata:
					component.Metadata = &cloudauth.AccountComponent_CloudLogsMetadata{CloudLogsMetadata: &cloudauth.CloudLogsMetadata{}}
					err = protojson.Unmarshal([]byte(value.(string)), component.GetCloudLogsMetadata())
				}
				if err != nil {
					diag.FromErr(err)
				}
			}
		}
		accountComponents = append(accountComponents, component)
	}

	return accountComponents
}

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountSecure {
	return &v2.CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:           data.Get(SchemaEnabled).(bool),
			OrganizationId:    data.Get(SchemaOrganizationIDKey).(string),
			ProviderId:        data.Get(SchemaCloudProviderId).(string),
			Provider:          cloudauth.Provider(cloudauth.Provider_value[data.Get(SchemaCloudProviderType).(string)]),
			Components:        constructAccountComponents(data),
			Feature:           constructAccountFeatures(data),
			ProviderTenantId:  data.Get(SchemaCloudProviderTenantId).(string),
			ProviderAlias:     data.Get(SchemaCloudProviderAlias).(string),
			ProviderPartition: cloudauth.ProviderPartition(cloudauth.ProviderPartition_value[data.Get(SchemaProviderPartition).(string)]),
		},
	}
}

/*
This helper function converts feature values from *cloudauth.AccountFeature to resource data schema.
*/
func featureValuesToResourceData(feature *cloudauth.AccountFeature) map[string]interface{} {
	valuesMap := make(map[string]interface{})

	valuesMap["type"] = feature.Type.String()
	valuesMap["enabled"] = feature.Enabled
	valuesMap["components"] = feature.Components

	return valuesMap
}

/*
This helper function converts the features data from *cloudauth.AccountFeatures to resource data schema.
This is needed to set the value in cloudauthAccountToResourceData().
*/
func featureToResourceData(features *cloudauth.AccountFeatures) []interface{} {
	// In the resource data, SchemaFeature field is a nested set[] of sets[] of individual features
	// Hence we need to return this uber level set[] to cloudauthAccountToResourceData
	featureMap := []interface{}{}

	featureFields := map[string]*cloudauth.AccountFeature{
		SchemaSecureThreatDetection:     features.SecureThreatDetection,
		SchemaSecureConfigPosture:       features.SecureConfigPosture,
		SchemaSecureIdentityEntitlement: features.SecureIdentityEntitlement,
		SchemaMonitorCloudMetrics:       features.MonitorCloudMetrics,
		SchemaSecureAgentlessScanning:   features.SecureAgentlessScanning,
	}

	allFeatures := make(map[string]interface{})
	for name, feature := range featureFields {
		if feature != nil {
			featureBlock := make([]map[string]interface{}, 0)
			value := featureValuesToResourceData(feature)
			featureBlock = append(featureBlock, value)

			allFeatures[name] = featureBlock
		}
	}

	// return featureMap only if there is any features data from *cloudauth.AccountFeatures, else return nil
	if len(allFeatures) > 0 {
		featureMap = append(featureMap, allFeatures)
		return featureMap
	}
	return nil
}

/*
This helper function converts the components data from []*cloudauth.AccountComponent to resource data schema.
This is needed to set the value in cloudauthAccountToResourceData().
*/
func componentsToResourceData(components []*cloudauth.AccountComponent) []map[string]interface{} {
	resourceList := []map[string]interface{}{}
	for _, component := range components {
		resourceData := map[string]interface{}{}
		resourceData[SchemaType] = component.GetType().String()
		resourceData[SchemaInstance] = component.GetInstance()
		resourceData[SchemaVersion] = component.GetVersion()

		switch component.GetType() {
		case cloudauth.Component_COMPONENT_CLOUD_CONNECTOR:
			resourceData[SchemaCloudConnectorMetadata] = getComponentMetadataString(component.GetCloudConnectorMetadata())
		case cloudauth.Component_COMPONENT_TRUSTED_ROLE:
			resourceData[SchemaTrustedRoleMetadata] = getComponentMetadataString(component.GetTrustedRoleMetadata())
		case cloudauth.Component_COMPONENT_EVENT_BRIDGE:
			resourceData[SchemaEventBridgeMetadata] = getComponentMetadataString(component.GetEventBridgeMetadata())
		case cloudauth.Component_COMPONENT_SERVICE_PRINCIPAL:
			// special handling for GCP service principal if it has keys which are to be base64 encoded
			if component.GetServicePrincipalMetadata().GetGcp() != nil {
				resourceData[SchemaServicePrincipalMetadata] = getGcpServicePrincipalMetadata(component.GetServicePrincipalMetadata())
			} else {
				resourceData[SchemaServicePrincipalMetadata] = getComponentMetadataString(component.GetServicePrincipalMetadata())
			}
		case cloudauth.Component_COMPONENT_WEBHOOK_DATASOURCE:
			resourceData[SchemaWebhookDatasourceMetadata] = getComponentMetadataString(component.GetWebhookDatasourceMetadata())
		case cloudauth.Component_COMPONENT_CRYPTO_KEY:
			resourceData[SchemaCryptoKeyMetadata] = getComponentMetadataString(component.GetCryptoKeyMetadata())
		case cloudauth.Component_COMPONENT_CLOUD_LOGS:
			resourceData[SchemaCloudLogsMetadata] = getComponentMetadataString(component.GetCloudLogsMetadata())
		}
		resourceList = append(resourceList, resourceData)
	}

	return resourceList
}

func getComponentMetadataString(message protoreflect.ProtoMessage) string {
	// marshal through protojson get correct snake case keys
	protoJsonMessage, err := protojson.MarshalOptions{UseProtoNames: true}.Marshal(message)
	if err != nil {
		diag.FromErr(err)
	}
	// re-marshal through encoding/json to get consistent key ordering, avoiding diff errors with TF internals
	metadataMap := make(map[string]interface{})
	err = json.Unmarshal(protoJsonMessage, &metadataMap)
	if err != nil {
		diag.FromErr(err)
	}
	jsonMessage, err := json.Marshal(metadataMap)
	if err != nil {
		diag.FromErr(err)
	}
	return string(jsonMessage)
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	err := data.Set(SchemaIDKey, cloudAccount.Id)
	if err != nil {
		return err
	}

	err = data.Set(SchemaEnabled, cloudAccount.Enabled)
	if err != nil {
		return err
	}

	err = data.Set(SchemaCloudProviderId, cloudAccount.ProviderId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaCloudProviderType, cloudAccount.Provider.String())
	if err != nil {
		return err
	}

	err = data.Set(SchemaFeature, featureToResourceData(cloudAccount.Feature))
	if err != nil {
		return err
	}

	err = data.Set(SchemaComponent, componentsToResourceData(cloudAccount.GetComponents()))
	if err != nil {
		return err
	}

	err = data.Set(SchemaOrganizationIDKey, cloudAccount.OrganizationId)
	if err != nil {
		return err
	}

	if cloudAccount.Provider == cloudauth.Provider_PROVIDER_AZURE {
		err = data.Set(SchemaCloudProviderTenantId, cloudAccount.ProviderTenantId)
		if err != nil {
			return err
		}

		err = data.Set(SchemaCloudProviderAlias, cloudAccount.ProviderAlias)
		if err != nil {
			return err
		}
	}

	if cloudAccount.Provider == cloudauth.Provider_PROVIDER_ORACLECLOUD {
		err = data.Set(SchemaCloudProviderTenantId, cloudAccount.ProviderTenantId)
		if err != nil {
			return err
		}
	}

	if !(cloudAccount.ProviderPartition.String() == cloudauth.ProviderPartition_PROVIDER_PARTITION_UNSPECIFIED.String()) {
		err = data.Set(SchemaProviderPartition, cloudAccount.ProviderPartition.String())
		if err != nil {
			return err
		}
	}

	return nil
}

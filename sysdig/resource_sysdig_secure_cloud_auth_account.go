package sysdig

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
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

	accountComponents := &schema.Resource{
		Schema: map[string]*schema.Schema{
			SchemaType: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaInstance: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaCloudConnectorMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaTrustedRoleMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaEventBridgeMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaServicePrincipalMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaWebhookDatasourceMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaCryptoKeyMetadata: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaCloudLogsMetadata: {
				Type:     schema.TypeString,
				Optional: true,
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
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String()}, false),
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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     accountComponents,
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
		if strings.Contains(errStatus, "404") {
			return nil
		}
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
		if strings.Contains(errStatus, "404") {
			return nil
		}
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
		if strings.Contains(errStatus, "404") {
			return nil
		}
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
func constructAccountFeatures(accountFeatures *cloudauth.AccountFeatures, featureData interface{}) *cloudauth.AccountFeatures {
	featureMap := convertSchemaSetToMap(featureData.(*schema.Set))

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
func constructAccountComponents(accountComponents []*cloudauth.AccountComponent, data *schema.ResourceData) []*cloudauth.AccountComponent {
	provider := data.Get(SchemaCloudProviderType).(string)

	for _, rc := range data.Get(SchemaComponent).([]interface{}) {
		resourceComponent := rc.(map[string]interface{})
		component := &cloudauth.AccountComponent{}

		for key, value := range resourceComponent {
			if value != nil && value.(string) != "" {
				switch key {
				case SchemaType:
					component.Type = cloudauth.Component(cloudauth.Component_value[value.(string)])
				case SchemaInstance:
					component.Instance = value.(string)
				case SchemaCloudConnectorMetadata:
					component.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{
						CloudConnectorMetadata: &cloudauth.CloudConnectorMetadata{},
					}
				case SchemaTrustedRoleMetadata:
					// TODO: Make it more generic than just for GCP
					if provider == cloudauth.Provider_PROVIDER_GCP.String() {
						component.Metadata = &cloudauth.AccountComponent_TrustedRoleMetadata{
							TrustedRoleMetadata: &cloudauth.TrustedRoleMetadata{
								Provider: &cloudauth.TrustedRoleMetadata_Gcp{
									Gcp: &cloudauth.TrustedRoleMetadata_GCP{
										RoleName: value.(string),
									},
								},
							},
						}
					}
				case SchemaEventBridgeMetadata:
					if provider == cloudauth.Provider_PROVIDER_AZURE.String() {
						eventBridgeMetadata := parseResourceMetadataJson(value.(string))
						eventHubMetadata, ok := eventBridgeMetadata["azure"].(map[string]interface{})["event_hub_metadata"].(map[string]interface{})

						if !ok {
							fmt.Printf("Resource input for component metadata for provider %s is invalid and not as expected", provider)
							break
						}

						component.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{
							EventBridgeMetadata: &cloudauth.EventBridgeMetadata{
								Provider: &cloudauth.EventBridgeMetadata_Azure_{
									Azure: &cloudauth.EventBridgeMetadata_Azure{
										EventHubMetadata: &cloudauth.EventBridgeMetadata_Azure_EventHubMetadata{
											EventHubName:      eventHubMetadata["event_hub_name"].(string),
											EventHubNamespace: eventHubMetadata["event_hub_namespace"].(string),
											ConsumerGroup:     eventHubMetadata["consumer_group"].(string),
										},
									},
								},
							},
						}
					} else {
						component.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{
							EventBridgeMetadata: &cloudauth.EventBridgeMetadata{},
						}
					}
				case SchemaServicePrincipalMetadata:
					// TODO: Make it more generic than just for GCP
					servicePrincipalMetadata := parseResourceMetadataJson(value.(string))

					if provider == cloudauth.Provider_PROVIDER_GCP.String() {
						if metadata, err := getServicePrincipalMetadataForGCP(servicePrincipalMetadata); err == nil {
							component.Metadata = metadata
						}
					} else if provider == cloudauth.Provider_PROVIDER_AZURE.String() {
						servicePrincipalAzureKey, ok := servicePrincipalMetadata["azure"].(map[string]interface{})["active_directory_service_principal"].(map[string]interface{})

						if !ok {
							fmt.Printf("Resource input for component metadata for provider %s is invalid and not as expected", provider)
							break
						}

						component.Metadata = &cloudauth.AccountComponent_ServicePrincipalMetadata{
							ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{
								Provider: &cloudauth.ServicePrincipalMetadata_Azure_{
									Azure: &cloudauth.ServicePrincipalMetadata_Azure{
										ActiveDirectoryServicePrincipal: &cloudauth.ServicePrincipalMetadata_Azure_ActiveDirectoryServicePrincipal{
											AccountEnabled:         servicePrincipalAzureKey["account_enabled"].(bool),
											AppDisplayName:         servicePrincipalAzureKey["app_display_name"].(string),
											AppId:                  servicePrincipalAzureKey["app_id"].(string),
											AppOwnerOrganizationId: servicePrincipalAzureKey["app_owner_organization_id"].(string),
											DisplayName:            servicePrincipalAzureKey["display_name"].(string),
											Id:                     servicePrincipalAzureKey["id"].(string),
										},
									},
								},
							},
						}
					}
				case SchemaWebhookDatasourceMetadata:
					component.Metadata = &cloudauth.AccountComponent_WebhookDatasourceMetadata{
						WebhookDatasourceMetadata: &cloudauth.WebhookDatasourceMetadata{},
					}
				case SchemaCryptoKeyMetadata:
					component.Metadata = &cloudauth.AccountComponent_CryptoKeyMetadata{
						CryptoKeyMetadata: &cloudauth.CryptoKeyMetadata{},
					}
				case SchemaCloudLogsMetadata:
					component.Metadata = &cloudauth.AccountComponent_CloudLogsMetadata{
						CloudLogsMetadata: &cloudauth.CloudLogsMetadata{},
					}
				}
			}
		}
		accountComponents = append(accountComponents, component)
	}

	return accountComponents
}

func getServicePrincipalMetadataForGCP(servicePrincipalMetadata map[string]interface{}) (*cloudauth.AccountComponent_ServicePrincipalMetadata, error) {
	encodedServicePrincipal, ok := servicePrincipalMetadata["gcp"].(map[string]interface{})
	if !ok {
		fmt.Printf("Failed to parse service principal metadata, for provider GCP")
		return nil, fmt.Errorf("failed to parse service principal metadata, for provider GCP")
	}

	rawKeyData, okKey := encodedServicePrincipal["key"].(string)
	wifData, okWif := encodedServicePrincipal["workload_identity_federation"].(map[string]interface{})

	if !okKey && !okWif {
		fmt.Printf("Failed to parse on of either service principal types, key or workload_identity_federation")
		return nil, fmt.Errorf("failed to parse on of either service principal types, key or workload_identity_federation")
	}

	servcicePrincipalMetadata := &cloudauth.AccountComponent_ServicePrincipalMetadata{
		ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{
			Provider: &cloudauth.ServicePrincipalMetadata_Gcp{
				Gcp: &cloudauth.ServicePrincipalMetadata_GCP{},
			},
		},
	}

	if okKey {
		keyData := decodeServicePrincipalKeyToMap(rawKeyData)
		servcicePrincipalMetadata.ServicePrincipalMetadata.GetGcp().Key = &cloudauth.ServicePrincipalMetadata_GCP_Key{
			Type:                    keyData["type"],
			ProjectId:               keyData["project_id"],
			PrivateKeyId:            keyData["private_key_id"],
			PrivateKey:              keyData["private_key"],
			ClientEmail:             keyData["client_email"],
			ClientId:                keyData["client_id"],
			AuthUri:                 keyData["auth_uri"],
			TokenUri:                keyData["token_uri"],
			AuthProviderX509CertUrl: keyData["auth_provider_x509_cert_url"],
			ClientX509CertUrl:       keyData["client_x509_cert_url"],
		}
	}

	if okWif {
		servcicePrincipalMetadata.ServicePrincipalMetadata.GetGcp().WorkloadIdentityFederation = &cloudauth.ServicePrincipalMetadata_GCP_WorkloadIdentityFederation{
			PoolProviderId: wifData["pool_provider_id"].(string),
		}
		servcicePrincipalMetadata.ServicePrincipalMetadata.GetGcp().Email = encodedServicePrincipal["email"].(string)
	}

	return servcicePrincipalMetadata, nil
}

/*
This helper function parses the provided component resource metadata in opaque Json string format into a map
*/
func parseResourceMetadataJson(value string) map[string]interface{} {
	var metadataJSON map[string]interface{}
	err := json.Unmarshal([]byte(value), &metadataJSON)
	if err != nil {
		fmt.Printf("Failed to parse component metadata: %v", err)
		return nil
	}

	return metadataJSON
}

/*
This helper function decodes the base64 encoded Service Principal Key obtained from cloud
and parses it from Json format into a map
*/
func decodeServicePrincipalKeyToMap(encodedKey string) map[string]string {
	bytesData, err := b64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		fmt.Printf("Failed to decode service principal key: %v", err)
		return nil
	}
	var privateKeyMap map[string]string
	err = json.Unmarshal(bytesData, &privateKeyMap)
	if err != nil {
		fmt.Printf("Failed to parse service principal key: %v", err)
		return nil
	}

	return privateKeyMap
}

/*
This helper function encodes the Service Principal Key returned by Sysdig
and returns a base64 encoded string
*/
func encodeServicePrincipalKey(key []byte) string {
	encodedKey := b64.StdEncoding.EncodeToString(key)
	return encodedKey
}

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountSecure {
	accountComponents := constructAccountComponents([]*cloudauth.AccountComponent{}, data)

	featureData := data.Get(SchemaFeature)
	accountFeatures := constructAccountFeatures(&cloudauth.AccountFeatures{}, featureData)

	return &v2.CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:          data.Get(SchemaEnabled).(bool),
			OrganizationId:   data.Get(SchemaOrganizationIDKey).(string),
			ProviderId:       data.Get(SchemaCloudProviderId).(string),
			Provider:         cloudauth.Provider(cloudauth.Provider_value[data.Get(SchemaCloudProviderType).(string)]),
			Components:       accountComponents,
			Feature:          accountFeatures,
			ProviderTenantId: data.Get(SchemaCloudProviderTenantId).(string),
			ProviderAlias:    data.Get(SchemaCloudProviderAlias).(string),
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
func componentsToResourceData(components []*cloudauth.AccountComponent, dataComponentsOrder []string) []map[string]interface{} {
	// In the resource data, SchemaComponent field is a list of component sets[] / block
	// Hence we need to return this uber level list in same order to cloudauthAccountToResourceData
	componentsList := []map[string]interface{}{}

	allComponents := make(map[string]interface{})
	for _, comp := range components {
		componentBlock := map[string]interface{}{}

		componentBlock[SchemaType] = comp.Type.String()
		componentBlock[SchemaInstance] = comp.Instance

		metadata := comp.GetMetadata()
		if metadata != nil {
			switch metadata.(type) {
			case *cloudauth.AccountComponent_ServicePrincipalMetadata:
				provider := metadata.(*cloudauth.AccountComponent_ServicePrincipalMetadata).ServicePrincipalMetadata.GetProvider()
				// TODO: Make it more generic than just for GCP
				if servicePrincipalMetadata, ok := provider.(*cloudauth.ServicePrincipalMetadata_Gcp); ok {

					metadataContent := map[string]interface{}{
						"gcp": map[string]interface{}{},
					}

					gcpKey := servicePrincipalMetadata.Gcp.GetKey()
					if gcpKey != nil {
						keyBytes, err := deserializeServiceMetadata_GCP_Key(servicePrincipalMetadata)
						if err != nil {
							fmt.Printf("Failed to deserializeServiceMetadata_GCP_Key %s: %v", SchemaServicePrincipalMetadata, err)
							break
						}
						metadataContent["gcp"].(map[string]interface{})["key"] = keyBytes
					}

					workloadIdentityData := servicePrincipalMetadata.Gcp.GetWorkloadIdentityFederation()
					if workloadIdentityData != nil {
						metadataContent["gcp"].(map[string]interface{})["workload_identity_federation"] = map[string]interface{}{
							"pool_provider_id": workloadIdentityData.GetPoolProviderId(),
						}

						metadataContent["gcp"].(map[string]interface{})["email"] = servicePrincipalMetadata.Gcp.GetEmail()
					}

					// encode the key to base64 and add to the component block
					schemaData, err := json.Marshal(metadataContent)
					if err != nil {
						fmt.Printf("Failed to populate %s: %v", SchemaServicePrincipalMetadata, err)
						break
					}

					componentBlock[SchemaServicePrincipalMetadata] = string(schemaData)
				}

				if providerKey, ok := provider.(*cloudauth.ServicePrincipalMetadata_Azure_); ok {
					schemaData, err := json.Marshal(map[string]interface{}{
						"azure": map[string]interface{}{
							"active_directory_service_principal": map[string]interface{}{
								"account_enabled":           providerKey.Azure.GetActiveDirectoryServicePrincipal().GetAccountEnabled(),
								"app_display_name":          providerKey.Azure.GetActiveDirectoryServicePrincipal().GetAppDisplayName(),
								"app_id":                    providerKey.Azure.GetActiveDirectoryServicePrincipal().GetAppId(),
								"app_owner_organization_id": providerKey.Azure.GetActiveDirectoryServicePrincipal().GetAppOwnerOrganizationId(),
								"display_name":              providerKey.Azure.GetActiveDirectoryServicePrincipal().GetDisplayName(),
								"id":                        providerKey.Azure.GetActiveDirectoryServicePrincipal().GetId(),
							},
						},
					})
					if err != nil {
						fmt.Printf("Failed to populate %s: %v", SchemaServicePrincipalMetadata, err)
						break
					}

					componentBlock[SchemaServicePrincipalMetadata] = string(schemaData)
				}
			case *cloudauth.AccountComponent_EventBridgeMetadata:
				provider := metadata.(*cloudauth.AccountComponent_EventBridgeMetadata).EventBridgeMetadata.GetProvider()

				if providerKey, ok := provider.(*cloudauth.EventBridgeMetadata_Azure_); ok {
					schemaData, err := json.Marshal(map[string]interface{}{
						"azure": map[string]interface{}{
							"event_hub_metadata": map[string]interface{}{
								"event_hub_name":      providerKey.Azure.GetEventHubMetadata().GetEventHubName(),
								"event_hub_namespace": providerKey.Azure.GetEventHubMetadata().GetEventHubNamespace(),
								"consumer_group":      providerKey.Azure.GetEventHubMetadata().GetConsumerGroup(),
							},
						},
					})
					if err != nil {
						fmt.Printf("Failed to populate %s: %v", SchemaEventBridgeMetadata, err)
						break
					}

					componentBlock[SchemaEventBridgeMetadata] = string(schemaData)
				}
			}
		}

		allComponents[comp.Instance] = componentBlock
	}

	// return componentsList only if there is any components data from *[]cloudauth.AccountComponent, else return nil
	if len(allComponents) > 0 {
		// add the component blocks in same order to maintain ordering
		for _, c := range dataComponentsOrder {
			componentItem := allComponents[c].(map[string]interface{})
			componentsList = append(componentsList, componentItem)
		}
		return componentsList
	}

	return nil
}

func deserializeServiceMetadata_GCP_Key(servicePrincipalMetadata *cloudauth.ServicePrincipalMetadata_Gcp) (string, error) {
	// convert key struct to jsonified key with all the expected fields
	jsonifiedKey := struct {
		Type                    string `json:"type"`
		ProjectId               string `json:"project_id"`
		PrivateKeyId            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientId                string `json:"client_id"`
		AuthUri                 string `json:"auth_uri"`
		TokenUri                string `json:"token_uri"`
		AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
		ClientX509CertUrl       string `json:"client_x509_cert_url"`
		UniverseDomain          string `json:"universe_domain"`
	}{
		Type:                    servicePrincipalMetadata.Gcp.GetKey().GetType(),
		ProjectId:               servicePrincipalMetadata.Gcp.GetKey().GetProjectId(),
		PrivateKeyId:            servicePrincipalMetadata.Gcp.GetKey().GetPrivateKeyId(),
		PrivateKey:              servicePrincipalMetadata.Gcp.GetKey().GetPrivateKey(),
		ClientEmail:             servicePrincipalMetadata.Gcp.GetKey().GetClientEmail(),
		ClientId:                servicePrincipalMetadata.Gcp.GetKey().GetClientId(),
		AuthUri:                 servicePrincipalMetadata.Gcp.GetKey().GetAuthUri(),
		TokenUri:                servicePrincipalMetadata.Gcp.GetKey().GetTokenUri(),
		AuthProviderX509CertUrl: servicePrincipalMetadata.Gcp.GetKey().GetAuthProviderX509CertUrl(),
		ClientX509CertUrl:       servicePrincipalMetadata.Gcp.GetKey().GetClientX509CertUrl(),
		UniverseDomain:          "googleapis.com",
	}
	bytesKey, err := json.Marshal(jsonifiedKey)
	if err != nil {
		return "", fmt.Errorf("failed to populate %s: %v", SchemaServicePrincipalMetadata, err)
	}

	// update the json with proper indentation
	var keyOut bytes.Buffer
	if err := json.Indent(&keyOut, bytesKey, "", "  "); err != nil {
		return "", fmt.Errorf("failed to populate %s: %v", SchemaServicePrincipalMetadata, err)
	}
	keyOut.WriteByte('\n')

	keyBytes := encodeServicePrincipalKey(keyOut.Bytes())
	return keyBytes, nil
}

func getResourceComponentsOrder(dataComponents interface{}) []string {
	var dataComponentsOrder []string
	for _, rc := range dataComponents.([]interface{}) {
		resourceComponent := rc.(map[string]interface{})
		dataComponentsOrder = append(dataComponentsOrder, resourceComponent[SchemaInstance].(string))
	}
	return dataComponentsOrder
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

	dataComponentsOrder := getResourceComponentsOrder(data.Get(SchemaComponent))
	err = data.Set(SchemaComponent, componentsToResourceData(cloudAccount.Components, dataComponentsOrder))
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

	return nil
}

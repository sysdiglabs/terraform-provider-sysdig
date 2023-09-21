package sysdig

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCloudauthAccount() *schema.Resource {
	timeout := 5 * time.Minute

	var accountFeature = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	var accountFeatures = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"secure_config_posture": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			"secure_identity_entitlement": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			"secure_threat_detection": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			"secure_agentless_scanning": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
			"monitor_cloud_metrics": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeature,
			},
		},
	}

	var accountComponents = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_connector_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trusted_role_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_bridge_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_principal_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"webhook_datasource_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"crypto_key_metadata": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_logs_metadata": {
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
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{cloudauth.Provider_PROVIDER_AWS.String(), cloudauth.Provider_PROVIDER_GCP.String(), cloudauth.Provider_PROVIDER_AZURE.String()}, false),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"feature": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     accountFeatures,
			},
			"components": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     accountComponents,
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

	cloudauthAccount, err := client.CreateCloudauthAccountSecure(ctx, cloudauthAccountFromResourceData(data))
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(cloudauthAccount.Id)

	return nil
}

func resourceSysdigSecureCloudauthAccountRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccount, err := client.GetCloudauthAccountSecure(ctx, data.Id())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
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

	_, err = client.UpdateCloudauthAccountSecure(ctx, data.Id(), cloudauthAccountFromResourceData(data))

	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudauthAccountSecure(ctx, data.Id())

	if err != nil {
		return diag.FromErr(err)
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
		case "enabled":
			target.Elem().FieldByName("Enabled").SetBool(value.(bool))
		case "components":
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
			case "secure_config_posture":
				accountFeatures.SecureConfigPosture = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureConfigPosture", cloudauth.Feature_FEATURE_SECURE_CONFIG_POSTURE, valueMap)
			case "secure_identity_entitlement":
				accountFeatures.SecureIdentityEntitlement = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureIdentityEntitlement", cloudauth.Feature_FEATURE_SECURE_IDENTITY_ENTITLEMENT, valueMap)
			case "secure_threat_detection":
				accountFeatures.SecureThreatDetection = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureThreatDetection", cloudauth.Feature_FEATURE_SECURE_THREAT_DETECTION, valueMap)
			case "secure_agentless_scanning":
				accountFeatures.SecureAgentlessScanning = &cloudauth.AccountFeature{}
				setAccountFeature(accountFeatures, "SecureAgentlessScanning", cloudauth.Feature_FEATURE_SECURE_AGENTLESS_SCANNING, valueMap)
			case "monitor_cloud_metrics":
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
	provider := data.Get("cloud_provider_type").(string)

	for _, rc := range data.Get("components").([]interface{}) {
		resourceComponent := rc.(map[string]interface{})
		component := &cloudauth.AccountComponent{}

		for key, value := range resourceComponent {
			if value != nil && value.(string) != "" {
				switch key {
				case "type":
					component.Type = cloudauth.Component(cloudauth.Component_value[value.(string)])
				case "instance":
					component.Instance = value.(string)
				case "cloud_connector_metadata":
					component.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{
						CloudConnectorMetadata: &cloudauth.CloudConnectorMetadata{},
					}
				case "trusted_role_metadata":
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
				case "event_bridge_metadata":
					component.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{
						EventBridgeMetadata: &cloudauth.EventBridgeMetadata{},
					}
				case "service_principal_metadata":
					// TODO: Make it more generic than just for GCP
					// service_principal_private_key := getServicePrincipalKeyObject(value.(string))
					component.Metadata = &cloudauth.AccountComponent_ServicePrincipalMetadata{
						ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{
							Provider: &cloudauth.ServicePrincipalMetadata_Gcp{
								Gcp: &cloudauth.ServicePrincipalMetadata_GCP{
									Key: &cloudauth.ServicePrincipalMetadata_GCP_Key{
										ProjectId:    data.Get("cloud_provider_id").(string),
										PrivateKeyId: "testKeyID",
										PrivateKey:   "testKey",
										// PrivateKeyId: service_principal_private_key["private_key_id"],
										// PrivateKey:   service_principal_private_key["private_key"],
									},
								},
							},
						},
					}
				case "webhook_datasource_metadata":
					component.Metadata = &cloudauth.AccountComponent_WebhookDatasourceMetadata{
						WebhookDatasourceMetadata: &cloudauth.WebhookDatasourceMetadata{},
					}
				case "crypto_key_metadata":
					component.Metadata = &cloudauth.AccountComponent_CryptoKeyMetadata{
						CryptoKeyMetadata: &cloudauth.CryptoKeyMetadata{},
					}
				case "cloud_logs_metadata":
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

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountSecure {
	accountComponents := constructAccountComponents([]*cloudauth.AccountComponent{}, data)

	featureData := data.Get("feature").(interface{})
	accountFeatures := constructAccountFeatures(&cloudauth.AccountFeatures{}, featureData)

	return &v2.CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:    data.Get("enabled").(bool),
			ProviderId: data.Get("cloud_provider_id").(string),
			Provider:   cloudauth.Provider(cloudauth.Provider_value[data.Get("cloud_provider_type").(string)]),
			Components: accountComponents,
			Feature:    accountFeatures,
		},
	}
}

func AccountFeatureToMap(feature *cloudauth.AccountFeature) map[string]interface{} {
	log.Printf("DEBUG individual feature => %v", feature)
	log.Printf("DEBUG individual feature components => %v", feature.Components)
	// featureMap := make(map[string]interface{})
	objectMap := make(map[string]interface{})

	objectMap["type"] = feature.Type.String()
	objectMap["enabled"] = feature.Enabled
	objectMap["components"] = feature.Components

	// if feature.CreatedAt != nil {
	// 	objectMap["createdAt"] = feature.CreatedAt.AsTime().Format(time.RFC3339)
	// }

	// featureMap["secure_config_posture"] = objectMap

	log.Printf("DEBUG objectMap => %v", objectMap)

	return objectMap
}

func mapToSet(features *cloudauth.AccountFeatures) map[string]interface{} {
	// featuresMap := featureToResourceData(features)
	// set := schema.NewSet(schema.HashString, []interface{})

	return map[string]interface{}{
		"secure_config_posture": AccountFeatureToMap(features.SecureConfigPosture),
	}

	// if features.SecureThreatDetection != nil {
	// 	set.Add(AccountFeatureToMap(features.SecureThreatDetection))
	// }

	// if features.SecureConfigPosture != nil {
	// 	featureMap := AccountFeatureToMap(features.SecureConfigPosture)
	// 	jsonData, _ := json.Marshal(featureMap)
	// 	log.Printf("DEBUG JSON data => %v", jsonData)
	// 	log.Printf("DEBUG string JSON data => %v", string(jsonData))

	// 	// set.Add(string(jsonData))
	// 	set.Add(featureMap)
	// }

	// if features.SecureIdentityEntitlement != nil {
	// 	set.Add(AccountFeatureToMap(features.SecureIdentityEntitlement))
	// }

	// if features.MonitorCloudMetrics != nil {
	// 	set.Add(AccountFeatureToMap(features.MonitorCloudMetrics))
	// }

	// if features.SecureAgentlessScanning != nil {
	// 	set.Add(AccountFeatureToMap(features.SecureAgentlessScanning))
	// }

	// log.Printf("DEBUG feature set => %v", set)
	// log.Printf("DEBUG feature set type => %T", set)
	// return set
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	for _, err := range []error{
		data.Set("id", cloudAccount.Id),
		data.Set("enabled", cloudAccount.Enabled),
		data.Set("cloud_provider_id", cloudAccount.ProviderId),
		data.Set("cloud_provider_type", cloudAccount.Provider.String()),
		data.Set("components", cloudAccount.Components),
		// data.Set("feature", cloudAccount.Feature),
		data.Set("feature", []map[string]interface{}{
			mapToSet(cloudAccount.Feature),
		}),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

func getServicePrincipalKeyObject(value string) map[string]string {
	bytes, err := b64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Printf("Failed to decode service principal key: %v", err)
		return nil
	}
	var privateKeyJSON map[string]string
	err = json.Unmarshal(bytes, &privateKeyJSON)
	if err != nil {
		fmt.Printf("Failed to parse service principal key: %v", err)
		return nil
	}

	return privateKeyJSON
}

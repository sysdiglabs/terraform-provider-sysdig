package sysdig

import (
	"context"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/protobuf/encoding/protojson"
)

func resourceSysdigSecureCloudauthAccount() *schema.Resource {
	timeout := 5 * time.Minute

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
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cloudConnectorMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"trustedRoleMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"eventBridgeMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"servicePrincipalMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"webhookDatasourceMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cryptoKeyMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cloudLogsMetadata": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

var accountFeatures = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"secureConfigPosture": {
			Type: schema.TypeSet,
			Elem: accountFeature,
		},
		"secureIdentityEntitlement": {
			Type: schema.TypeSet,
			Elem: accountFeature,
		},
		"secureThreatDetection": {
			Type: schema.TypeSet,
			Elem: accountFeature,
		},
		"secureAgentlessScanning": {
			Type: schema.TypeSet,
			Elem: accountFeature,
		},
		"monitorCloudMetrics": {
			Type: schema.TypeSet,
			Elem: accountFeature,
		},
	},
}
var accountFeature = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"components": {
			Type: schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	},
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
	data.Set("cloud_provider_type", cloudauthAccount.Provider.String())

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

func cloudauthAccountFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountSecure {
	components := []*cloudauth.AccountComponent{}
	for _, rc := range data.Get("components").([]interface{}) {
		resourceComponent := rc.(map[string]interface{})
		component := &cloudauth.AccountComponent{}
		for key, value := range resourceComponent {
			switch key {
			case "type":
				component.Type = cloudauth.Component(cloudauth.Component_value[value.(string)])
			case "instance":
				component.Instance = value.(string)
			case "cloudConnectorMetadata":
				cloudConnectorMetadata := &cloudauth.CloudConnectorMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), cloudConnectorMetadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{
						CloudConnectorMetadata: cloudConnectorMetadata,
					}
				}
			case "trustedRoleMetadata":
				metadata := &cloudauth.TrustedRoleMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_TrustedRoleMetadata{
						TrustedRoleMetadata: metadata,
					}
				}
			case "eventBridgeMetadata":
				metadata := &cloudauth.EventBridgeMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{
						EventBridgeMetadata: metadata,
					}
				}
			case "servicePrincipalMetadata":
				metadata := &cloudauth.CloudConnectorMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{
						CloudConnectorMetadata: metadata,
					}
				}
			case "webhookDatasourceMetadata":
				metadata := &cloudauth.WebhookDatasourceMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_WebhookDatasourceMetadata{
						WebhookDatasourceMetadata: metadata,
					}
				}
			case "cryptoKeyMetadata":
				metadata := &cloudauth.CryptoKeyMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_CryptoKeyMetadata{
						CryptoKeyMetadata: metadata,
					}
				}
			case "cloudLogsMetadata":
				metadata := &cloudauth.CloudLogsMetadata{}
				if err := protojson.Unmarshal([]byte(value.(string)), metadata); err == nil {
					component.Metadata = &cloudauth.AccountComponent_CloudLogsMetadata{
						CloudLogsMetadata: metadata,
					}
				}
			}
		}
		components = append(components, component)
	}

	accountFeatures := &cloudauth.AccountFeatures{}
	for name, value := range data.Get("feature").(map[string]interface{}) {
		switch name {
		case "secureConfigPosture":
			accountFeatures.SecureConfigPosture = &cloudauth.AccountFeature{}
			for name2, value2 := range value.(map[string]interface{}) {
				switch name2 {
				case "type":
					accountFeatures.SecureConfigPosture.Type = cloudauth.Feature(cloudauth.Feature_value[value2.(string)])
				case "enabled":
					accountFeatures.SecureConfigPosture.Enabled = value2.(bool)
				case "components":
					for _, componentID := range value2.([]interface{}) {
						accountFeatures.SecureConfigPosture.Components = append(accountFeatures.SecureConfigPosture.Components, componentID.(string))
					}
				}
			}
		case "secureIdentityEntitlement":
			accountFeatures.SecureIdentityEntitlement = &cloudauth.AccountFeature{}
			for name2, value2 := range value.(map[string]interface{}) {
				switch name2 {
				case "type":
					accountFeatures.SecureIdentityEntitlement.Type = cloudauth.Feature(cloudauth.Feature_value[value2.(string)])
				case "enabled":
					accountFeatures.SecureIdentityEntitlement.Enabled = value2.(bool)
				case "components":
					for _, componentID := range value2.([]interface{}) {
						accountFeatures.SecureIdentityEntitlement.Components = append(accountFeatures.SecureIdentityEntitlement.Components, componentID.(string))
					}
				}
			}
		case "secureThreatDetection":
			accountFeatures.SecureThreatDetection = &cloudauth.AccountFeature{}
			for name2, value2 := range value.(map[string]interface{}) {
				switch name2 {
				case "type":
					accountFeatures.SecureThreatDetection.Type = cloudauth.Feature(cloudauth.Feature_value[value2.(string)])
				case "enabled":
					accountFeatures.SecureThreatDetection.Enabled = value2.(bool)
				case "components":
					for _, componentID := range value2.([]interface{}) {
						accountFeatures.SecureThreatDetection.Components = append(accountFeatures.SecureThreatDetection.Components, componentID.(string))
					}
				}
			}
		case "secureAgentlessScanning":
			accountFeatures.SecureAgentlessScanning = &cloudauth.AccountFeature{}
			for name2, value2 := range value.(map[string]interface{}) {
				switch name2 {
				case "type":
					accountFeatures.SecureAgentlessScanning.Type = cloudauth.Feature(cloudauth.Feature_value[value2.(string)])
				case "enabled":
					accountFeatures.SecureAgentlessScanning.Enabled = value2.(bool)
				case "components":
					for _, componentID := range value2.([]interface{}) {
						accountFeatures.SecureAgentlessScanning.Components = append(accountFeatures.SecureAgentlessScanning.Components, componentID.(string))
					}
				}
			}
		case "monitorCloudMetrics":
			accountFeatures.MonitorCloudMetrics = &cloudauth.AccountFeature{}
			for name2, value2 := range value.(map[string]interface{}) {
				switch name2 {
				case "type":
					accountFeatures.MonitorCloudMetrics.Type = cloudauth.Feature(cloudauth.Feature_value[value2.(string)])
				case "enabled":
					accountFeatures.MonitorCloudMetrics.Enabled = value2.(bool)
				case "components":
					for _, componentID := range value2.([]interface{}) {
						accountFeatures.MonitorCloudMetrics.Components = append(accountFeatures.MonitorCloudMetrics.Components, componentID.(string))
					}
				}
			}
		}
	}

	return &v2.CloudauthAccountSecure{
		CloudAccount: cloudauth.CloudAccount{
			Enabled:    data.Get("enabled").(bool),
			ProviderId: data.Get("cloud_provider_id").(string),
			Provider:   cloudauth.Provider(cloudauth.Provider_value[data.Get("cloud_provider_type").(string)]),
			Components: components,
			Feature:    accountFeatures,
		},
	}
}

func cloudauthAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudauthAccountSecure) error {
	for _, err := range []error{
		data.Set("id", cloudAccount.Id),
		data.Set("enabled", cloudAccount.Enabled),
		data.Set("cloud_provider_id", cloudAccount.ProviderId),
		data.Set("cloud_provider_type", cloudAccount.Provider.String()),
		data.Set("components", cloudAccount.Components),
		data.Set("feature", cloudAccount.Feature),
	} {
		if err != nil {
			return err
		}
	}
	return nil
}

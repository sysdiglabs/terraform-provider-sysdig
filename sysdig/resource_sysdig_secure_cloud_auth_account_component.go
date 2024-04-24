package sysdig

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/protobuf/encoding/protojson"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	cloudauth "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2/cloudauth/go"
)

func resourceSysdigSecureCloudauthAccountComponent() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureCloudauthAccountComponentCreate,
		UpdateContext: resourceSysdigSecureCloudauthAccountComponentUpdate,
		ReadContext:   resourceSysdigSecureCloudauthAccountComponentRead,
		DeleteContext: resourceSysdigSecureCloudauthAccountComponentDelete,
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
			SchemaAccountId: {
				Type:     schema.TypeString,
				Required: true,
			},
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
}

func getSecureCloudauthAccountComponentClient(client SysdigClients) (v2.CloudauthAccountComponentSecureInterface, error) {
	return client.sysdigSecureClientV2()
}

func resourceSysdigSecureCloudauthAccountComponentCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountComponentClient((meta.(SysdigClients)))
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := data.Get(SchemaAccountId).(string)
	cloudauthAccountComponent, errStatus, err := client.CreateCloudauthAccountComponentSecure(ctx, accountId, cloudauthAccountComponentFromResourceData(data))
	if err != nil {
		return diag.Errorf("Error creating resource: %s %s", errStatus, err)
	}

	// using tuple 'accountId/componentType/componentInstance' as TF resource identifier
	data.SetId(accountId + "/" + cloudauthAccountComponent.GetType().String() + "/" + cloudauthAccountComponent.GetInstance())
	err = data.Set(SchemaAccountId, accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountComponentRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountComponentClient((meta.(SysdigClients)))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudauthAccountComponent, errStatus, err := client.GetCloudauthAccountComponentSecure(
		ctx, data.Get(SchemaAccountId).(string), data.Get(SchemaType).(string), data.Get(SchemaInstance).(string))

	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	err = cloudauthAccountComponentToResourceData(data, cloudauthAccountComponent)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountComponentUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountComponentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	accountId := data.Get(SchemaAccountId).(string)
	existingCloudAccountComponent, errStatus, err := client.GetCloudauthAccountComponentSecure(
		ctx, accountId, data.Get(SchemaType).(string), data.Get(SchemaInstance).(string))
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error reading resource: %s %s", errStatus, err)
	}

	newCloudAccountComponent := cloudauthAccountComponentFromResourceData(data)

	// validate and reject non-updatable resource schema fields upfront
	err = validateCloudauthAccountComponentUpdate(existingCloudAccountComponent, newCloudAccountComponent)
	if err != nil {
		return diag.Errorf("Error updating resource: %s", err)
	}

	_, errStatus, err = client.UpdateCloudauthAccountComponentSecure(
		ctx, accountId, data.Get(SchemaType).(string), data.Get(SchemaInstance).(string), newCloudAccountComponent)
	if err != nil {
		if strings.Contains(errStatus, "404") {
			return nil
		}
		return diag.Errorf("Error updating resource: %s %s", errStatus, err)
	}

	return nil
}

func resourceSysdigSecureCloudauthAccountComponentDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCloudauthAccountComponentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	errStatus, err := client.DeleteCloudauthAccountComponentSecure(
		ctx, data.Get(SchemaAccountId).(string), data.Get(SchemaType).(string), data.Get(SchemaInstance).(string))
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
func validateCloudauthAccountComponentUpdate(existingComponent *v2.CloudauthAccountComponentSecure, newComponent *v2.CloudauthAccountComponentSecure) error {
	if existingComponent.Type != newComponent.Type || existingComponent.Instance != newComponent.Instance {
		errorInvalidResourceUpdate := fmt.Sprintf("Bad Request. Updating restricted fields not allowed: %s", []string{"type", "instance"})
		return errors.New(errorInvalidResourceUpdate)
	}

	return nil
}

func cloudauthAccountComponentFromResourceData(data *schema.ResourceData) *v2.CloudauthAccountComponentSecure {
	cloudAccountComponent := &v2.CloudauthAccountComponentSecure{
		AccountComponent: cloudauth.AccountComponent{
			Type:     cloudauth.Component(cloudauth.Component_value[data.Get(SchemaType).(string)]),
			Instance: data.Get(SchemaInstance).(string),
		},
	}
	// XXX: naive but simple approach to read resource data, and check for the metadata schema type passed (only one of the types will be passed)
	// then, populate the respective appropriate metadata proto in accountComponent object for cloudauth.
	var err error
	if resourceMetadata := data.Get(SchemaCloudConnectorMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_CloudConnectorMetadata{CloudConnectorMetadata: &cloudauth.CloudConnectorMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetCloudConnectorMetadata())

	} else if resourceMetadata = data.Get(SchemaTrustedRoleMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_TrustedRoleMetadata{TrustedRoleMetadata: &cloudauth.TrustedRoleMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetTrustedRoleMetadata())

	} else if resourceMetadata = data.Get(SchemaEventBridgeMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_EventBridgeMetadata{EventBridgeMetadata: &cloudauth.EventBridgeMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetEventBridgeMetadata())

	} else if resourceMetadata = data.Get(SchemaServicePrincipalMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_ServicePrincipalMetadata{ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetServicePrincipalMetadata())
		// special handling for GCP service principal if it has keys which are base64 encoded
		if cloudAccountComponent.GetServicePrincipalMetadata().GetGcp() != nil {
			cloudAccountComponent.Metadata = constructGcpServicePrincipalMetadata(resourceMetadata)
		}

	} else if resourceMetadata = data.Get(SchemaWebhookDatasourceMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_WebhookDatasourceMetadata{WebhookDatasourceMetadata: &cloudauth.WebhookDatasourceMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetWebhookDatasourceMetadata())

	} else if resourceMetadata = data.Get(SchemaCryptoKeyMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_CryptoKeyMetadata{CryptoKeyMetadata: &cloudauth.CryptoKeyMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetCryptoKeyMetadata())

	} else if resourceMetadata = data.Get(SchemaCloudLogsMetadata).(string); resourceMetadata != "" {
		cloudAccountComponent.Metadata = &cloudauth.AccountComponent_CloudLogsMetadata{CloudLogsMetadata: &cloudauth.CloudLogsMetadata{}}
		err = protojson.Unmarshal([]byte(resourceMetadata), cloudAccountComponent.GetCloudLogsMetadata())

	} else {
		diag.FromErr(errors.New("bad request, invalid account component metadata type"))
	}

	if err != nil {
		diag.FromErr(err)
	}

	return cloudAccountComponent
}

func cloudauthAccountComponentToResourceData(data *schema.ResourceData, cloudAccountComponent *v2.CloudauthAccountComponentSecure) error {

	accountId := data.Get(SchemaAccountId).(string)
	data.SetId(accountId + "/" + cloudAccountComponent.GetType().String() + "/" + cloudAccountComponent.GetInstance())

	err := data.Set(SchemaAccountId, accountId)
	if err != nil {
		return err
	}

	err = data.Set(SchemaType, cloudAccountComponent.GetType().String())
	if err != nil {
		return err
	}

	err = data.Set(SchemaInstance, cloudAccountComponent.GetInstance())
	if err != nil {
		return err
	}

	// XXX: naive but simple approach to read accountComponent object from cloudauth, and check for the metadata proto type (only one of the types will be returned)
	// then, populate the respective appropriate metadata schema in resource data.
	switch cloudAccountComponent.GetType() {
	case cloudauth.Component_COMPONENT_CLOUD_CONNECTOR:
		err = data.Set(SchemaCloudConnectorMetadata, getComponentMetadataString(cloudAccountComponent.GetCloudConnectorMetadata()))
	case cloudauth.Component_COMPONENT_TRUSTED_ROLE:
		err = data.Set(SchemaTrustedRoleMetadata, getComponentMetadataString(cloudAccountComponent.GetTrustedRoleMetadata()))
	case cloudauth.Component_COMPONENT_EVENT_BRIDGE:
		err = data.Set(SchemaEventBridgeMetadata, getComponentMetadataString(cloudAccountComponent.GetEventBridgeMetadata()))
	case cloudauth.Component_COMPONENT_SERVICE_PRINCIPAL:
		// special handling for GCP service principal if it has keys which are to be base64 encoded
		if cloudAccountComponent.GetServicePrincipalMetadata().GetGcp() != nil {
			err = data.Set(SchemaServicePrincipalMetadata, getGcpServicePrincipalMetadata(cloudAccountComponent.GetServicePrincipalMetadata()))
		} else {
			err = data.Set(SchemaServicePrincipalMetadata, getComponentMetadataString(cloudAccountComponent.GetServicePrincipalMetadata()))
		}
	case cloudauth.Component_COMPONENT_WEBHOOK_DATASOURCE:
		err = data.Set(SchemaWebhookDatasourceMetadata, getComponentMetadataString(cloudAccountComponent.GetWebhookDatasourceMetadata()))
	case cloudauth.Component_COMPONENT_CRYPTO_KEY:
		err = data.Set(SchemaCryptoKeyMetadata, getComponentMetadataString(cloudAccountComponent.GetCryptoKeyMetadata()))
	case cloudauth.Component_COMPONENT_CLOUD_LOGS:
		err = data.Set(SchemaCloudLogsMetadata, getComponentMetadataString(cloudAccountComponent.GetCloudLogsMetadata()))
	}

	if err != nil {
		return err
	}

	return nil
}

// internal type redefintion for GCP service principals.
// This exists because in terraform, the key is originally provided in the form of a base64 encoded json string

// note; caution with order of fields, they have to go in alphabetical ASC so that the json marshalled on the tf read phase produces no drift https://github.com/golang/go/issues/27179
type internalServicePrincipalMetadata_GCP struct {
	Email                      string                                                             `json:"email,omitempty"`
	Key                        string                                                             `json:"key,omitempty"` // base64 encoded
	WorkloadIdentityFederation *cloudauth.ServicePrincipalMetadata_GCP_WorkloadIdentityFederation `json:"workload_identity_federation,omitempty"`
}
type internalServicePrincipalMetadata struct {
	Gcp *internalServicePrincipalMetadata_GCP `json:"gcp,omitempty"`
}

/*
This helper function is for special handling of GCP service principal metadata if it is key-based.
It takes the metadata coming from resource data and returns populated metadata proto, with decoded key if present.
*/
func constructGcpServicePrincipalMetadata(metadata string) *cloudauth.AccountComponent_ServicePrincipalMetadata {
	spGcp := &internalServicePrincipalMetadata{}
	err := json.Unmarshal([]byte(metadata), spGcp)
	if err != nil {
		diag.FromErr(err)
	}
	// special handling if GCP service principal key is present, decode and unmarshal it before populating all the metadata
	var spGcpKey *cloudauth.ServicePrincipalMetadata_GCP_Key
	if len(spGcp.Gcp.Key) > 0 {
		var spGcpKeyBytes []byte
		spGcpKeyBytes, err = base64.StdEncoding.DecodeString(spGcp.Gcp.Key)
		if err != nil {
			diag.FromErr(err)
		}
		err = json.Unmarshal(spGcpKeyBytes, &spGcpKey)
		if err != nil {
			diag.FromErr(err)
		}
	}
	return &cloudauth.AccountComponent_ServicePrincipalMetadata{
		ServicePrincipalMetadata: &cloudauth.ServicePrincipalMetadata{
			Provider: &cloudauth.ServicePrincipalMetadata_Gcp{
				Gcp: &cloudauth.ServicePrincipalMetadata_GCP{
					Key:                        spGcpKey,
					WorkloadIdentityFederation: spGcp.Gcp.WorkloadIdentityFederation,
					Email:                      spGcp.Gcp.Email,
				},
			},
		},
	}
}

/*
This helper function is for special handling of GCP service principal metadata if it is key-based.
It takes the metadata coming from cloudauth metadata proto and returns component metadata string, with encoded key if present.
*/
func getGcpServicePrincipalMetadata(metadata *cloudauth.ServicePrincipalMetadata) string {
	var gcpKeyBytes []byte
	if metadata.GetGcp().GetKey() != nil {
		var err error
		gcpKeyBytes, err = protojson.MarshalOptions{UseProtoNames: true}.Marshal(metadata.GetGcp().GetKey())
		if err != nil {
			diag.FromErr(err)
		}
		var gcpKeyBytesBuffer bytes.Buffer
		err = json.Indent(&gcpKeyBytesBuffer, gcpKeyBytes, "", "  ")
		if err != nil {
			diag.FromErr(err)
		}
		gcpKeyBytes = append(gcpKeyBytesBuffer.Bytes(), '\n')
	}
	spGcpBytes, err := json.Marshal(&internalServicePrincipalMetadata{
		Gcp: &internalServicePrincipalMetadata_GCP{
			Key:                        base64.StdEncoding.EncodeToString(gcpKeyBytes),
			WorkloadIdentityFederation: metadata.GetGcp().GetWorkloadIdentityFederation(),
			Email:                      metadata.GetGcp().GetEmail(),
		},
	})
	if err != nil {
		diag.FromErr(err)
	}

	return string(spGcpBytes)
}

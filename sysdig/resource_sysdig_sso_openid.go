package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSSOOpenID() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigSSOOpenIDRead,
		CreateContext: resourceSysdigSSOOpenIDCreate,
		UpdateContext: resourceSysdigSSOOpenIDUpdate,
		DeleteContext: resourceSysdigSSOOpenIDDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		CustomizeDiff: validateSSOOpenIDMetadata,
		Schema: map[string]*schema.Schema{
			// Required fields
			"issuer_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OpenID Connect issuer URL (e.g., https://accounts.google.com)",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The OAuth 2.0 client ID",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The OAuth 2.0 client secret",
			},

			// Optional base SSO fields
			"product": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "secure",
				ValidateFunc: validation.StringInSlice([]string{"monitor", "secure"}, false),
				Description:  "The Sysdig product (monitor or secure)",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the SSO configuration is active",
			},
			"create_user_on_login": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to create a new user upon first login",
			},
			"is_single_logout_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether single logout is enabled",
			},
			"is_group_mapping_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether group mapping is enabled",
			},
			"group_mapping_attribute_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "groups",
				Description: "The attribute name for group mapping",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A custom name for this SSO integration",
			},

			// OpenID specific optional fields
			"is_metadata_discovery_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to use automatic metadata discovery from the issuer URL",
			},
			"group_attribute_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "groups",
				Description: "The attribute name for groups in the ID token",
			},
			"is_additional_scopes_check_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether additional scopes check is enabled",
			},
			"additional_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Additional OAuth scopes to request",
			},

			// Metadata block (required if is_metadata_discovery_enabled = false)
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The issuer identifier",
						},
						"authorization_endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The authorization endpoint URL",
						},
						"token_endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The token endpoint URL",
						},
						"jwks_uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The JWKS URI for token verification",
						},
						"token_auth_method": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CLIENT_SECRET_BASIC", "CLIENT_SECRET_POST"}, false),
							Description:  "The token authentication method (CLIENT_SECRET_BASIC or CLIENT_SECRET_POST)",
						},
						"end_session_endpoint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The end session endpoint URL for logout",
						},
						"user_info_endpoint": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user info endpoint URL",
						},
					},
				},
				Description: "Manual metadata configuration (required when is_metadata_discovery_enabled is false)",
			},

			// Computed field
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the SSO configuration (used for optimistic locking)",
			},
		},
	}
}

func validateSSOOpenIDMetadata(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	isMetadataDiscoveryEnabled := diff.Get("is_metadata_discovery_enabled").(bool)
	metadata := diff.Get("metadata").([]any)

	if !isMetadataDiscoveryEnabled && len(metadata) == 0 {
		return diff.ForceNew("metadata")
	}

	return nil
}

func resourceSysdigSSOOpenIDRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sso, err := client.GetSSOOpenID(ctx, id)
	if err != nil {
		if err == v2.ErrSSOOpenIDNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return ssoOpenIDToResourceData(sso, d)
}

func resourceSysdigSSOOpenIDCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	sso := ssoOpenIDFromResourceData(d)
	sso.Type = "OPENID"

	created, err := client.CreateSSOOpenID(ctx, sso)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(created.ID))

	return resourceSysdigSSOOpenIDRead(ctx, d, m)
}

func resourceSysdigSSOOpenIDUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sso := ssoOpenIDFromResourceData(d)
	sso.ID = id
	sso.Type = "OPENID"
	sso.Version = d.Get("version").(int)

	_, err = client.UpdateSSOOpenID(ctx, id, sso)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSysdigSSOOpenIDRead(ctx, d, m)
}

func resourceSysdigSSOOpenIDDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteSSOOpenID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ssoOpenIDFromResourceData(d *schema.ResourceData) *v2.SSOOpenID {
	sso := &v2.SSOOpenID{
		IssuerURL:                      d.Get("issuer_url").(string),
		ClientID:                       d.Get("client_id").(string),
		ClientSecret:                   d.Get("client_secret").(string),
		Product:                        d.Get("product").(string),
		IsActive:                       d.Get("is_active").(bool),
		CreateUserOnLogin:              d.Get("create_user_on_login").(bool),
		IsSingleLogoutEnabled:          d.Get("is_single_logout_enabled").(bool),
		IsGroupMappingEnabled:          d.Get("is_group_mapping_enabled").(bool),
		GroupMappingAttributeName:      d.Get("group_mapping_attribute_name").(string),
		IntegrationName:                d.Get("integration_name").(string),
		IsMetadataDiscoveryEnabled:     d.Get("is_metadata_discovery_enabled").(bool),
		GroupAttributeName:             d.Get("group_attribute_name").(string),
		IsAdditionalScopesCheckEnabled: d.Get("is_additional_scopes_check_enabled").(bool),
	}

	// Handle additional scopes
	if v, ok := d.GetOk("additional_scopes"); ok {
		scopesInterface := v.([]any)
		scopes := make([]string, len(scopesInterface))
		for i, s := range scopesInterface {
			scopes[i] = s.(string)
		}
		sso.AdditionalScopes = scopes
	}

	// Handle metadata block
	if v, ok := d.GetOk("metadata"); ok {
		metadataList := v.([]any)
		if len(metadataList) > 0 {
			metadata := metadataList[0].(map[string]any)
			sso.Metadata = &v2.OpenIDMetadata{
				Issuer:                metadata["issuer"].(string),
				AuthorizationEndpoint: metadata["authorization_endpoint"].(string),
				TokenEndpoint:         metadata["token_endpoint"].(string),
				JwksURI:               metadata["jwks_uri"].(string),
				TokenAuthMethod:       metadata["token_auth_method"].(string),
				EndSessionEndpoint:    metadata["end_session_endpoint"].(string),
				UserInfoEndpoint:      metadata["user_info_endpoint"].(string),
			}
		}
	}

	return sso
}

func ssoOpenIDToResourceData(sso *v2.SSOOpenID, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("issuer_url", sso.IssuerURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_id", sso.ClientID); err != nil {
		return diag.FromErr(err)
	}
	// Note: client_secret is not returned by the API, so we don't set it here
	// to avoid diff issues with the sensitive value

	if err := d.Set("product", sso.Product); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_active", sso.IsActive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("create_user_on_login", sso.CreateUserOnLogin); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_single_logout_enabled", sso.IsSingleLogoutEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_group_mapping_enabled", sso.IsGroupMappingEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_mapping_attribute_name", sso.GroupMappingAttributeName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("integration_name", sso.IntegrationName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_metadata_discovery_enabled", sso.IsMetadataDiscoveryEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("group_attribute_name", sso.GroupAttributeName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_additional_scopes_check_enabled", sso.IsAdditionalScopesCheckEnabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("additional_scopes", sso.AdditionalScopes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("version", sso.Version); err != nil {
		return diag.FromErr(err)
	}

	// Handle metadata block
	if sso.Metadata != nil {
		metadata := []map[string]any{
			{
				"issuer":                 sso.Metadata.Issuer,
				"authorization_endpoint": sso.Metadata.AuthorizationEndpoint,
				"token_endpoint":         sso.Metadata.TokenEndpoint,
				"jwks_uri":               sso.Metadata.JwksURI,
				"token_auth_method":      sso.Metadata.TokenAuthMethod,
				"end_session_endpoint":   sso.Metadata.EndSessionEndpoint,
				"user_info_endpoint":     sso.Metadata.UserInfoEndpoint,
			},
		}
		if err := d.Set("metadata", metadata); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

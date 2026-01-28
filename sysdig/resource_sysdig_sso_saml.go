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

func resourceSysdigSSOSaml() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigSSOSamlRead,
		CreateContext: resourceSysdigSSOSamlCreate,
		UpdateContext: resourceSysdigSSOSamlUpdate,
		DeleteContext: resourceSysdigSSOSamlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		CustomizeDiff: validateSSOSamlMetadata,
		Schema: map[string]*schema.Schema{
			// SAML metadata - mutually exclusive
			"metadata_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The URL to fetch SAML metadata from the IdP",
				ExactlyOneOf: []string{"metadata_url", "metadata_xml"},
			},
			"metadata_xml": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The raw SAML metadata XML from the IdP",
				ExactlyOneOf: []string{"metadata_url", "metadata_xml"},
			},

			// Required field
			"email_parameter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The SAML attribute name that contains the user's email address",
			},

			// Optional base SSO fields (shared with OpenID)
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
				Description: "The SAML attribute name for group mapping",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A custom name for this SSO integration",
			},

			// SAML specific optional fields (security)
			"is_signature_validation_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether SAML response signature validation is enabled",
			},
			"is_signed_assertion_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether signed SAML assertions are required",
			},
			"is_destination_verification_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether destination verification is enabled",
			},
			"is_encryption_support_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether encryption support is enabled",
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

func validateSSOSamlMetadata(_ context.Context, diff *schema.ResourceDiff, _ any) error {
	// ExactlyOneOf already handles mutual exclusion, no additional validation needed
	return nil
}

func resourceSysdigSSOSamlRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sso, err := client.GetSSOSaml(ctx, id)
	if err != nil {
		if err == v2.ErrSSOSamlNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return ssoSamlToResourceData(sso, d)
}

func resourceSysdigSSOSamlCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	sso := ssoSamlFromResourceData(d)
	sso.Type = "SAML"

	created, err := client.CreateSSOSaml(ctx, sso)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(created.ID))

	return resourceSysdigSSOSamlRead(ctx, d, m)
}

func resourceSysdigSSOSamlUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	sso := ssoSamlFromResourceData(d)
	sso.ID = id
	sso.Type = "SAML"
	sso.Version = d.Get("version").(int)

	_, err = client.UpdateSSOSaml(ctx, id, sso)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSysdigSSOSamlRead(ctx, d, m)
}

func resourceSysdigSSOSamlDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteSSOSaml(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ssoSamlFromResourceData(d *schema.ResourceData) *v2.SSOSaml {
	sso := &v2.SSOSaml{
		MetadataURL:               d.Get("metadata_url").(string),
		MetadataXML:               d.Get("metadata_xml").(string),
		EmailParameter:            d.Get("email_parameter").(string),
		Product:                   d.Get("product").(string),
		IsActive:                  d.Get("is_active").(bool),
		CreateUserOnLogin:         d.Get("create_user_on_login").(bool),
		IsSingleLogoutEnabled:     d.Get("is_single_logout_enabled").(bool),
		IsGroupMappingEnabled:     d.Get("is_group_mapping_enabled").(bool),
		GroupMappingAttributeName: d.Get("group_mapping_attribute_name").(string),
		IntegrationName:           d.Get("integration_name").(string),
	}

	// Handle SAML security fields (using pointers to distinguish unset from false)
	isSignatureValidationEnabled := d.Get("is_signature_validation_enabled").(bool)
	sso.IsSignatureValidationEnabled = &isSignatureValidationEnabled

	isSignedAssertionEnabled := d.Get("is_signed_assertion_enabled").(bool)
	sso.IsSignedAssertionEnabled = &isSignedAssertionEnabled

	isDestinationVerificationEnabled := d.Get("is_destination_verification_enabled").(bool)
	sso.IsDestinationVerificationEnabled = &isDestinationVerificationEnabled

	isEncryptionSupportEnabled := d.Get("is_encryption_support_enabled").(bool)
	sso.IsEncryptionSupportEnabled = &isEncryptionSupportEnabled

	return sso
}

func ssoSamlToResourceData(sso *v2.SSOSaml, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := d.Set("metadata_url", sso.MetadataURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("metadata_xml", sso.MetadataXML); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email_parameter", sso.EmailParameter); err != nil {
		return diag.FromErr(err)
	}
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
	if err := d.Set("version", sso.Version); err != nil {
		return diag.FromErr(err)
	}

	// Handle SAML security fields
	if sso.IsSignatureValidationEnabled != nil {
		if err := d.Set("is_signature_validation_enabled", *sso.IsSignatureValidationEnabled); err != nil {
			return diag.FromErr(err)
		}
	}
	if sso.IsSignedAssertionEnabled != nil {
		if err := d.Set("is_signed_assertion_enabled", *sso.IsSignedAssertionEnabled); err != nil {
			return diag.FromErr(err)
		}
	}
	if sso.IsDestinationVerificationEnabled != nil {
		if err := d.Set("is_destination_verification_enabled", *sso.IsDestinationVerificationEnabled); err != nil {
			return diag.FromErr(err)
		}
	}
	if sso.IsEncryptionSupportEnabled != nil {
		if err := d.Set("is_encryption_support_enabled", *sso.IsEncryptionSupportEnabled); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

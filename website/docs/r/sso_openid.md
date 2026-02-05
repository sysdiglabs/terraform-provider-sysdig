---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_sso_openid"
description: |-
  Creates an SSO OpenID Connect configuration in Sysdig.
---

# Resource: sysdig_sso_openid

Creates an SSO OpenID Connect configuration in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Basic Configuration with Metadata Discovery

```terraform
resource "sysdig_sso_openid" "google" {
  issuer_url       = "https://accounts.google.com"
  client_id        = "your-client-id.apps.googleusercontent.com"
  client_secret    = "your-client-secret"
  integration_name = "Google SSO"

  is_active                   = true
  create_user_on_login        = true
  is_metadata_discovery_enabled = true
}
```

### Configuration with Manual Metadata

When using an identity provider that doesn't support metadata discovery, you can provide the metadata manually:

```terraform
resource "sysdig_sso_openid" "custom_idp" {
  issuer_url       = "https://idp.example.com"
  client_id        = "your-client-id"
  client_secret    = "your-client-secret"
  integration_name = "Custom IDP"

  is_active                     = true
  is_metadata_discovery_enabled = false

  metadata {
    issuer                 = "https://idp.example.com"
    authorization_endpoint = "https://idp.example.com/oauth2/authorize"
    token_endpoint         = "https://idp.example.com/oauth2/token"
    jwks_uri               = "https://idp.example.com/.well-known/jwks.json"
    token_auth_method      = "CLIENT_SECRET_BASIC"
    end_session_endpoint   = "https://idp.example.com/oauth2/logout"
    user_info_endpoint     = "https://idp.example.com/userinfo"
  }
}
```

### Configuration with Group Mapping and Additional Scopes

```terraform
resource "sysdig_sso_openid" "okta" {
  issuer_url       = "https://your-org.okta.com"
  client_id        = "your-client-id"
  client_secret    = "your-client-secret"
  integration_name = "Okta SSO"

  is_active                          = true
  create_user_on_login               = true
  is_group_mapping_enabled           = true
  group_mapping_attribute_name       = "groups"
  is_single_logout_enabled           = true

  is_additional_scopes_check_enabled = true
  additional_scopes                  = ["groups", "profile", "email"]
}
```

## Argument Reference

### Required Arguments

* `issuer_url` - (Required) The OpenID Connect issuer URL (e.g., `https://accounts.google.com`).

* `client_id` - (Required) The OAuth 2.0 client ID.

* `client_secret` - (Required, Sensitive) The OAuth 2.0 client secret.

### Optional Arguments

* `product` - (Optional) The Sysdig product to configure SSO for. Valid values are `monitor` or `secure`. Default is `secure`.

* `is_active` - (Optional) Whether the SSO configuration is active. Default is `true`.

* `create_user_on_login` - (Optional) Whether to create a new user upon first login. Default is `false`.

* `is_single_logout_enabled` - (Optional) Whether single logout (SLO) is enabled. Default is `false`.

* `is_group_mapping_enabled` - (Optional) Whether group mapping is enabled. Default is `false`.

* `group_mapping_attribute_name` - (Optional) The attribute name for group mapping in the ID token claims. Default is `groups`.

* `integration_name` - (Optional) A name to distinguish different SSO integrations. Users can select this integration on the login page.

* `is_metadata_discovery_enabled` - (Optional) Whether to use automatic metadata discovery from the issuer URL. Default is `true`.

* `is_additional_scopes_check_enabled` - (Optional) Whether additional scopes check is enabled. Default is `false`.

* `additional_scopes` - (Optional) A list of additional OAuth scopes to request.

* `metadata` - (Optional) Manual metadata configuration. Required when `is_metadata_discovery_enabled` is `false`. See [Metadata](#metadata) below for details.

### Metadata

The `metadata` block supports the following arguments:

* `issuer` - (Required) The issuer identifier.

* `authorization_endpoint` - (Required) The authorization endpoint URL.

* `token_endpoint` - (Required) The token endpoint URL.

* `jwks_uri` - (Required) The JWKS URI for token verification.

* `token_auth_method` - (Required) The token authentication method. Valid values are `CLIENT_SECRET_BASIC` or `CLIENT_SECRET_POST`.

* `end_session_endpoint` - (Optional) The end session endpoint URL for logout.

* `user_info_endpoint` - (Optional) The user info endpoint URL.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The version of the SSO configuration (used for optimistic locking).

## Import

Sysdig SSO OpenID configurations can be imported using the ID, e.g.

```
$ terraform import sysdig_sso_openid.example 12345
```

~> **Note:** The `client_secret` attribute cannot be imported and must be set in the configuration after import.

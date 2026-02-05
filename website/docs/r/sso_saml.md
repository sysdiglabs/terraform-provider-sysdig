---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_sso_saml"
description: |-
  Creates a SAML SSO configuration in Sysdig.
---

# Resource: sysdig_sso_saml

Creates a SAML Single Sign-On (SSO) configuration in Sysdig.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Basic example with metadata URL

```terraform
resource "sysdig_sso_saml" "example" {
  metadata_url     = "https://idp.example.com/app/sysdig/sso/saml/metadata"
  email_parameter  = "email"
  integration_name = "Corporate SAML SSO"
  is_active        = true
}
```

### Example with inline metadata XML

```terraform
resource "sysdig_sso_saml" "example_xml" {
  metadata_xml = <<-EOF
<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://idp.example.com">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://idp.example.com/sso"/>
  </IDPSSODescriptor>
</EntityDescriptor>
EOF

  email_parameter  = "email"
  integration_name = "Corporate SAML SSO"
  is_active        = true
}
```

### Example with group mapping enabled

```terraform
resource "sysdig_sso_saml" "example_groups" {
  metadata_url                  = "https://idp.example.com/app/sysdig/sso/saml/metadata"
  email_parameter               = "email"
  integration_name              = "Corporate SAML SSO"
  is_active                     = true
  is_group_mapping_enabled      = true
  group_mapping_attribute_name  = "groups"
}
```

### Example with custom security settings

```terraform
resource "sysdig_sso_saml" "example_security" {
  metadata_url                        = "https://idp.example.com/app/sysdig/sso/saml/metadata"
  email_parameter                     = "email"
  integration_name                    = "Corporate SAML SSO"
  is_active                           = true
  is_signature_validation_enabled     = true
  is_signed_assertion_enabled         = true
  is_destination_verification_enabled = true
  is_encryption_support_enabled       = false
}
```

## Argument Reference

### Required

* `email_parameter` - (Required) The SAML attribute name that contains the user's email address.

### Metadata (exactly one required)

* `metadata_url` - (Optional) The URL to fetch SAML metadata from the Identity Provider. Mutually exclusive with `metadata_xml`.
* `metadata_xml` - (Optional) The raw SAML metadata XML from the Identity Provider. Mutually exclusive with `metadata_url`.

### Optional

* `product` - (Optional) The Sysdig product to configure SSO for. Valid values are `monitor` or `secure`. Default: `secure`.
* `is_active` - (Optional) Whether the SSO configuration is active. Default: `true`.
* `create_user_on_login` - (Optional) Whether to automatically create a new user upon first login via SSO. Default: `false`.
* `is_single_logout_enabled` - (Optional) Whether SAML Single Logout (SLO) is enabled. Default: `false`.
* `is_group_mapping_enabled` - (Optional) Whether group mapping from SAML attributes is enabled. Default: `false`.
* `group_mapping_attribute_name` - (Optional) The SAML attribute name that contains group membership information. Default: `groups`.
* `integration_name` - (Optional) A name to distinguish different SSO integrations. Users can select this integration on the login page.

### Security Settings (Optional)

* `is_signature_validation_enabled` - (Optional) Whether SAML response signature validation is enabled. Default: `true`.
* `is_signed_assertion_enabled` - (Optional) Whether signed SAML assertions are required. Default: `true`.
* `is_destination_verification_enabled` - (Optional) Whether destination verification in SAML responses is enabled. Default: `true`.
* `is_encryption_support_enabled` - (Optional) Whether SAML encryption support is enabled. Default: `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The version of the SSO configuration, used for optimistic locking during updates.

## Import

SAML SSO configurations can be imported using the SSO configuration ID:

```
$ terraform import sysdig_sso_saml.example 12345
```

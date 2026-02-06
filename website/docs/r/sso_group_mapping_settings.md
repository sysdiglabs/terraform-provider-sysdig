---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_sso_group_mapping_settings"
description: |-
  Manages SSO group mapping settings in Sysdig using the Platform API.
---

# Resource: sysdig_sso_group_mapping_settings

Manages SSO group mapping conflict resolution settings in Sysdig using the Platform API. This resource replaces the deprecated `sysdig_group_mapping_config` resource.

This is a singleton resource — only one instance should exist per Sysdig account. The resource cannot be deleted; removing it from Terraform configuration will only remove it from state.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Basic configuration

```terraform
resource "sysdig_sso_group_mapping_settings" "default" {
  no_mapping_strategy              = "UNAUTHORIZED"
  different_roles_same_team_strategy = "UNAUTHORIZED"
}
```

### With error redirect

```terraform
resource "sysdig_sso_group_mapping_settings" "with_redirect" {
  no_mapping_strategy              = "NO_MAPPINGS_ERROR_REDIRECT"
  different_roles_same_team_strategy = "HIGHEST_ROLE"
  no_mappings_error_redirect_url   = "https://example.com/sso-error"
}
```

## Argument Reference

* `no_mapping_strategy` - (Required) Strategy when no group mapping matches a user. Valid values:
  * `UNAUTHORIZED` - Deny access.
  * `DEFAULT_TEAM_DEFAULT_ROLE` - Assign default team and role.
  * `NO_MAPPINGS_ERROR_REDIRECT` - Redirect to an error URL (requires `no_mappings_error_redirect_url`).

* `different_roles_same_team_strategy` - (Required) Strategy when a user matches multiple mappings with different roles for the same team. Valid values:
  * `UNAUTHORIZED` - Deny access.
  * `HIGHEST_ROLE` - Use the highest-privilege role.
  * `LOWEST_ROLE` - Use the lowest-privilege role.

* `no_mappings_error_redirect_url` - (Optional) URL to redirect users when `no_mapping_strategy` is `NO_MAPPINGS_ERROR_REDIRECT`. Maximum 2048 characters.

## Attributes Reference

No additional attributes are exported.

## Import

SSO group mapping settings can be imported using the static ID `sso_group_mapping_settings`:

```
$ terraform import sysdig_sso_group_mapping_settings.default sso_group_mapping_settings
```

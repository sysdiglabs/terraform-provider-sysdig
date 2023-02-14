---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_scanning_policy"
description: |-
  Creates a Sysdig Secure Scanning Policy for Legacy Scanning Engine.
---

# Resource: sysdig_secure_scanning_policy

Creates a Sysdig Secure Policy.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.  

## Example Usage

```terraform
resource "sysdig_secure_scanning_policy" "scanning_policy_example" {
  name = "Scanning Policy Name"
  comment = "Scanning Policy Description"

  // Scanning Policy Rules (gates) and parameters for configuration
  rules {
    gate = "dockerfile"
    trigger = "effective_user"
    action = "WARN"
    params {
        name = "users"
        value = "docker"
    }
    params {
        name = "type"
        value = "blacklist"
    }
  }

  rules {
    gate = "files"
    trigger = "attribute_match"
    action = "WARN"
    params {
        name = "filename"
        value = "/etc/passwd"
    }
  }

  rules {
    gate = "vulnerabilities"
    trigger = "package"
    action = "WARN"
    params {
      name = "package_type"
      value = "all"
    }
    params {
      name = "severity"
      value = "medium"
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Secure policy. It must be unique.

* `comment` - (Required) The description of Secure scanning policy.

* `rules` - (Optional) Define all rules included in the Policy for scanning detection.

- - -

### Rules block

* `gate` - (Required) Must be one of `always`, `dockerfile`, `files`, `licenses`, `metadata`, `npms`, `packages`, `passwd_file`, `retrieved_files`, `vulnerabilities`, `secret_scans`, `ruby_gems`. You can see the description of each gate in this [link](https://docs.sysdig.com/en/docs/sysdig-secure/scanning/manage-scanning-policies/scanning-policy-gates-and-triggers/).
* `trigger` - (Required) Each gate have different trigger options and parameters. Check possible triggers per gate in the previous link.
* `params` - (Required) Each gate and trigger options have different parameter configurations. Review the previous link to see all options.
* `action` - (Required) define the action to take if one gate triggers what would affect the policy results. Must be `WARN` or `STOP`.

- - -

## Attributes Reference

No additional attributes are exported.

## Import

Secure scanning policies can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_scanning_policy.example policy_123456
```

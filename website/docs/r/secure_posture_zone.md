---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_zone"
description: |-
  Creates a Sysdig Secure Posture Zone.
---

# Resource: sysdig_secure_posture_zone

Creates a Sysdig Secure Posture Zone.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_posture_policies" "all" {}

locals {
  fedramp_policies = [
    for p in data.sysdig_secure_posture_policies.all.policies :
    p if length(regexall(".*FedRAMP.*", p.name)) > 0
  ]
}

resource "sysdig_secure_posture_zone" "example" {
  name        = "Zone with FedRAMP policies"
  description = "Zone description"
  policy_ids  = [for p in local.fedramp_policies : p.id]

  scopes {
    scope {
      target_type = "aws"
      rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
    }

    scope {
      target_type = "azure"
      rules       = "organization contains \"o1\""
    }
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Posture Zone.
* `description` - (Optional) The description of the Posture Zone.
* `policy_ids` - (Optional) The list of Posture Policy IDs attached to Zone.
* `scopes` - (Optional) Scopes block defines list of scopes attached to Zone.

### Scopes block

* `target_type` - (Required) The target type for the scope.
* `rules` - (Optional) Rules attached to scope.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `author` - (Computed) The zone author.
* `last_modified_by` - (Computed) By whom is last modification made.
* `last_updated` - (Computed) Timestamp of last modification of zone.

## Import

Posture zone can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_zone.example 12345
```

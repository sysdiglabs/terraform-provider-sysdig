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

resource "sysdig_secure_posture_zone" "z1" {
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

---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_zone_posture_policy_assignment"
description: |-
  Manages the association between a Sysdig Secure Zone and a set of posture policies.
---

# Resource: sysdig_secure_zone_posture_policy_assignment

Manages the association between a [`sysdig_secure_zone`](secure_zone.html) and a set of posture policy IDs.

Each zone can have at most one assignment. Updating the resource replaces the entire policy list (PUT semantics).

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_zone" "production" {
  name = "Production"
  scope {
    target_type = "aws"
    expression {
      field    = "account"
      operator = "in"
      values   = ["111111111111"]
    }
  }
}

data "sysdig_secure_posture_policy" "cis_k8s" {
  name = "CIS Kubernetes V1.24 Benchmark"
}

resource "sysdig_secure_zone_posture_policy_assignment" "production" {
  zone_id    = sysdig_secure_zone.production.id
  policy_ids = [data.sysdig_secure_posture_policy.cis_k8s.id]
}
```

## Argument Reference

- `zone_id` - (Required, ForceNew) The ID of the zone to associate policies with. Changing this forces a new resource.
- `policy_ids` - (Required) Set of posture policy IDs to associate with the zone. Updates replace the entire list.

## Attributes Reference

No additional attributes are exported beyond the arguments.

## Import

The resource can be imported using the zone ID:

```
$ terraform import sysdig_secure_zone_posture_policy_assignment.example 12345
```

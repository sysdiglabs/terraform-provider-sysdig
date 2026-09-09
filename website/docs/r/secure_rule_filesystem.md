---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_filesystem"
description: |-
  Creates a Sysdig Secure Filesystem Rule.
---

# Resource: sysdig_secure_rule_filesystem

Creates a Sysdig Secure Filesystem Rule.

!> **Removed:** `sysdig_secure_rule_filesystem` is no longer functional. List-matching ("fast engine") rule support was deprecated on 2025-12-15 and removed from the Sysdig backend on 2026-02-28. `POST /api/secure/rules` now rejects `ruleType FILESYSTEM` with `HTTP 400: The field details has an unknown ruleType: FILESYSTEM`. The provider no longer attempts the call: any `create`, `read` or `update` on this resource fails immediately with migration guidance. Rewrite the detection as [`sysdig_secure_rule_falco`](secure_rule_falco.md) with an equivalent Falco condition. This resource will be deleted entirely in the next major version.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform

resource "sysdig_secure_rule_filesystem"  "example" {
  name = "Apache writing to non allowed directory"
  description = "Attempt to write to directories that should be immutable"
  tags = ["filesystem", "cis"]

  read_only {
    matching = true // default
    paths = ["/etc"]
  }

  read_write {
    matching = true // default
    paths = ["/var/log/apache2", "/dev/tty"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

### Read Only

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

### Read Write

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

## Migrating off this resource

The backend no longer stores rules of type `FILESYSTEM`, so there is nothing left to import or
refresh. Replace the `sysdig_secure_rule_filesystem` block with a
[`sysdig_secure_rule_falco`](secure_rule_falco.md) block expressing the same detection as a Falco
condition.

No state surgery is needed. The provider reports these resources as gone on refresh, so a stale
entry written by an older provider version drops out of state on the next `terraform plan`, and
`terraform destroy` completes cleanly (with a warning noting that no remote object was deleted).

If a rule of this type predates the backend removal, delete it by hand under
**Policies > Rules** in Sysdig Secure.

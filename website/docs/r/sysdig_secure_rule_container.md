---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_container"
sidebar_current: "docs-sysdig-secure-rule-container"
description: |-
  Creates a Sysdig Secure Container Rule.
---

# sysdig\_secure\_rule\_container

Creates a Sysdig Secure Container Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_rule_container" "sample" {
  name = "Nginx container spawned"
  description = "A container withthe nginx image spawned in the cluster."
  tags = ["container", "cis"]

  matching = true // default
  containers = ["nginx"]
}
```

## Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

### Matching

* `matching` - (Optional) Defines if the image name matches or not with the provided list. Default is true.
* `containers` - (Required) List of containers to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.
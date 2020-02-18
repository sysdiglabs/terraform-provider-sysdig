---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_process"
sidebar_current: "docs-sysdig-secure-rule-process"
description: |-
  Creates a Sysdig Secure Process Rule.
---

# sysdig\_secure\_rule\_process

Creates a Sysdig Secure Process Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

## Example usage

```hcl
resource "sysdig_secure_rule_process" "sample" {
  name = "Launch Suspicious Network Tool in Container" // ID
  description = "Detect network tools launched inside container"

  matching = true // default
  processes = ["nc", "ncat", "nmap", "dig", "tcpdump", "tshark", "ngrep"]
}

```

## Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

### Matching

* `matching` - (Optional) Defines if the process name matches or not with the provided list. Default is true.
* `processes` - (Required) List of processes to match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.
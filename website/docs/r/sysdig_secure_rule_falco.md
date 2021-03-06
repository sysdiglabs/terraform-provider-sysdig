---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_falco"
sidebar_current: "docs-sysdig-secure-rule-falco"
description: |-
  Creates a Sysdig Secure Falco Rule.
---

# sysdig\_secure\_rule\_falco

Creates a Sysdig Secure Falco Rule.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_rule_falco" "example" {
  name = "Terminal shell in container" // ID
  description = "A shell was used as the entrypoint/exec point into a container with an attached terminal."
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%user.name %container.info shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline terminal=%proc.tty container_id=%container.id image=%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.

- - -

### Conditions

* `condition` - (Required) A [Falco condition](https://falco.org/docs/rules/) is simply a Boolean predicate on Sysdig events expressed using the Sysdig [filter syntax](http://www.sysdig.org/wiki/sysdig-user-guide/#filtering) and macro terms. 
* `output` - (Optional) Add additional information to each Falco notification's output. Required if append is false.
* `priority` - (Optional) The priority of the Falco rule. It can be: "emergency", "alert", "critical", "error", "warning", "notice", "info" or "debug". By default is "warning".
* `source` - (Optional) The source of the event. It can be either "syscall", "k8s_audit" or "aws_cloudtrail". Required if append is false.
* `append` - (Optional) This indicates that the rule being created appends the condition to an existing Sysdig-provided rule. By default this is false. Appending to user-created rules is not supported by the API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

## Import

Secure Falco runtime rules can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_rule_falco.example 12345
```
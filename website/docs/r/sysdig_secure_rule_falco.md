---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_rule_falco"
description: |-
  Creates a Sysdig Secure Falco Rule.
---

# Resource: sysdig\_secure\_rule\_falco

Creates a Sysdig Secure Falco Rule.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_rule_falco" "example" {
  name        = "Terminal shell in container" // ID
  description = "A shell was used as the entrypoint/exec point into a container with an attached terminal."
  tags        = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output    = "A shell was spawned in a container with an attached terminal (user=%user.name %container.info shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline terminal=%proc.tty container_id=%container.id image=%container.image.repository)"
  priority  = "notice"
  source    = "syscall" // syscall or k8s_audit


  exceptions {
    name   = "proc_names"
    fields = ["proc.name"]
    comps  = ["in"]
    values = jsonencode(["python", "python2", "python3"]) # If only one element is provided, do not specify it a list of lists.
  }

  exceptions {
    name   = "container_proc_name"
    fields = ["container.id", "proc.name"]
    comps  = ["=", "in"]
    values = jsonencode([ # If more than one element is provided, you need to specify a list of lists.
      ["host", ["docker_binaries", "k8s_binaries", "lxd_binaries", "nsenter"]]
    ])
  }

  exceptions {
    name   = "proc_cmdline"
    fields = ["proc.name", "proc.cmdline"]
    comps  = ["in", "contains"]
    values = jsonencode([ # In this example, we are providing a pair of values for proc_cmdline, each one in a line.
      [["python", "python2", "python3"], "/opt/draios/bin/sdchecks"],
      [["java"], "sdjagent.jar"]
    ])
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Optional) The description of Secure rule. By default is empty.
* `tags` - (Optional) A list of tags for this rule.
* `condition` - (Required) A [Falco condition](https://falco.org/docs/rules/) is simply a Boolean predicate on Sysdig events expressed using the Sysdig [filter syntax](http://www.sysdig.org/wiki/sysdig-user-guide/#filtering) and macro terms. 
* `output` - (Optional) Add additional information to each Falco notification's output. Required if append is false.
* `priority` - (Optional) The priority of the Falco rule. It can be: "emergency", "alert", "critical", "error", "warning", "notice", "info" or "debug". By default is "warning".
* `source` - (Optional) The source of the event. It can be either "syscall", "k8s_audit" or "aws_cloudtrail". Required if append is false.
* `exceptions` - (Optional) The exceptions key is a list of identifier plus list of tuples of filtercheck fields. See below for details.
* `append` - (Optional) This indicates that the rule being created appends the condition to an existing Sysdig-provided rule. By default this is false. Appending to user-created rules is not supported by the API.

### Exceptions

Starting in 0.28.0, Falco supports an optional exceptions property to rules. The exceptions key is a list of identifier plus list of tuples of filtercheck fields.
For more information about the syntax of the exceptions, check the [official Falco documentation](https://falco.org/docs/rules/exceptions/).

Supported fields for exceptions:

* `name` - (Required) The name of the exception. Only used to provide a handy name, and to potentially link together values in a later rule that has `append = true`.
* `fields` - (Required) Contains one or more fields that will extract a value from the syscall/k8s_audit events.
* `comps` - (Required) Contains comparison operators that align 1-1 with the items in the fields property.
* `values` - (Required) Contains tuples of values. Each item in the tuple should align 1-1 with the corresponding field and comparison operator. Since the value can be a string, a list of strings or a list of a list of strings, the value of this field must be supplied in JSON format. You can use the default `jsonencode` function to provide this value. See the usage example on the top.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

## Import

Secure Falco runtime rules can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_rule_falco.example 12345
```
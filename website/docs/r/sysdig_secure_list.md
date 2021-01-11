---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_list"
sidebar_current: "docs-sysdig-secure-list"
description: |-
  Creates a Sysdig Secure Falco List.
---

# sysdig\_secure\_list

Creates a Sysdig Secure Falco List.

`~> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.`

## Example usage

```hcl
resource "sysdig_secure_list" "allowed_dev_files" {
  name = "allowed_dev_files"
  items = ["/dev/null", "/dev/stdin", "/dev/stdout", "/dev/stderr", "/dev/random", 
           "/dev/urandom", "/dev/console", "/dev/kmsg"]
  append = true # default: false
}
```

## Argument Reference

* `name` - (Required) The name of the Secure list. It must be unique if it's not in append mode.

* `items` - (Required) Elements in the list. Elements can be another lists.

* `append` - (Optional)  Adds these elements to an existing list. Used to extend existing lists provided by Sysdig.
    The rules can only be extended once, for example if there is an existing list called "foo", one can have another 
    append rule called "foo" but not a second one. By default this is false.

## Import

Secure lists can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_list.example 12345
```
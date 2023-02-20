---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_scanning_policy_assignment"
description: |-
  Creates a Sysdig Secure Scanning Policy Assignment for Legacy Scanning Engine.
---

# Resource: sysdig_secure_scanning_policy_assignment

Creates a Sysdig Secure Policy Assignment.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.  

## Example Usage

```terraform
resource "sysdig_secure_scanning_policy_assignment" "assignment_example" {
  items {
    name = "myassignment1"
    image {
      type = "tag"
      value = "latest"
    }
    registry = "docker.io"
    repository = "example"

    policy_ids = ["default"]
  }

  items {
    name = ""
    image {
      type = "tag"
      value = "latest"
    }
    registry = "*"
    repository = "*"

    policy_ids = [sysdig_secure_scanning_policy.scanning_policy_example.id]
    whitelist_ids = []
  }

  items {
    name = "default"
    image {
      type = "tag"
      value = "*"
    }
    registry = "*"
    repository = "*"

    policy_ids = [sysdig_secure_scanning_policy.scanning_policy_example.id, "default"]
  }
  
}
```

## Argument Reference

* `items` - (Required) List of scanning policy assignments. **Priority is defined from top to bottom with the order of the items**.

* `policy_bundle_id` - (Optional) Bundle for the policy assignment. The only value accepted is "default".

## Items block

* `name` - (Optional) The name of the Secure scanning policy assignment.

* `registry` - (Required) Any registry domain (e.g. quay.io). Wildcards are supported; an asterisk * specifies any registry.

* `repository` - (Required) Any repository (typically = name of the image). Wildcards are supported; an asterisk * specifies any repository.

* `image` - (Required) Block to define the image tag.

* `policy_ids` - (Required) Scanning policy IDs assigned to the given Registry/Repository:tag. At least 1 required.

* `whitelist_ids` - (Optional) List of vulnerability exception list associated with the assignment.

- - -

### Image block

* `type` - equal always to "tag"

* `value` - Image tag, any tag. Wildcards are supported; an asterisk * specifies any tag.

- - -

## Attributes Reference

No additional attributes are exported.

## Import

Secure scanning policies can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_scanning_policy_assignment.example default
```

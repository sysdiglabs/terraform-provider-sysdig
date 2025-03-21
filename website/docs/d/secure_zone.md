---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_zone"
description: |-
  Retrieves Zone by ID.
---

# sysdig\_secure\_zone Data Source

The `sysdig_secure_zone` data source allows you to retrieve information about a specific Sysdig Secure Zone.

## Example Usage

```hcl
resource "sysdig_secure_zone" "sample" {
  name        = "test-secure-zone"
  description = "Test secure zone"
  scope {
    target_type = "aws"
    rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
  }
}

data "sysdig_secure_zone" "test" {
  depends_on = [sysdig_secure_zone.sample]
  name       = sysdig_secure_zone.sample.name
}
```

## Argument Reference

The following arguments are supported, it is required that one of them is provided:

- `name` - The name of the Sysdig Secure Zone.
- `id` - The ID of the Sysdig Secure Zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `is_system` - (Computed) Whether the Zone is a system zone.
- `author` - (Computed) The zone author.
- `scope` - (Computed) The scope of the zone.
- `last_modified_by` - (Computed) By whom is last modification made.
- `last_updated` - (Computed) Timestamp of last modification of zone.

## Import

Zone can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_zone.example 12345
```

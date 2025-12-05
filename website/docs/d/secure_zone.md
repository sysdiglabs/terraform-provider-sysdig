---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_zone"
description: |-
  Retrieves Zone by ID or name.
---

# sysdig\_secure\_zone Data Source

The `sysdig_secure_zone` data source allows you to retrieve information about a specific Sysdig Secure Zone.

-> **Note:** The `rules` attribute supports both v2 and legacy (v1) syntax. Legacy v1 syntax (`labels`, `labelValues`, `agentTags`) is deprecated — use `expression` blocks instead. Rules using v2-compatible field names are fully supported. See the [resource documentation](../r/secure_zone.md) for migration guidance.

## Example Usage

### With expression-based scopes (recommended)

```hcl
resource "sysdig_secure_zone" "sample" {
  name        = "test-secure-zone"
  description = "Test secure zone"

  scope {
    target_type = "aws"

    expression {
      field    = "organization"
      operator = "in"
      values   = ["o1", "o2"]
    }

    expression {
      field    = "account"
      operator = "in"
      values   = ["a1", "a2"]
    }
  }
}

data "sysdig_secure_zone" "test" {
  depends_on = [sysdig_secure_zone.sample]
  name       = sysdig_secure_zone.sample.name
}
```

### With legacy `rules` scopes (deprecated)

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
- `last_modified_by` - (Computed) By whom is last modification made.
- `last_updated` - (Computed) Timestamp of last modification of zone.
- `scope` - (Computed) The scope of the zone. Each scope contains:
    - `id` - The ID of the scope.
    - `target_type` - The resource type this scope applies to (e.g., `aws`, `gcp`, `azure`, `kubernetes`, `host`, `image`, `git`, `ibm`, `oci`).
    - `rules` - (Computed) Query language expression. Legacy v1 syntax is deprecated; v2-compatible syntax is fully supported.
    - `expression` - List of filter expressions, each containing:
        - `field` - Field name to filter on (e.g., `organization`, `account`, `label.<key>`, `agent.tag.<key>`).
        - `operator` - Operator applied (e.g., `in`, `contains`, `not_in`, `not_contains`, `is_not`).
        - `value` - Single value (for operators like `contains`).
        - `values` - List of values (for operators like `in`).

### Understanding scope structure

- **Within a single scope**: all `expression` blocks are combined with **AND**.
- **Between scopes**: multiple scopes are combined with **OR**.
- **Within `in` operator**: values are combined with **OR**.

### Expression fields (v2)

The v2 expression model uses explicit field names for labels and agent tags:

- `label.<labelKey>` - Filters on label values (replaces legacy `labels` and `labelValues`).
- `agent.tag.<tagKey>` - Filters on agent tag values (replaces legacy `agentTags`).

For detailed migration examples from `rules` to `expression`, see the [resource documentation](../r/secure_zone.md#migrating-from-rules-to-expression).

-> **Note:** The data source returns **both** `rules` and `expression` for each scope when available. This differs from the resource, where they are mutually exclusive. The `rules` field contains the v1/v2 string representation, while `expression` contains the structured equivalent. Use whichever representation is more convenient for your use case.

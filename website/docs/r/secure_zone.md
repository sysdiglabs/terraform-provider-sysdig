---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_zone"
description: |-
  Creates a Sysdig Secure Zone.
---

# Resource: sysdig_secure_zone

Creates a Sysdig Secure Zone.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

### Expression-based scopes (recommended)

```terraform
resource "sysdig_secure_zone" "example" {
  name        = "example-zone"
  description = "An example Sysdig zone"

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

  scope {
    target_type = "azure"

    expression {
      field    = "organization"
      operator = "contains"
      value    = "o1"
    }
  }
}
```

### Legacy rule-based scopes (deprecated)

```terraform
resource "sysdig_secure_zone" "legacy" {
  name        = "example-zone-legacy"
  description = "Legacy rules test"

  scope {
    target_type = "kubernetes"
    rules       = "agentTags != \"environment: production\" and not agentTags contains \"team: platform\""
  }
}
```

## Argument Reference

- `name` - (Required) The name of the Zone.
- `description` - (Optional) The description of the Zone.
- `scope` - (Required) One or more `scope` blocks attached to the Zone.

### Scope block

A `scope` defines what resources belong to this zone.

- `id` - (Computed) The ID of the scope.

- `target_type` - (Required) The resource type this scope applies to. Supported types:

    - AWS - `aws`
    - GCP - `gcp`
    - Azure - `azure`
    - Kubernetes - `kubernetes`
    - Image - `image`
    - Host - `host`
    - Git - `git`
    - IBM - `ibm`
    - OCI - `oci`

- `rules` - (Optional) Query language expression for filtering results.

  ~> **Note:** The `rules` field supports both v2 and legacy (v1) syntax. When using legacy v1 attributes (`labels`, `labelValues`, `agentTags`), a deprecation warning will be shown — migrate to `expression` blocks with v2 field names (`label.<key>`, `agent.tag.<key>`). Rules using v2-compatible syntax (e.g., `organization`, `account`, `cluster`) are fully supported and produce no warning. `rules` and `expression` cannot be used together within the same `scope`.

- `expression` - One or more blocks that define the scope as a list of filter expressions.

  A scope must specify either `rules` or at least one `expression` block.

#### Expression block

Each `expression` block represents a single condition.

- `field` - (Required) Field name to filter on. See the "Supported fields" section below.
- `operator` - (Required) Operator to apply.
- `value` - (Optional) Single value for operators that take one argument.
- `values` - (Optional) List of values for operators such as `in`.

~> **Note:** Provide either `value` or `values` for an `expression` block (depending on the operator). If both are set, `values` takes precedence.

## Migrating from `rules` to `expression`

The `rules` attribute is deprecated and will be removed in a future version. New zones should be created using `expression` blocks. Existing zones that use `rules` can be migrated.

### What to expect during migration

- Migration is done by updating your Terraform configuration.
- An **update in place** is expected. Terraform may show changes under `scope` because the representation changes from a single `rules` string to structured `expression` blocks.
- Within a single `scope`, `rules` and `expression` are **mutually exclusive**.

### Understanding scope logic

To migrate correctly, you must understand how expressions combine:

- **Within a single scope**: all `expression` blocks are combined with **AND**.
- **Between scopes**: multiple scopes are combined with **OR**.
- **Within `in` operator**: values are combined with **OR** (e.g., `field in ("a", "b")` means `a OR b`).

This means that some legacy rules that use `in` with multiple values for different keys will need to be **split into multiple scopes** to preserve the same semantic behavior.

### Semantic change: labels / agentTags

In legacy `rules`, fields like `labels` (cloud accounts), `labelValues` (kubernetes), and `agentTags` (kubernetes/host) were represented as a single string in the form `"key: value"`.
This had important side effects:

- `labels in ("key: value")` filtered on both key and value combined.
- `labels contains "e"` could match the `e` in either the key or the value.

In the v2 expression model this ambiguity is removed.
Instead of querying the combined `"key: value"` string, you explicitly select the key in the field name and filters apply to the **value only**:

- `label.<labelKey>` (replaces legacy `labels` and `labelValues`)
- `agent.tag.<tagKey>` (replaces legacy `agentTags`)

### Migration strategy

1. Pick one `scope` block at a time.
2. Translate each condition in `rules` into one or more `expression` blocks.
3. Replace the `rules = "..."` line with `expression { ... }` blocks.
4. If your legacy rule uses `in` with multiple **different keys**, split into multiple scopes (see examples below).
5. Run `terraform plan` and verify the plan shows an in-place update.

### Example migrations

#### Simple case: single key

When all values in an `in` clause share the same key, they stay in a single scope:

Legacy:

```hcl
rules = "agentTags in (\"cluster: auto-do-not-delete-683\", \"cluster: qa-integrations\")"
```

Migrated (single scope):

```terraform
scope {
  target_type = "kubernetes"

  expression {
    field    = "agent.tag.cluster"
    operator = "in"
    values   = ["auto-do-not-delete-683", "qa-integrations"]
  }
}
```

#### Multiple different keys: split into multiple scopes

When an `in` clause contains values with **different keys**, they must be split into separate scopes to preserve the OR semantics:

Legacy:

```hcl
rules = "agentTags in (\"env: prod\", \"region: us-west\")"
```

Migrated (two scopes, combined with OR):
```terraform
scope {
  target_type = "kubernetes"

  expression {
    field    = "agent.tag.env"
    operator = "in"
    values   = ["prod"]
  }
}

scope {
  target_type = "kubernetes"

  expression {
    field    = "agent.tag.region"
    operator = "in"
    values   = ["us-west"]
  }
}
```

#### Mixed tags and static filters with same key

When tags/labels share the same key, they can be combined in a single scope:

Legacy:

```hcl
rules = "agentTags in (\"team: a\", \"team: b\") and labelValues in (\"env: dev\", \"env: staging\") and namespace in (\"core\") and clusterId in (\"dev\")"
```

Migrated (single scope):

```terraform
scope {
  target_type = "kubernetes"

  expression {
    field    = "agent.tag.team"
    operator = "in"
    values   = ["a", "b"]
  }

  expression {
    field    = "label.env"
    operator = "in"
    values   = ["dev", "staging"]
  }

  expression {
    field    = "namespace"
    operator = "in"
    values   = ["core"]
  }

  expression {
    field    = "clusterId"
    operator = "in"
    values   = ["dev"]
  }
}
```

### Operator mapping notes

The legacy query language and the new structured expressions use different operator spellings.

Common patterns:

- `field in ("a", "b")` → `operator = "in"` with `values = ["a", "b"]`
- `field contains "x"` → `operator = "contains"` with `value = "x"`
- `not field contains "x"` → `operator = "not_contains"` with `value = "x"`
- `field != "x"` → `operator = "is_not"` with `value = "x"`

~> **Note:** Operators with multiple words use underscore notation: `is_not`, `not_contains`, `not_in`. The backend normalizes space-separated forms (e.g., `"is not"`) to underscores.

## Supported fields and legacy query language notes

When using `rules` (**deprecated**), the following operators are supported:

- `and` logical operators
- `in`
- `contains` to check partial values of attributes

When using `expression`, you specify each condition in a dedicated block and the provider translates it to the backend model.

### Legacy-only fields (`rules` only)

The following fields are supported only by the deprecated `rules` syntax and are **not** available in `expression`.

- `labels` (cloud target types)
- `labelValues` (kubernetes)
- `agentTags` (kubernetes and host)

Use the v2 expression fields described below instead.

### Expression fields (v2)

The v2 expression model avoids ambiguous matching by requiring the label/tag key to be encoded in the `field`:

- `label.<labelKey>` (replaces legacy `labels` and `labelValues`)
- `agent.tag.<tagKey>` (replaces legacy `agentTags`)

### Supported fields by target type (legacy `rules` reference)

The following list is kept as a reference for the legacy `rules` language.

-> **Forward compatibility:** If the backend introduces new fields that the provider does not yet recognize, they are silently accepted during validation. Only fields the provider knows about are checked against the target_type allowlist.

- `aws`:
    - `account`
        - Type: string
        - Description: AWS account ID
        - Example query: `account in ("123456789012")`
    - `organization`
        - Type: string
        - Description: AWS organization ID
        - Example query: `organization in ("o-1234567890")`
    - `labels`
        - Type: string
        - Description: AWS account labels (legacy `rules` only)
        - Example query: `labels in ("key: value")`
    - `location`
        - Type: string
        - Description: AWS account location
        - Example query: `location in ("us-east-1")`
- `gcp`:
    - `account`
        - Type: string
        - Description: GCP account ID
        - Example query: `account in ("123456789012")`
    - `organization`
        - Type: string
        - Description: GCP organization ID
        - Example query: `organization in ("1234567890")`
    - `labels`
        - Type: string
        - Description: GCP account labels (legacy `rules` only)
        - Example query: `labels in ("key: value")`
    - `location`
        - Type: string
        - Description: GCP account location
        - Example query: `location in ("us-east-1")`
- `azure`:
    - `account`
        - Type: string
        - Description: Azure account ID
        - Example query: `account in ("123456789012")`
    - `organization`
        - Type: string
        - Description: Azure organization ID
        - Example query: `organization in ("1234567890")`
    - `labels`
        - Type: string
        - Description: Azure account labels (legacy `rules` only)
        - Example query: `labels in ("key: value")`
    - `location`
        - Type: string
        - Description: Azure account location
        - Example query: `location in ("us-east-1")`
- `kubernetes`:
    - `clusterId`
        - Type: string
        - Description: Kubernetes cluster ID
        - Example query: `clusterId in ("cluster")`
    - `namespace`
        - Type: string
        - Description: Kubernetes namespace
        - Example query: `namespace in ("namespace")`
    - `labelValues`
        - Type: string
        - Description: Kubernetes label values (legacy `rules` only)
        - Example query: `labelValues in ("label1")`
    - `distribution`
        - Type: string
        - Description: Kubernetes distribution
        - Example query: `distribution in ("eks")`
    - `agentTags`
        - Type: string
        - Description: Agent tags in the form `"key: value"` (legacy `rules` only)
        - Example query: `agentTags contains "key: value"`
- `host`:
    - `clusterId`
        - Type: string
        - Description: Kubernetes cluster ID
        - Example query: `clusterId in ("cluster")`
    - `name`
        - Type: string
        - Description: Host name
        - Example query: `name in ("host")`
    - `agentTags`
        - Type: string
        - Description: Agent tags in the form `"key: value"` (legacy `rules` only)
        - Example query: `agentTags contains "key: value"`
- `image`:
    - `registry`
        - Type: string
        - Description: Image registry
        - Example query: `registry in ("registry")`
    - `repository`
        - Type: string
        - Description: Image repository
        - Example query: `repository in ("repository")`
- `git`:
    - `gitIntegrationId`
        - Type: string
        - Description: Git integration ID
        - Example query: `gitIntegrationId in ("gitIntegrationId")`
    - `gitSourceId`
        - Type: string
        - Description: Git source ID
        - Example query: `gitSourceId in ("gitSourceId")`
- `ibm`:
    - `account`
        - Type: string
        - Description: IBM account ID
        - Example query: `account in ("123456789012")`
    - `organization`
        - Type: string
        - Description: IBM organization ID
        - Example query: `organization in ("1234567890")`
    - `labels`
        - Type: string
        - Description: IBM account labels (legacy `rules` only)
        - Example query: `labels in ("key: value")`
    - `location`
        - Type: string
        - Description: IBM account location
        - Example query: `location in ("us-east-1")`
    - `resourceGroupId`
        - Type: string
        - Description: IBM resource group ID
        - Example query: `resourceGroupId in ("rg-1234")`
    - `accountGroupId`
        - Type: string
        - Description: IBM account group ID
        - Example query: `accountGroupId in ("ag-1234")`
    - `accountGroupName`
        - Type: string
        - Description: IBM account group name
        - Example query: `accountGroupName in ("my-group")`
- `oci`:
    - `account`
        - Type: string
        - Description: OCI account ID
        - Example query: `account in ("ocid1.tenancy.oc1..example")`
    - `organization`
        - Type: string
        - Description: OCI organization ID
        - Example query: `organization in ("1234567890")`
    - `labels`
        - Type: string
        - Description: OCI account labels (legacy `rules` only)
        - Example query: `labels in ("key: value")`
    - `location`
        - Type: string
        - Description: OCI account location
        - Example query: `location in ("us-ashburn-1")`

**Note**: Whenever filtering for values with special characters, the values need to be encoded.
When `"` or `\` are the special characters, they need to be escaped with `\` and then encoded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - (Computed) The ID of the Zone.
- `is_system` - (Computed) Whether the Zone is a system zone.
- `author` - (Computed) The zone author.
- `last_modified_by` - (Computed) By whom is last modification made.
- `last_updated` - (Computed) Timestamp of last modification of zone.

## How state is managed (drift prevention)

When reading a zone from the API, the provider preserves the representation format from your configuration:

- If your config uses `expression` blocks, the state stores expressions.
- If your config uses `rules`, the state stores the rules string.

This prevents perpetual plan diffs when the backend returns both representations. On import (where there is no prior config), the state defaults to `rules`.

## Import

Zone can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_zone.example 12345
```

~> **Note:** Imported zones are always represented using the `rules` string in initial state. If your configuration uses `expression` blocks, the first `terraform plan` after import will show changes to converge to the expression-based representation. Apply once to align state with your config.

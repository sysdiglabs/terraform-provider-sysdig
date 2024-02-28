---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_zone"
description: |-
  Creates a Sysdig Secure Posture Zone.
---

# Resource: sysdig_secure_posture_zone

Creates a Sysdig Secure Posture Zone.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data "sysdig_secure_posture_policies" "all" {}

locals {
  fedramp_policies = [
    for p in data.sysdig_secure_posture_policies.all.policies :
    p if length(regexall(".*FedRAMP.*", p.name)) > 0
  ]
}

resource "sysdig_secure_posture_zone" "example" {
  name        = "Zone with FedRAMP policies"
  description = "Zone description"
  policy_ids  = [for p in local.fedramp_policies : p.id]

  scopes {
    scope {
      target_type = "aws"
      rules       = "organization in (\"o1\", \"o2\") and account in (\"a1\", \"a2\")"
    }

    scope {
      target_type = "azure"
      rules       = "organization contains \"o1\""
    }
  }
}
```

## Argument Reference

- `name` - (Required) The name of the Posture Zone.
- `description` - (Optional) The description of the Posture Zone.
- `policy_ids` - (Optional) The list of Posture Policy IDs attached to Zone.
- `scopes` - (Optional) Scopes block defines list of scopes attached to Zone.

### Scopes block

- `target_type` - (Required) The target type for the scope. Supported types:

  - AWS - `aws`
  - GCP - `gcp`
  - Azure - `azure`
  - Kubernetes - `kubernetes`
  - Image - `image`
  - Host - `host`
  - Git - `git`

- `rules` - (Optional) Query language expression for filtering results. Empty rules means no filtering.

  Operators:

  - `and`, `or` logical operators
  - `in`
  - `contains` to check partial values of attributes

  List of supported fields by target type:

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
      - Description: AWS account labels
      - Example query: `labels in ("label1")`
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
      - Description: GCP account labels
      - Example query: `labels in ("label1")`
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
      - Description: Azure account labels
      - Example query: `labels in ("label1")`
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
      - Description: Kubernetes label values
      - Example query: `labelValues in ("label1")`
    - `distribution`
      - Type: string
      - Description: Kubernetes distribution
      - Example query: `distribution in ("eks")`
    - `agentTags`
      - Type: string
      - Description: Tags of the Sysdig Agent running in Kubernetes
      - Example query: `agentTags in ("env:prod")` _(tag values are formatted as `key:value`)_
  - `host`:
    - `clusterId`
      - Type: string
      - Description: Kubernetes cluster ID
      - Example query: `clusterId in ("cluster")`
    - `name`
      - Type: string
      - Description: Host name
      - Example query: `name in ("host")`
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

  **Note**: Whenever filtering for values with special characters, the values need to be encoded.
  When “ or \ are the special characters, they need to be escaped with \ and then encoded.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The zone author.
- `last_modified_by` - (Computed) By whom is last modification made.
- `last_updated` - (Computed) Timestamp of last modification of zone.

## Import

Posture zone can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_zone.example 12345
```

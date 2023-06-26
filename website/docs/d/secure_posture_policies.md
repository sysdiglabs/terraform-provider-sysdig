---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_policies"
description: |-
  Retrieves all Posture policies.
---

# Data Source: sysdig_secure_posture_policies

Retrieves the information of all Posture policies.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data sysdig_secure_posture_policies policies {}
```

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policies` - List of posture policy objects.

### Policies reference

* `id` - Policy ID
* `name` - Policy name, eg. `CIS Docker Benchmark`
* `type` - Policy type as int value, can be one of the following:
  - 0 - UNKNOWN
  - 1 - KUBERNETES
  - 2 - DOCKER
  - 3 - LINUX
  - 4 - AWS
  - 5 - GCP
  - 6 - AZURE
* `kind` - Policy kind as int value, can be one of the following:
  - 0 - None
  - 1 - BestPractice
  - 2 - Compliance
  - 3 - Corporate
* `description` - Policy description, eg. `CIS Docker Benchmark`
* `version` - Policy version, eg. `1.0.0`
* `api_version` - Policy API version `1.0.0`
* `link` - Policy link
* `authors` - Policy authors, eg. `John Doe`
* `published_date` - Policy published date, eg. `1588617600000`
* `min_kube_version` - Policy minimum Kubernetes version, eg. `1.16`
* `max_kube_version` - Policy maximum Kubernetes version, eg. `1.18`
* `is_custom` - Policy is custom flag
* `is_active` - Policy is active flag
* `platform` - Policy platform, eg. `Kubernetes`
* `zones` - List of policy zones

### Zones reference

* `id` - Zone ID
* `name` - Zone Name, eg. `Entire Infrastructure`

---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_policy"
description: |-
  Retrieves Posture policy by ID.
---

# Data Source: sysdig_secure_posture_policies

Retrieves the information of all Posture policies.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
data sysdig_secure_posture_policies policy {
  id = "454678"
}
```

## Argument Reference

- `id` - (Required) The ID of the Posture Policy, eg. `2`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Posture Policy, eg. `452`
- `name` - The name of the Posture Policy, eg. `CIS Docker Benchmark`
- `description` - The description of the Posture Poliy,  eg. `CIS Docker Benchmark`
* `link` - Policy link
* `type` - Policy type:
  - AWS - `aws`
  - GCP - `gcp`
  - Azure - `azure`
  - Kubernetes - `kubernetes`
  - Linux - `linux`
  - Docker - `docker`
  - OCI = `oci`
* `min_kube_version` - Policy minimum Kubernetes version, eg. `1.24`
* `max_kube_version` - Policy maximum Kubernetes version, eg. `1.26`
* `is_active` - Policy is active flag (active means policy is published, not active means policy is draft). by default is true.
* `platform` - Policy platform: 
    - IKS -     `iks`,
    - GKE -     `gke`,
    - Vanilla -  `vanilla`,
    - AKS -     `aks`,
    - RKE2 -     `rke2`,
    - OCP4  -     `ocp4`,
    - MKE  -      `mke`,
    - EKS  -     `eks`,
* `groups` - Group block defines list of groups attached to Policy

### Groups block
- `id` - The ID of the Group, eg. `15000`
- `name` - The name of the Posture Policy Group.
- `description` - The description of the Posture Policy Group.
- `requirements` - Requirements block defines list of requirements attached to Group

### Requirements block
- `id` - The ID of the Requirement, eg. `15000`
- `name` - The name of the Posture Policy Requirement.
- `description` - The description of the Posture Policy Requirement.
- `controls` - Controls block defines list of controls linked to requirments

### Controls block
- `name` - The name of the Posture Control.
- `enabled` - The 'Control is enabled' flag indicates whether the control will affect the policy evaluation or not. By default, it is set to true

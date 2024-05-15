---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_policy"
description: |-
  Creates Sysdig Secure Posture Policy.
---

# Resource: sysdig_secure_posture_policy

Creates a Sysdig Secure Posture Policy.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_posture_policy" "example" {
    name = "demo policy"
    type = "kuberenetes"
    platform = "vanilla"
    max_kube_version = 2.0
    description = "demo create policy from terraform"
      group {
        name = "Security"
        description = "Security description"
        requirement{
          name = "Security Enforce access control"
          description = "Enforce description"
          control {
              name = "Create Pods"
              enabled = false
          }
          control {
              name = "Kubelet - Disabled AlwaysAllowed Authorization"
          }
        }
      }
      group {
          name = "Data protection"
          description = "Data protection description"
          requirement{
            name = "Enforce access control"
            description = "Enforce description"
            control {
                name = "Create Pods"
            }
            control {
                name = "Kubelet - Disabled AlwaysAllowed Authorization"
            }
          }     
      }
}
```

## Argument Reference

- `name` - (Required) The name of the Posture Policy, eg. `CIS Docker Benchmark`
- `description` - (Required) The description of the Posture Poliy,  eg. `CIS Docker Benchmark`
* `link` -  (Optional) Policy link
* `type` -  (Optional) Policy type:
  - AWS - `aws`
  - GCP - `gcp`
  - Azure - `azure`
  - Kubernetes - `kubernetes`
  - Linux - `linux`
  - Docker - `docker`
  - OCI = `oci`
* `min_kube_version` -  (Optional) Policy minimum Kubernetes version, eg. `1.24`
* `max_kube_version` -  (Optional) Policy maximum Kubernetes version, eg. `1.26`
* `is_active` -  (Optional) Policy is active flag (active means policy is published, not active means policy is draft). by default is true.
* `platform` - (Optional) Policy platform: 
    - IKS -     `iks`,
    - GKE -     `gke`,
    - Vanilla -  `vanilla`,
    - AKS -     `aks`,
    - RKE2 -     `rke2`,
    - OCP4  -     `ocp4`,
    - MKE  -      `mke`,
    - EKS  -     `eks`,
* `groups` - (Optional) Group block defines list of groups attached to Policy

### Groups block
- `name` - (Required) The name of the Posture Policy Group.
- `description` - (Required) The description of the Posture Policy Group.
- `requirements` -  (Optional) Requirements block defines list of requirements attached to Group

### Requirements block
- `name` - (Required) The name of the Posture Policy Requirement.
- `description` - (Required) The description of the Posture Policy Requirement.
- `controls` -  (Optional) Controls block defines list of controls linked to requirments

### Controls block
- `name` - (Required) The name of the Posture Control.
- `enabled` - (Optional) The 'Control is enabled' flag indicates whether the control will affect the policy evaluation or not. By default, it is set to true

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The zone author.

## Import

Posture policy can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_policy.example p
```

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
  name             = "demo policy"
  type             = "kubernetes"
  platform         = "Vanilla" // Currently supported, but will be deprecated in the future
  min_kube_version = 1.5       // Currently supported, but will be deprecated in the future
  max_kube_version = 2.0       // Currently supported, but will be deprecated in the future
  description      = "demo create policy from terraform"

  // New targets field to specify version constraints
  target
    {
      platform   = "Vanilla"
      minVersion = 1.5
      maxVersion = 2.0
    }

  group {
    name        = "Security"
    description = "Security description"

    requirement {
      name        = "Security Enforce access control"
      description = "Enforce description"

      control {
        name    = "Create Pods"
        enabled = false
      }

      control {
        name = "Kubelet - Disabled AlwaysAllowed Authorization"
      }
    }
  }

  group {
    name        = "Data protection"
    description = "Data protection description"

    requirement {
      name        = "Enforce access control"
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
  - OCI - `oci`
 * `platform`: (Optional) Platform for which the policy applies. This field will be deprecated in the future, and you should use the targets field instead to describe policy platform and version. Supported platforms include:

    IKS - iks
    GKE - gke
    Vanilla - vanilla
    AKS - aks
    RKE2 - rke2
    OCP4 - ocp4
    MKE - mke
    EKS - eks
    OCI - oci

* `minKubeVersion`: (Optional) Policy minimum Kubernetes version, e.g., 1.24. This field will be deprecated in the future, and you should use the targets field instead to describe policy platform and version.

* `maxKubeVersion`: (Optional) Policy maximum Kubernetes version, e.g., 1.26. This field will be deprecated in the future, and you should use the targets field instead to describe policy platform and version.

* `target`:(Optional) Specifies target platforms and version ranges. This field should replace Platform, MinKubeVersion, and MaxKubeVersion for more flexible and detailed policy descriptions.

  Note: The fields Platform, MinKubeVersion, and MaxKubeVersion will be deprecated in the future. We recommend using the targets field now to describe policy platform and version constraints

* `group` - (Optional) Group block defines list of groups attached to Policy

### Targets block
 - `platform` (Optional): Name of the target platform (e.g., IKS, AWS).
 - `minVersion` (Optional): Minimum version of the platform.(e.g., 1.24)
 - `maxVersion` (Optional): Maximum version of the platform. (e.g., 1.26)

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

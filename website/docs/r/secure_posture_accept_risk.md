---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_accept_risk"
description: |-
  Accepts Sysdig Secure Posture Risk.
---

# Resource: sysdig_secure_posture_accept_risk

Creates a Sysdig Secure Posture Accept Risk.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_posture_accept_risk" "accept_risk_global" {
    description  = "Accept risk for zone"
    control_name = "Network - Enabled Endpoint Private Access in Existing Clusters (EKS)"
    reason       = "Risk Transferred"
    expires_in   = "30 Days"
    zone_name = "Entire Infrastructure"
}

resource "sysdig_secure_posture_accept_risk" "accept_risk_resource" {
    description  = "Accept risk for resource"
    control_name = "Fargate - Untrusted Workloads"
    reason       = "Risk Transferred"
    expires_in   = "30 Days"
    filter       = "name in (\"aws-int-01-cicd-aws-eks-workloads-shield\") and providerType in (\"AWS\") and kind in (\"AWS_EKS_CLUSTER\") and location in (\"us-east-1\")"
}


resource "sysdig_secure_posture_accept_risk" "scheduler_set_to_loopback_bind_address" {
    description = "This is custom risk acceptance for scheduler_set_to_loopback_bind_address"
    control_name = "Scheduler - Set to Loopback bind-address"
    reason = "Custom"
    expires_in = "Custom"
    end_time = "1730293523000"
    zone_name = "Entire Infrastructure"
}
```

## Argument Reference

- `id` - (Computed) The unique identifier for the risk acceptance.
- `control_name` - (Required) The name of the posture control being accepted.
- `zone_name` - (Optional) The zone associated with the risk acceptance.
- `description` - (Required) A description of the risk acceptance.
- `filter` - (Optional) A filter for identifying the resources affected by the acceptance.
  
   ##### List of supported fields:
   - name
      - Type: string
      - Example: name in ("cf-templates-1s951ca3qbh1-us-west-2")
      - Description: The name of the resource to accept risk for
  
  - namespace
      - Type: string
      - Example: namespace in ("my-namespace")
      - Description: The namespace to accept risk for
  
   - kind
      - Type: string
      - Example: kind in ("AWS_S3_BUCKET")
      - Description: The resource kind to accept risk for

    - location
      - Type: string
      - Example: location in ("ap-southeast-2")
      - Description: The cloud location/region to accept risk for

    - providerType
      - Type: string
      - Example: providerType in ("AWS")
      - Description: The cloud provider to accept risk for (AWS/GCP/Azure)

- `reason` - (Required) The reason for accepting the risk. Possible values are:
  - `Risk Owned`
  - `Risk Transferred`
  - `Risk Avoided`
  - `Risk Mitigated`
  - `Risk Not Relevant`
  - `Custom`
- `expires_in` - (Required) The duration for which the risk acceptance is valid. Possible values are:
  - `7 Days`
  - `30 Days`
  - `60 Days`
  - `90 Days`
  - `Custom`
  - `Never`
- `expires_at` - (Computed) This timestamp indicates when the acceptance expires, formatted in UTC time (milliseconds since epoch).
- `end_time` - (Optional)  This timestamp indicates the custom time, when the acceptance expires, formatted in UTC time (milliseconds since epoch).
 If you choose expires_in=Custom, you should provide future end_time, which specifies the expiration date in milliseconds.
- `is_expired` - (Computed) Indicates whether the acceptance is expired.
- `acceptance_date` - (Computed) The date when the risk was accepted.
- `username` - (Computed) The username of the user who accepted the risk.
- `type` - (Computed) The type of risk acceptance.
- `is_system` - (Computed) Indicates whether the acceptance is sysdig-accepts.
- `accept_period` - (Computed) The period for which the risk is accepted.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The custom control author.

## Import

Posture accept risk can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_accept_risk.example c 12345
```

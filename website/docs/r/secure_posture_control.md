---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_posture_control"
description: |-
  Creates Sysdig Secure Posture Control.
---

# Resource: sysdig_secure_posture_control

Creates a Sysdig Secure Posture Control.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

## Example Usage

```terraform
resource "sysdig_secure_posture_control" "c"{
        name = "S3 - Enabled Versioning"
        description = "S3 - Enabled Versioning"
        resource_kind = "AWS_S3_BUCKET"
        severity = "Low"
        rego          = <<-EOF

            package sysdig

            import future.keywords.if
            import future.keywords.in

            default risky := false

            risky if {
              count(input.Versioning) == 0
            }

            risky if {
              some version in input.Versioning
              lower(version.Status) != "enabled"
            }
        EOF
     
     remediation_details = <<-EOF 
      **Using AWS CLI**\n1. Run **put-bucket-versioning** command (OSX/Linux/UNIX) using the name of the Amazon S3 bucket that you want to reconfigure as the identifier parameter, to enable S3 object versioning for the selected bucket. If the request is successful, the **put-bucket-versioning** command should not return an output:\n```bash\nawsaws s3api put-bucket-versioning\n  --bucket cc-prod-web-data\n  --versioning-configuration Status=Enabled\n```\n2. Repeat step no. 1 to enable S3 object versioning for other Amazon S3 buckets available within your AWS cloud account.
      
    EOF
}
```

## Argument Reference

- `name` - (Required) The name of the Posture Control. The name must be unique, e.g. `EC2 - Instances should not have a public IP address`
- `description` - (Required) The description of the Posture Control, eg. `EC2 - Instances should not have a public IP address`
- `rego` - (Required) The Posture control Rego. `package sysdig\ndefault risky = false\nrisky {\n    input.NetworkInterfaces[_].Association.PublicIp\n    input.      NetworkInterfaces[_].Association.PublicIp != \"\"\n}`
- `remediation_details`- (Required) The Posture control Remediation details. `Use a non-default VPC so that your instance is not assigned a public IP address by default`
- `resource_kind` - (Required) The resource type this control evaluates. Must be a supported resource kind string matching
  a resource type in the Sysdig CSPM inventory. The format varies by platform:

  - **AWS**: `AWS_S3_BUCKET`, `AWS_EC2_INSTANCE`, `AWS_IAM_ROLE`, `AWS_LAMBDA_FUNCTION`, ...
  - **GCP**: `GCP_STORAGE_GOOGLEAPIS_COM_BUCKET`, `GCP_COMPUTE_GOOGLEAPIS_COM_INSTANCE`, ...
  - **Azure**: `AZURE_MICROSOFT_COMPUTE_VIRTUALMACHINES`, `AZURE_MICROSOFT_STORAGE_STORAGEACCOUNTS`, ...
  - **Kubernetes**: `DEPLOYMENT`, `SERVICE`, `NAMESPACE`, `CLUSTERROLE`, ...
  - **IBM Cloud**: `IBM_USER-MANAGEMENT_USER`, `IBM_IS_VPC_INSTANCE`, `IBM_CLOUD-OBJECT-STORAGE_BUCKET`, ...
  - **Host** (Linux/Windows/Docker): `host`

  To list all valid values, query the CSPM API:
  ```
  GET /api/cspm/v1/policy/controls/resource-template/kinds
  ```
  See the [Sysdig API Swagger docs](https://docs.sysdig.com/en/docs/developer-tools/sysdig-api/#swagger-documentation) and
  the [posture controls API documentation](https://docs.sysdig.com/en/sysdig-secure/posture_controls/#sysdig-api-endpoint) for more details.
- `severity` - (Required) The Posture Control Severity [`High`, `Medium`, `Low`], case sensitive, e.g., `High`.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The custom control author.

## Import

Posture custom control can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_control.example 12345
```

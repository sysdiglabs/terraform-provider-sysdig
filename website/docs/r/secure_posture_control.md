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
- `resource_kind` - (Required) The Posture Control Resource kind. It should be a supported resource kind, eg. `AWS_S3_BUCKET` 
- `severity` - (Required) The Posture Control Severity [`High`, `Medium`, `Low`], case sensitive, e.g., `High`.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The custom control author.

## Import

Posture custom control can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_control.example c 12345
```

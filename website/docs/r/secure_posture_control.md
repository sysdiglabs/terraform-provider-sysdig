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
        name = "EC2 - Instances should not have a public IP address"
        description = "EC2 - Instances should not have a public IP address"
        rego = "package sysdig\ndefault risky = false\nrisky {\n    input.NetworkInterfaces[_].Association.PublicIp\n    input.NetworkInterfaces[_].Association.PublicIp != \"\"\n}"
        remediation_details= "Use a non-default VPC so that your instance is not assigned a public IP address by default.\n\nWhen you launch an EC2 instance into a default VPC, it is assigned a public IP address. When you launch an EC2 instance into a non-default VPC, the subnet configuration determines whether it receives a public IP address. The subnet has an attribute to determine if new EC2 instances in the subnet receive a public IP address from the public IPv4 address pool.\n\nYou cannot manually associate or disassociate an automatically-assigned public IP address from your EC2 instance. To control whether your EC2 instance receives a public IP address, do one of the following:\n\nModify the public IP addressing attribute of your subnet. For more information, see Modifying the public IPv4 addressing attribute for your subnet in the Amazon VPC User Guide.\nEnable or disable the public IP addressing feature during launch. This overrides the subnet's public IP addressing attribute. For more information, see Assign a public IPv4 address during instance launch in the Amazon EC2 User Guide for Linux Instances.\nFor more information, see Public IPv4 addresses and external DNS hostnames in the Amazon EC2 User Guide for Linux Instances.\n\nIf your EC2 instance is associated with an Elastic IP address, then your EC2 instance is reachable from the internet. You can disassociate an Elastic IP address from an instance or network interface at any time.\n\nTo disassociate an Elastic IP address\nOpen the Amazon EC2 console at https://console.aws.amazon.com/ec2/.\nIn the navigation pane, choose Elastic IPs.\nSelect the Elastic IP address to disassociate.\nFrom Actions, choose Disassociate Elastic IP address.\nChoose Disassociate."
        resource_kind = "AWS_S3_BUCKET"
        severity= "High"
}
```

## Argument Reference

- `name` -- (Required) The name of the Posture Control. The name must be unique, e.g. `EC2 - Instances should not have a public IP address`
- `description` - (Required) The description of the Posture Control, eg. `EC2 - Instances should not have a public IP address`
- `rego` - (Required) The Posture control Rego. `package sysdig\ndefault risky = false\nrisky {\n    input.NetworkInterfaces[_].Association.PublicIp\n    input.      NetworkInterfaces[_].Association.PublicIp != \"\"\n}`
- `remediation_details`- (Required) The Posture control Remediation details. `Use a non-default VPC so that your instance is not assigned a public IP address by default`
- `resource_kind` - (Required) The Posture Control Resource kind. It should be a supported resource kind, eg. `AWS_S3_BUCKET` 
- `severity` - (Required) The Posture Control Severity [`High`, `Medium`, `Low`], case sensitive, e.g., `High`.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `author` - (Computed) The Control author.

## Import

Posture policy can be imported using the ID, e.g.

```
$ terraform import sysdig_secure_posture_control.example c 12345
```

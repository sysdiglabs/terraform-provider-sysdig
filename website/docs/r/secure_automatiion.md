---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_automation"
description: "Creates a Sysdig Secure Automation"
---

# Resource: sysdig_secure_automation

Creates a Sysdig Secure Automation for automated response to security events.

-> **Note:** Sysdig Terraform Provider is under rapid development at this point. If you experience any issue or discrepancy while using it, please make sure you have the latest version. If the issue persists, or you have a Feature Request to support an additional set of resources, please open a [new issue](https://github.com/sysdiglabs/terraform-provider-sysdig/issues/new) in the GitHub repository.

### <i>Data for this resource is currently driven from data emanating primarily from the UI. Please only use this resource if you are comfortable with retrieving data from APIs</i>

## Example Usage
### Basic Email Automation

```terraform
resource "sysdig_secure_automation" "email_alerts" {
  name = "High Severity Email Alerts"
  
  automation_json = jsonencode({
    automation = {
      name    = "name-from-ui"  # This will be replaced with "High Severity Email Alerts"
      enabled = false
      version = "v1"
      
      nodes = {
        Send_Email_1 = {
          action = {
            type = "email"
            inputs = {
              channelId = "55179"
            }
          }
          outboundEdges = []
          onError       = []
        }
      }
      
      trigger = {
        on   = "new_findings"
        when = "finding.severity in (0, 1, 2, 3)"
        outboundEdges = [
          {
            node = "Send_Email_1"
          }
        ]
      }
    }
  })
}
```

## Argument Reference

* `name` - (Required) Name for the automation. This will override any name specified in the `automation_json`.  This is to ensure that the terraform objects name is the single source of truth

* `automation_json` - (Required) JSON configuration for the automation, typically exported from the Sysdig UI. The `name` field within this JSON will be replaced with the value from the `name` argument.  Test this in a lower environment to obtain without affecting your production environment.

## Usage Workflow

1. **Design in UI**: Use the Sysdig Secure automation builder to create and test your automation logic.
2. **Export JSON**: Copy the automation JSON from the browser's developer tools when the automation is saved, or retrieve it via the API.
3. **Configure Terraform**: Paste the JSON into the `automation_json` field and set a meaningful `name`.
4. **Apply**: Terraform will create and manage the automation, with the Terraform `name` taking precedence.

## Import

Automations can be imported using the automation ID, e.g.
```
$ terraform import sysdig_secure_automation.example 6a8eec76-1807-42d2-87e0-9468791cb03a
```

## Important Notes

* The `name` field in the Terraform configuration will always override the name in the `automation_json`. This ensures consistency and prevents drift.
* Complex automation logic (conditions, actions, triggers) should be designed in the Sysdig UI first, then exported to Terraform for state management.
* The automation JSON structure is flexible and supports all features available in the Sysdig UI automation builder.
* Changes to the automation logic require updating the `automation_json` field. Terraform will detect differences and update the automation accordingly.

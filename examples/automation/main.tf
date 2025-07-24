# examples/automation/main.tf
# Complete working example for testing the automation resource

terraform {
  required_version = ">= 0.13"
  required_providers {
    sysdig = {
      source  = "sysdiglabs/sysdig"
      version = ">= 0.4.0"
    }
  }
}

# Variables for configuration
variable "secure_api_token" {
  description = "Sysdig Secure API token"
  type        = string
  sensitive   = true
}

variable "secure_url" {
  description = "Sysdig Secure URL"
  type        = string
  default     = "https://secure.sysdig.com"
}

variable "notification_channel_id" {
  description = "ID of notification channel to use"
  type        = string
  default     = "12345"  # Replace with your actual channel ID
}

# Provider configuration
provider "sysdig" {
  sysdig_secure_api_token = var.secure_api_token
  sysdig_secure_url       = var.secure_url
}

# Data source to get current user info
data "sysdig_current_user" "current" {}

# Simple email automation
resource "sysdig_secure_automation" "simple_email" {
  name = "Simple Email Alert Automation"

  automation_json = jsonencode({
    automation = {
      name    = "UI Generated Name"  # This will be replaced
      enabled = false # Set this to true or false as appropriate
      version = "v1"

      nodes = {
        Send_Email_1 = {
          action = {
            type = "email"
            inputs = {
              channelId = var.notification_channel_id
            }
          }
          outboundEdges = []
          onError       = []
        }
      }

      trigger = {
        on   = "new_findings"
        when = "finding.severity in (0, 1, 2, 3, 4, 5, 6, 7)"
        outboundEdges = [
          {
            node = "Send_Email_1"
          }
        ]
      }
    }
  })
}

# Outputs for verification
output "current_user" {
  description = "Current user information"
  value = {
    id    = data.sysdig_current_user.current.id
    email = data.sysdig_current_user.current.email
  }
}

output "simple_automation" {
  description = "Simple automation details"
  value = {
    id         = sysdig_secure_automation.simple_email.automation_id
    name       = sysdig_secure_automation.simple_email.name
    enabled    = sysdig_secure_automation.simple_email.enabled
    created_at = sysdig_secure_automation.simple_email.created_at
  }
}

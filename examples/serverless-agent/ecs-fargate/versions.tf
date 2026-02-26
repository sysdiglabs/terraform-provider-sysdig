terraform {
  required_version = ">=1.7.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~>6.32.0"
    }
    sysdig = {
      source  = "sysdiglabs/sysdig"
      version = "~>3.4.0"
    }
  }
}
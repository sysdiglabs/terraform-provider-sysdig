terraform {
  required_version = ">=1.7.2"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.35.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.4.1"
    }
    sysdig = {
      source  = "sysdiglabs/sysdig"
      version = "~> 1.24.5"
    }
  }
}
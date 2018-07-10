---
layout: "sysdig"
page_title: "Provider: Sysdig"
sidebar_current: "docs-sysdig-index"
description: |-
  The Sysdig provider is used to interact with Sysdig products. The provider needs to be configured with proper API token before it can be used.
---

# Sysdig Provider

The Sysdig provider is used to interact with
[Sysdig Secure](https://sysdig.com/product/secure/) and
[Sysdig Monitor](https://sysdig.com/product/monitor/) products. The provider
needs to be configure with the proper API token before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
// Configure the Sysdig provider
provider "sysdig" {
  sysdig_secure_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

// Create a new secure policy
resource "sysdig_secure_policy" "unexpected_inbound_tcp_connection_traefik" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `sysdig_secure_api_token` - (Required) The Sysdig Secure API token, it must be
  present, but you can get it from the `SYSDIG_SECURE_API_TOKEN` environment variable.

* `sysdig_secure_url` - (Optional) This is the target Sysdig Secure base API
  endpoint. It's intended to be used with OnPrem installations. By defaults it
  points to `https://secure.sysdig.com`, and notice that should not be ended
  with an slash. It can also be sourced from the `SYSDIG_SECURE_URL` environment
  variable.

---
subcategory: "Sysdig Secure"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_connection"
description: |-
  Provides Sysdig Secure connection details.
---

# Data Source: sysdig_secure_connection

Provides information about current secure connection details.

## Example Usage

```terraform
data "sysdig_secure_connection" "current" {
}
```

## Attributes Reference

The following attributes are exported:

* `secure_url` - Returns `sysdig_secure_url` provider configuration attribute
* `secure_api_token` - Returns `sysdig_secure_api_token` provider configuration sensitive attribute
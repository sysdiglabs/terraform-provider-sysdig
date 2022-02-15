---
subcategory: "Sysdig Platform"
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_connection"
description: |-
  Provides secure connection details.
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

* `secure_url` - Sysdig Secure Endpoint URL basepath.
* `secure_api_token` - Sysdig Api Token for authentication (Sensitive).
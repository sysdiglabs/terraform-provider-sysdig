---
layout: "sysdig"
page_title: "Sysdig: sysdig_secure_user_rules_file"
sidebar_current: "docs-sysdig-secure-user-rules-file"
description: |-
  Manages the User Rules File.
---

# sysdig\_secure\_user\_rules\_file

Manages custom Falco rules availables for policies in Sysdig Secure.

With this resource you can upload a file to Rules Editor:

![Rules Editor][/docs/assets/rules-editor.jpg]

## Example usage

```hcl
resource "sysdig_secure_user_rules_file" "this" {
  content = "${file("${path.module}/rules-traefik.yaml")}"
}
```

## Argument Reference

* `content` - (Required) The custom Falco rules which will be available for policies in Secure

Note: You must concatenate all rules in one file and upload that file.

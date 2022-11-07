---
layout: "sysdig"
page_title: "Provider: Sysdig"
description: |-
  The Sysdig provider is used to interact with Sysdig products. The provider needs to be configured with proper API token before it can be used.
---

# Sysdig Provider

The Sysdig provider is used to interact with
[Sysdig Secure](https://sysdig.com/product/secure/) and
[Sysdig Monitor](https://sysdig.com/product/monitor/) products.

For Sysdig provider authentication **one of Monitor or Secure authentication is required**, being the other
optional.

For either options, the corresponding **URL** and **API Token** must be configured.
See options bellow.

Use the navigation to the left to read about the available resources.

## Example Usage

### Monitor Authentication

```terraform
provider "sysdig" {
  sysdig_monitor_url = "https://app.sysdigcloud.com"
  sysdig_monitor_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"  
}
```

### Secure Authentication

```terraform
provider "sysdig" {
  sysdig_secure_url="https://secure.sysdig.com"
  sysdig_secure_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

###  Both Secure and Monitor Provider Authentication

```terraform
provider "sysdig" {
  
  sysdig_secure_url="https://secure.sysdig.com"
  sysdig_secure_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  
  sysdig_monitor_url = "https://app.sysdigcloud.com"
  sysdig_monitor_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  
  extra_headers = {
    "Proxy-Authorization": "Basic xxxxxxxxxxxxxxxx"
  }
}


// create a new secure policy
resource "sysdig_secure_policy" "unexpected_inbound_tcp_connection_traefik" {
# ...
}
```

## Configuration Reference

For Sysdig provider authentication **one of Monitor or Secure authentication is required**, being the other
optional.


### Monitor authentication

When Monitor resources are to be created, this authentication must be in place.

* `sysdig_monitor_url` - (Required) This is the target Sysdig Monitor base API
  endpoint. It's intended to be used with OnPrem installations. 
  <br/>By default, it  points to `https://app.sysdigcloud.com`. [Find your Sysdig Saas region url](https://docs.sysdig.com/en/docs/administration/saas-regions-and-ip-ranges/#saas-regions-and-ip-ranges) (Sysdig Monitor endpoint)</br>
  It can also be sourced from the `SYSDIG_MONITOR_URL` environment variable.<br/>Notice: it should not be ended with a 
  slash.<br/><br/>

* `sysdig_monitor_api_token` - (Required) The Sysdig Monitor API token.
   <br/>[Find API Token](https://docs.sysdig.com/en/docs/administration/on-premises-deployments/find-the-super-admin-credentials-and-api-token/#find-sysdig-api-token)
   <br/>It can also be configured from the `SYSDIG_MONITOR_API_TOKEN` environment variable.
   <br/>Required if any `sysdig_monitor_*` resource or data source is used.<br/><br/>

* `sysdig_monitor_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the Monitor API. It can be useful for
  on-prem installations.<br/> It can also be sourced from the `SYSDIG_MONITOR_INSECURE_TLS`
  environment variable. By default, this is false.


###  Secure Authentication

When Secure resources are to be created, this authentication must be in place.

* `sysdig_secure_url` - (Required) This is the target Sysdig Secure base API
  endpoint. It's intended to be used with OnPrem installations.
  <br/>By default, it  points to `https://secure.sysdig.com`. 
  <br/>[Find your Sysdig Saas region url](https://docs.sysdig.com/en/docs/administration/saas-regions-and-ip-ranges/#saas-regions-and-ip-ranges) (Secure endpoint)
  <br/> It can also be sourced from the `SYSDIG_SECURE_URL` environment variable.
  <br/>Notice: it should not be ended with a slash.<br/><br/>

* `sysdig_secure_api_token` - (Required) The Sysdig Secure API token
  <br/>[Find API Token](https://docs.sysdig.com/en/docs/administration/on-premises-deployments/find-the-super-admin-credentials-and-api-token/#find-sysdig-api-token)
  <br/>It can also be configured from the `SYSDIG_SECURE_API_TOKEN` environment variable.
  <br/>Required if any `sysdig_secure_*` resource or data source is used.<br/><br/>

* `sysdig_secure_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the Secure API. It can be useful for 
  on-prem installations. It can also be sourced from the `SYSDIG_SECURE_INSECURE_TLS`
  environment variable. By default, this is false.<br/><br/>

###  Others
* `extra_headers` - (Optional) Defines extra HTTP headers that will be added to the client
  while performing HTTP API calls.

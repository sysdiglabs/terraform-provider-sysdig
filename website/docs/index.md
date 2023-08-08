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
See options below.

Sysdig provider can also be used to interact with [IBM Cloud Monitoring](https://cloud.ibm.com/docs/monitoring?topic=monitoring-getting-started). For more details check `IBM Cloud Monitoring Authentication` example and configuration reference below.

Use the navigation to the left to read about the available resources.

## Example Usage

Include a `providers` block to declare the requirement of the `sysdig` provider:

```terraform
terraform {
  required_providers {
    sysdig = {
      source = "sysdiglabs/sysdig"
      version = ">=0.5"
    }
  }
}
```

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

### IBM Cloud Monitoring Authentication

```terraform
provider "sysdig" {
  sysdig_monitor_team_id = 1234
  sysdig_monitor_team_name = "My team" # or use this as alternative to `sysdig_monitor_team_id`
  sysdig_monitor_url = "https://us-south.monitoring.cloud.ibm.com"
  ibm_monitor_iam_url = "https://iam.cloud.ibm.com"
  ibm_monitor_instance_id = "xxxxxx"
  ibm_monitor_api_key = "xxxxxx"
}
```

### IBM Workload Protection Authentication

```terraform
provider "sysdig" {
  sysdig_secure_team_id = 1234
  sysdig_secure_team_name = "My team" # or use this as alternative to `sysdig_monitor_team_id`
  sysdig_secure_url = "https://us-south.monitoring.cloud.ibm.com"
  ibm_secure_iam_url = "https://iam.cloud.ibm.com"
  ibm_secure_instance_id = "xxxxxx"
  ibm_secure_api_key = "xxxxxx"
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


### IBM Cloud Monitoring Authentication

When IBM Cloud Monitoring resources are to be created, this authentication must be in place.

* `sysdig_monitor_url` - (Required) This is the target IBM Cloud Monitoring API
  endpoint. It can also be sourced from the `SYSDIG_MONITOR_URL` environment variable. [Find your IBM Cloud Monitoring region url](https://cloud.ibm.com/docs/monitoring?topic=monitoring-endpoints#endpoints_monitoring).
  <br/>Notice: it should not be ended with a slash.<br/><br/>
* `ibm_monitor_iam_url` - (Required) This is the target IAM endpoint used to issue IBM IAM token by consuming `ibm_monitor_api_key`.
  Provider will handle token expiration and refresh it when needed.
  <br/>It can also be configured from the `SYSDIG_IBM_MONITOR_IAM_URL` environment variable.<br/><br/>
* `ibm_monitor_instance_id` (Required) This is the target instance ID (GUID format) of IBM instance which is hosting IBM Cloud Monitoring.
  <br/>It can also be configured from the `SYSDIG_IBM_MONITOR_INSTANCE_ID` environment variable.
  <br/><br/>
* `ibm_monitor_api_key` (Required) An API key is a unique code that is passed to an IBM IAM service to generate IAM token used for making HTTP request against IBM endpoints.
  This argument can be used to specify any kind of IBM API keys (User API key, Service ID, ...).
  <br/>It can also be configured from the `SYSDIG_IBM_MONITOR_API_KEY` environment variable.
  <br/><br/>
* `sysdig_monitor_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the IBM Monitoring Cloud API.
  <br/> It can also be sourced from the `SYSDIG_MONITOR_INSECURE_TLS`
  environment variable. By default, this is false.<br/><br/>
* `sysdig_monitor_team_id` - (Optional) Use this argument to specify team in which you will be logged in.
  If not specified, default team will be used. This argument has precedence over `sysdig_monitor_team_name` if both are specified.<br/>
  It can also be configured from the `SYSDIG_MONITOR_TEAM_ID` environment variable.<br/><br/>
* `sysdig_monitor_team_name` - (Optional) This argument is the alternative way of specifying team in which you will be logged in.
  It has exactly the same meaning as `sysdig_monitor_team_id`, but instead of specifying team ID you are specifying a team name.</br>
  It can also be configured from the `SYSDIG_MONITOR_TEAM_NAME` environment variable.<br/><br/>

### IBM Workload Protection Authentication

When IBM Workload Protection resources are to be created, this authentication must be in place.

* `sysdig_secure_url` - (Required) This is the target IBM Workload Protection API
  endpoint. It can also be sourced from the `SYSDIG_SECURE_URL` environment variable. [Find your Workload Protection region url](https://cloud.ibm.com/docs/workload-protection?topic=workload-protection-endpoints#endpoints_monitoring).
  <br/>Notice: it should not be ended with a slash.<br/><br/>
* `ibm_secure_iam_url` - (Required) This is the target IAM endpoint used to issue IBM IAM token by consuming `ibm_secure_api_key`.
  Provider will handle token expiration and refresh it when needed.
  <br/>It can also be configured from the `SYSDIG_IBM_SECURE_IAM_URL` environment variable.<br/><br/>
* `ibm_secure_instance_id` (Required) This is the target instance ID (GUID format) of IBM instance which is hosting IBM Workload Protection.
  <br/>It can also be configured from the `SYSDIG_IBM_SECURE_INSTANCE_ID` environment variable.
  <br/><br/>
* `ibm_secure_api_key` (Required) An API key is a unique code that is passed to an IBM IAM service to generate IAM token used for making HTTP request against IBM endpoints.
  This argument can be used to specify any kind of IBM API keys (User API key, Service ID, ...).
  <br/>It can also be configured from the `SYSDIG_IBM_SECURE_API_KEY` environment variable.
  <br/><br/>
* `sysdig_secure_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the IBM Workload Protection API.
  <br/> It can also be sourced from the `SYSDIG_SECURE_INSECURE_TLS`
  environment variable. By default, this is false.<br/><br/>
* `sysdig_secure_team_id` - (Optional) Use this argument to specify team in which you will be logged in.
  If not specified, default team will be used. This argument has precedence over `sysdig_secure_team_name` if both are specified.<br/>
  It can also be configured from the `SYSDIG_SECURE_TEAM_ID` environment variable.<br/><br/>
* `sysdig_secure_team_name` - (Optional) This argument is the alternative way of specifying team in which you will be logged in.
  It has exactly the same meaning as `sysdig_secure_team_id`, but instead of specifying team ID you are specifying a team name.</br>
  It can also be configured from the `SYSDIG_SECURE_TEAM_NAME` environment variable.<br/><br/>

> **Note**
> Enabling resources and data sources on IBM is under active development.
>
> For now, you can manage following resources:
> - `sysdig_monitor_team`
> - `sysdig_secure_team`
> - `sysdig_monitor_notification_channel_email`
> - `sysdig_secure_notification_channel_email`
> - `sysdig_monitor_notification_channel_opsgenie`
> - `sysdig_secure_notification_channel_opsgenie`
> - `sysdig_monitor_notification_channel_pagerduty`
> - `sysdig_secure_notification_channel_pagerduty`
> - `sysdig_monitor_notification_channel_slack`
> - `sysdig_secure_notification_channel_slack`
> - `sysdig_monitor_notification_channel_sns`
> - `sysdig_secure_notification_channel_sns`
> - `sysdig_monitor_notification_channel_victorops`
> - `sysdig_secure_notification_channel_victorops`
> - `sysdig_monitor_notification_channel_webhook`
> - `sysdig_secure_notification_channel_webhook`
> - `sysdig_monitor_alert_downtime`
> - `sysdig_monitor_alert_event`
> - `sysdig_monitor_alert_metric`
> - `sysdig_monitor_alert_promql`
> - `sysdig_monitor_alert_anomaly`
> - `sysdig_monitor_alert_group_outlier`
> - `sysdig_monitor_alert_v2_downtime`
> - `sysdig_monitor_alert_v2_event`
> - `sysdig_monitor_alert_v2_metric`
> - `sysdig_monitor_alert_v2_prometheus`
> - `sysdig_monitor_dashboard`
> - `sysdig_secure_posture_zone`
>
> And data sources:
> - `sysdig_monitor_notification_channel_pagerduty`
> - `sysdig_monitor_notification_channel_email`
> - `sysdig_current_user`
> - `sysdig_secure_notification_channel`
> - `sysdig_secure_posture_policies`

###  Others
* `extra_headers` - (Optional) Defines extra HTTP headers that will be added to the client
  while performing HTTP API calls.

## Troubleshooting

If you get a:

```
panic: Invalid diagnostic: empty summary. This is always a bug in the provider implementation
```

Please check:

If you are using a monitor resource those variables should be correctly set:
```
sysdig_monitor_url = "https://app.sysdigcloud.com"
sysdig_monitor_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```
If you are using a secure resource those variables should be correctly set:
```
sysdig_secure_url="https://secure.sysdig.com"
sysdig_secure_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```
Ensure for both url variables your region is correctly set.
For more info on regions [check here](https://docs.sysdig.com/en/docs/administration/saas-regions-and-ip-ranges/).

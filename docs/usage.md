
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
  sysdig_monitor_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  sysdig_secure_api_token = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}

// Create a new secure policy
resource "sysdig_secure_policy" "unexpected_inbound_tcp_connection_traefik" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `sysdig_monitor_api_token` - (Optional) The Sysdig Secure API token, it must be
  present, but you can get it from the `SYSDIG_MONITOR_API_TOKEN` environment variable.
  Required if any `sysdig_monitor_*` resource or data source is used. 

* `sysdig_secure_api_token` - (Optional) The Sysdig Secure API token, it must be
  present, but you can get it from the `SYSDIG_SECURE_API_TOKEN` environment variable.
  Required if any `sysdig_secure_*` resource or data source is used.

* `sysdig_monitor_url` - (Optional) This is the target Sysdig Secure base API
  endpoint. It's intended to be used with OnPrem installations. By defaults it
  points to `https://app.sysdigcloud.com`, and notice that should not be ended
  with an slash. It can also be sourced from the `SYSDIG_MONITOR_URL` environment
  variable.
  
* `sysdig_secure_url` - (Optional) This is the target Sysdig Secure base API
  endpoint. It's intended to be used with OnPrem installations. By defaults it
  points to `https://secure.sysdig.com`, and notice that should not be ended
  with an slash. It can also be sourced from the `SYSDIG_SECURE_URL` environment
  variable.
  
* `sysdig_monitor_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the Monitor API. It can be useful for 
  on-prem installations. It can also be sourced from the `SYSDIG_MONITOR_INSECURE_TLS`
  environment variable. By default this is false.

* `sysdig_secure_insecure_tls` - (Optional) Defines if the HTTP client can ignore
  the use of invalid HTTPS certificates in the Secure API. It can be useful for 
  on-prem installations. It can also be sourced from the `SYSDIG_SECURE_INSECURE_TLS`
  environment variable. By default this is false.

# Data Sources


## sysdig\_secure\_notification_channel

Retrieves the information of an existing Sysdig Secure Notification Channel.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
data "sysdig_secure_notification_channel" "sample-email" {
  name                 = "Example Channel - Email"
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel.

### Attributes Reference

* `enabled` - If false, the channel will not emit notifications.

* `type` - Will be one of the following:  "EMAIL", "SNS", "OPSGENIE", 
    "VICTOROPS", "WEBHOOK", "SLACK", "PAGER_DUTY".

* `notify_when_ok` - Send a new notification when the alert condition is 
    no longer triggered.

* `notify_when_resolved` - Send a new notification when the alert is manually 
    acknowledged by a user.

* `send_test_notification` - Send an initial test notification to check
    if the notification channel is working.

#### Attributes for type EMAIL

* `recipients` - Comma-separated list of recipients that will receive 
    the message.
    
#### Attributes for type Amazon SNS

* `topics` - List of ARNs from the SNS topics.

#### Attributes for type VICTOROPS

* `api_key` - Key for the API.

* `routing_key` - Routing key for VictorOps. 

#### Attributes for type OPSGENIE

* `api_key` - Key for the API.

#### Attributes for type WEBHOOK

* `url` - URL to send the event.

#### Attributes for type SLACK

* `url` - URL of the Slack.

* `channel` - Channel name from this Slack.

#### Attributes for type PAGERDUTY

* `account` - Pagerduty account.

* `service_key` - Service Key for the Pagerduty account.

* `service_name` - Service name for the Pagerduty account.


---



# Resources


## sysdig\_monitor\_alert\_anomaly

Creates a Sysdig Monitor Anomaly Alert. Monitor hosts based on their historical behaviors and alert when they deviate.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_alert_anomaly" "sample" {
	name = "[Kubernetes] Anomaly Detection Alert"
	description = "Detects an anomaly in the cluster"
	severity = 6

	monitor = ["cpu.used.percent", "memory.bytes.used"]

    trigger_after_minutes = 10

	multiple_alerts_by = ["kubernetes.cluster.name", 
                          "kubernetes.namespace.name", 
                          "kubernetes.deployment.name", 
                          "kubernetes.pod.name"]
}
```

### Argument Reference

#### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. 
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
 

##### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

#### Metric alert arguments

* `monitor` - (Required) Array of metrics to monitor and alert on. Example: `["cpu.used.percent", "cpu.cores.used", "memory.bytes.used", "fs.used.percent", "thread.count", "net.request.count.in"]`.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

### Attributes Reference

#### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.


---


## sysdig\_monitor\_alert\_downtime

Creates a Sysdig Monitor Downtime Alert. Monitor any type of entity - host, container, process, service, etc - and alert when the entity goes down.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_alert_downtime" "sample" {
	name = "[Kubernetes] Downtime Alert"
	description = "Detects a downtime in the Kubernetes cluster"
	severity = 2

	entities_to_monitor = ["kubernetes.namespace.name"]
	
	trigger_after_minutes = 10
	trigger_after_pct = 100
}
```

### Argument Reference

#### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. 
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
 

##### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

#### Metric alert arguments

* `entities_to_monitor` - (Required) List of metrics to monitor downtime and alert on. Example: `["kubernetes.namespace.name"]` to detect namespace removal or `["host.hostName"]` to detect host downtime.
* `trigger_after_pct` - (Optional) Below of this percentage of downtime the alert will be triggered. Defaults to 100.  

### Attributes Reference

#### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

---


## sysdig\_monitor\_alert\_event

Creates a Sysdig Monitor Event Alert. Monitor occurrences of specific events, and alert if the total 
number of occurrences violates a threshold. Useful for alerting on container, orchestration, and 
service events like restarts and deployments.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_alert_event" "sample" {
	name = "[Kubernetes] Failed to pull image"
	description = "A Kubernetes pod failed to pull an image from the registry"
	severity = 4

	event_name = "Failed to pull image"
	source = "kubernetes"
	event_rel = ">"
	event_count = 0

	multiple_alerts_by = ["kubernetes.pod.name"]
	
	trigger_after_minutes = 1
}
```

### Argument Reference

#### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. 
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
 

##### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

#### Metric alert arguments

* `event_name` - (Required) String that matches part of name, tag or the description of Sysdig Events.
* `source` - (Required) Source of the event. It can be `docker` or `kubernetes`. 
* `event_rel` - (Required) Relationship of the event count. It can be `>`, `>=`, `<`, `<=`, `=` or `!=`.
* `event_count` - (Required) Number of events to match with event_rel.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

### Attributes Reference

#### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

---


## sysdig\_monitor\_alert\_group\_outlier

Creates a Sysdig Monitor Group Outlier Alert. Monitor a group of hosts and be notified when one acts differently from the rest.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_alert_group_outlier" "sample" {
	name = "[Kubernetes] A node is using more CPU than the rest"
	description = "Monitors the cluster and checks when a node has more CPU usage than the others"
	severity = 6

	monitor = ["cpu.used.percent"]
	
	trigger_after_minutes = 10

	capture {
		filename = "TERRAFORM_TEST"
		duration = 15
	}
}
```

### Argument Reference

#### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. 
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
 
 
##### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

#### Metric alert arguments

* `monitor` - (Required) Array of metrics to monitor and alert on. Example: `["cpu.used.percent", "cpu.cores.used", "memory.bytes.used", "fs.used.percent", "thread.count", "net.request.count.in"]`.  

### Attributes Reference

#### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

---


## sysdig\_monitor\_alert\_metric

Creates a Sysdig Monitor Metric Alert. Monitor time-series metrics and alert if they violate user-defined thresholds.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_alert_metric" "sample" {
	name = "[Kubernetes] CrashLoopBackOff"
	description = "A Kubernetes pod failed to restart"
	severity = 6

	metric = "sum(timeAvg(kubernetes.pod.restart.count)) > 2"
	trigger_after_minutes = 1

	multiple_alerts_by = ["kubernetes.cluster.name",
                          "kubernetes.namespace.name",
                          "kubernetes.deployment.name",
                          "kubernetes.pod.name"]

	capture {
		filename = "CrashLoopBackOff"
		duration = 15
	}
}
```

### Argument Reference

#### Common alert arguments

These arguments are common to all alerts in Sysdig Monitor.

* `name` - (Required) The name of the Monitor alert. It must be unique.
* `description` - (Optional) The description of Monitor alert.
* `severity` - (Optional) Severity of the Monitor alert. It must be a value between 0 and 7,
               with 0 being the most critical and 7 the less critical. Defaults to 4.
* `trigger_after_minutes` - (Required) Threshold of time for the status to stabilize until the alert is fired.
* `scope` - (Optional) Part of the infrastructure where the alert is valid. Defaults to the entire infrastructure. 
* `enabled` - (Optional) Boolean that defines if the alert is enabled or not. Defaults to true.
* `notification_channels` - (Optional) List of notification channel IDs where an alert must be sent to once fired.
* `renotification_minutes` - (Optional) Number of minutes for the alert to re-notify until the status is solved.
 

##### Capture

Enables the creation of a capture file of the syscalls during the event.

* `filename` - (Required) Defines the name of the capture file.
* `duration` - (Required) Time frame in seconds of the capture.
* `filter` - (Optional) Additional filter to apply to the capture. For example: `proc.name contains nginx`.

#### Metric alert arguments

* `metric` - (Required) Metric to monitor and alert on. Example: `sum(timeAvg(kubernetes.pod.restart.count)) > 2` or `avg(avg(cpu.used.percent)) > 50`.
* `multiple_alerts_by` - (Optional) List of segments to trigger a separate alert on. Example: `["kubernetes.cluster.name", "kubernetes.namespace.name"]`.  

### Attributes Reference

#### Common alert attributes

In addition to all arguments above, the following attributes are exported, which are common to all the
alerts in Sysdig Monitor:

* `version` - Current version of the resource in Sysdig Monitor.
* `team` - Team ID that owns the alert.

---


## sysdig\_monitor\_notification_channel\_email

Creates a Sysdig Monitor Notification Channel of type Email.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_email" "sample_email" {
	name                    = "Example Channel - Email"
	recipients              = ["foo@localhost.com", "bar@localhost.com"]
	enabled                 = true
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `recipients` - (Required) List of recipients that will receive 
    the message.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_opsgenie

Creates a Sysdig Monitor Notification Channel of type OpsGenie.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_opsgenie" "sample-opsgenie" {
	name                    = "Example Channel - OpsGenie"
	enabled                 = true
	api_key                 = "2349324-342354353-5324-23"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_pagerduty

Creates a Sysdig Monitor Notification Channel of type Pagerduty.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_pagerduty" "sample-pagerduty" {
	name                    = "Example Channel - Pagerduty"
	enabled                 = true
	account                 = "account"
	service_key             = "XXXXXXXXXX"
	service_name            = "sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_slack

Creates a Sysdig Monitor Notification Channel of type Slack.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_slack" "sample-slack" {
	name                    = "Example Channel - Slack"
	enabled                 = true
	url                     = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel                 = "#sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL of the Slack.

* `channel` - (Required) Channel name from this Slack.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_sns

Creates a Sysdig Monitor Notification Channel of type Amazon SNS.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_sns" "sample-amazon-sns" {
	name                    = "Example Channel - Amazon SNS"
	enabled                 = true
	topics                  = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `topics` - (Required) List of ARNs from the SNS topics.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_victorops

Creates a Sysdig Monitor Notification Channel of type VictorOps.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_victorops" "sample-victorops" {
	name                    = "Example Channel - VictorOps"
	enabled                 = true
	api_key                 = "1234342-4234243-4234-2"
	routing_key             = "My team"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_monitor\_notification\_channel\_webhook

Creates a Sysdig Monitor Notification Channel of type Webhook.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_monitor_notification_channel_webhook" "sample-webhook" {
	name                    = "Example Channel - Webhook"
	enabled                 = true
	url                     = "localhost:8080"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL to send the event.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_list

Creates a Sysdig Secure Falco List.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_list" "allowed_dev_files" {
  name = "allowed_dev_files"
  items = ["/dev/null", "/dev/stdin", "/dev/stdout", "/dev/stderr", "/dev/random", 
           "/dev/urandom", "/dev/console", "/dev/kmsg"]
  append = true # default: false
}
```

### Argument Reference

* `name` - (Required) The name of the Secure list. It must be unique if it's not in append mode.

* `items` - (Required) Elements in the list. Elements can be another lists.

* `append` - (Optional)  Adds these elements to an existing list. Used to extend existing lists provided by Sysdig.
    The rules can only be extended once, for example if there is an existing list called "foo", one can have another 
    append rule called "foo" but not a second one. By default this is false.


---


## sysdig\_secure\_macro

Creates a Sysdig Secure Falco Macro.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_macro" "http_port" {
  name = "web_port"
  condition = "fd.sport=80"
}

resource "sysdig_secure_macro" "https_port" {
  name = "web_port"
  condition = "or fd.sport=443"
  append = true # default: false
}
```

### Argument Reference

* `name` - (Required) The name of the macro. It must be unique if it's not in append mode.

* `condition` - (Required) Macro condition. It can contain lists or other macros.

* `append` - (Optional)  Adds these elements to an existing macro. Used to extend existing macros provided by Sysdig.
    The macros can only be extended once, for example if there is an existing macro called "foo", one can have another 
    append macro called "foo" but not a second one. By default this is false.


---


## sysdig\_secure\_notification_channel

Creates a Sysdig Secure Notification Channel.

!> **Warning:** This resource is deprecated and will be removed. Please use the different `sysdig_secure_notification_channel_*` resources.

### Example usage

```hcl
resource "sysdig_secure_notification_channel" "sample-email" {
  name                 = "Example Channel - Email"
  enabled              = true
  type                 = "EMAIL"
  recipients           = "root@localhost.com"
  notify_when_ok       = false
  notify_when_resolved = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `enabled` - (Required) If false, the channel will not emit notifications.

* `type` - (Required) Must be one of the following:  "EMAIL", "SNS", "OPSGENIE", 
    "VICTOROPS", "WEBHOOK", "SLACK", "PAGER_DUTY".

* `notify_when_ok` - (Required) Send a new notification when the alert condition is 
    no longer triggered.

* `notify_when_resolved` - (Required) Send a new notification when the alert is manually 
    acknowledged by a user.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working.

#### Arguments for type EMAIL

* `recipients` - (Required) Comma-separated list of recipients that will receive 
    the message.
    
#### Arguments for type Amazon SNS

* `topics` - (Required) List of ARNs from the SNS topics.

#### Arguments for type VICTOROPS

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

#### Arguments for type OPSGENIE

* `api_key` - (Required) Key for the API.

#### Arguments for type WEBHOOK

* `url` - (Required) URL to send the event.

#### Arguments for type SLACK

* `url` - (Required) URL of the Slack.

* `channel` - (Required) Channel name from this Slack.

#### Arguments for type PAGERDUTY

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.


---


## sysdig\_secure\_notification_channel\_email

Creates a Sysdig Secure Notification Channel of type Email.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_email" "sample_email" {
	name                    = "Example Channel - Email"
	recipients              = ["foo@localhost.com", "bar@localhost.com"]
	enabled                 = true
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `recipients` - (Required) List of recipients that will receive 
    the message.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_opsgenie

Creates a Sysdig Secure Notification Channel of type OpsGenie.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_opsgenie" "sample-opsgenie" {
	name                    = "Example Channel - OpsGenie"
	enabled                 = true
	api_key                 = "2349324-342354353-5324-23"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_pagerduty

Creates a Sysdig Secure Notification Channel of type Pagerduty.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_pagerduty" "sample-pagerduty" {
	name                    = "Example Channel - Pagerduty"
	enabled                 = true
	account                 = "account"
	service_key             = "XXXXXXXXXX"
	service_name            = "sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `account` - (Required) Pagerduty account.

* `service_key` - (Required) Service Key for the Pagerduty account.

* `service_name` - (Required) Service name for the Pagerduty account.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_slack

Creates a Sysdig Secure Notification Channel of type Slack.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_slack" "sample-slack" {
	name                    = "Example Channel - Slack"
	enabled                 = true
	url                     = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
	channel                 = "#sysdig"
	notify_when_ok          = false
	notify_when_resolved    = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL of the Slack.

* `channel` - (Required) Channel name from this Slack.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_sns

Creates a Sysdig Secure Notification Channel of type Amazon SNS.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_sns" "sample-amazon-sns" {
	name                    = "Example Channel - Amazon SNS"
	enabled                 = true
	topics                  = ["arn:aws:sns:us-east-1:273489009834:my-alerts2", "arn:aws:sns:us-east-1:279948934544:my-alerts"]
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `topics` - (Required) List of ARNs from the SNS topics.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_victorops

Creates a Sysdig Secure Notification Channel of type VictorOps.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_victorops" "sample-victorops" {
	name                    = "Example Channel - VictorOps"
	enabled                 = true
	api_key                 = "1234342-4234243-4234-2"
	routing_key             = "My team"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `api_key` - (Required) Key for the API.

* `routing_key` - (Required) Routing key for VictorOps. 

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_notification\_channel\_webhook

Creates a Sysdig Secure Notification Channel of type Webhook.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_notification_channel_webhook" "sample-webhook" {
	name                    = "Example Channel - Webhook"
	enabled                 = true
	url                     = "localhost:8080"
	notify_when_ok          = false
	notify_when_resolved    = false
	send_test_notification  = false
}
```

### Argument Reference

* `name` - (Required) The name of the Notification Channel. Must be unique.

* `url` - (Required) URL to send the event.

* `enabled` - (Optional) If false, the channel will not emit notifications. Default is true.

* `notify_when_ok` - (Optional) Send a new notification when the alert condition is 
    no longer triggered. Default is false.

* `notify_when_resolved` - (Optional) Send a new notification when the alert is manually 
    acknowledged by a user. Default is false.

* `send_test_notification` - (Optional) Send an initial test notification to check
    if the notification channel is working. Default is false.


---


## sysdig\_secure\_policy

Creates a Sysdig Secure Policy.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_policy" "write_apt_database" {
  name = "Write apt database"
  description = "an attempt to write to the dpkg database by any non-dpkg related program"
  severity = 4
  enabled = true

  // Scope selection
  scope = "container.id != \"\""

  // Rule selection
  rule_names = ["Terminal shell in container"]

  actions {
    container = "stop"
    capture {
      seconds_before_event = 5
      seconds_after_event = 10
    }
  }

  notification_channels = [10000]

}
```

### Argument Reference

* `name` - (Required) The name of the Secure policy. It must be unique.

* `description` - (Required) The description of Secure policy.

* `severity` - (Optional) The severity of Secure policy. The accepted values
    are: 0 (High), 4 (Medium), 6 (Low) and 7 (Info). The default value is 4 (Medium).

* `enabled` - (Optional) Will secure process with this rule?. By default this is true.

- - -

#### Scope selection

* `scope` - (Optional) Limit appplication scope based in one expresion. For
    example: "host.ip.private = \\"10.0.23.1\\"". By default the rule won't be scoped
    and will target the entire infrastructure.

- - -

#### Actions block

The actions block is optional and supports:

* `container` - (Optional) The action applied to container when this Policy is
    triggered. Can be *stop* or *pause*.

* `capture` - (Optional) Captures with Sysdig the stream of system calls:
    * `seconds_before_event` - (Required) Captures the system calls during the
    amount of seconds before the policy was triggered.
    * `seconds_after_event` - (Required) Captures the system calls for the amount
    of seconds after the policy was triggered.

- - -

#### Falco rule selection

* `rule_names` - (Optional) Array with the name of the rules to match.

- - -

#### Notification

* `notification_channels` - (Optional) IDs of the notification channels to send alerts to
    when the policy is fired.


---


## sysdig\_secure\_rule\_container

Creates a Sysdig Secure Container Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_rule_container" "sample" {
  name = "Nginx container spawned"
  description = "A container withthe nginx image spawned in the cluster."
  tags = ["container", "cis"]

  matching = true // default
  containers = ["nginx"]
}
```

### Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

#### Matching

* `matching` - (Optional) Defines if the image name matches or not with the provided list. Default is true.
* `containers` - (Required) List of containers to match.

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_secure\_rule\_falco

Creates a Sysdig Secure Falco Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_rule_falco" "example" {
  name = "Terminal shell in container" // ID
  description = "A shell was used as the entrypoint/exec point into a container with an attached terminal."
  tags = ["container", "shell", "mitre_execution"]

  condition = "spawned_process and container and shell_procs and proc.tty != 0 and container_entrypoint"
  output = "A shell was spawned in a container with an attached terminal (user=%user.name %container.info shell=%proc.name parent=%proc.pname cmdline=%proc.cmdline terminal=%proc.tty container_id=%container.id image=%container.image.repository)"
  priority = "notice"
  source = "syscall" // syscall or k8s_audit
}

```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

- - -

#### Conditions

* `condition` - (Required) A [Falco condition](https://falco.org/docs/rules/) is simply a Boolean predicate on Sysdig events expressed using the Sysdig [filter syntax](http://www.sysdig.org/wiki/sysdig-user-guide/#filtering) and macro terms. 
* `output` - (Required) Add additional information to each Falco notification's output.
* `priority` - (Required) The priority of the Falco rule. It can be: "emergency", "alert", "critical", "error", "warning", "notice", "informational", "informational" or "debug".
* `source` - (Required) The source of the event. It can be either "syscall" or "k8s_audit".

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_secure\_rule\_filesystem

Creates a Sysdig Secure Filesystem Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl

resource "sysdig_secure_rule_filesystem"  "example" {
  name = "Apache writing to non allowed directory"
  description = "Attempt to write to directories that should be immutable"
  tags = ["filesystem", "cis"]

  read_only {
    matching = true // default
    paths = ["/etc"]
  }

  read_write {
    matching = true // default
    paths = ["/var/log/apache2", "/dev/tty"]
  }
}
```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

#### Read Only

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

#### Read Write

* `matching` - (Optional) Defines if the path matches or not with the provided list. Default is true.
* `paths` - (Required) List of paths to match.

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_secure\_rule\_network

Creates a Sysdig Secure Network Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_rule_network" "example" {
  name = "Disallowed SSH Connection"
  description = "Detect any new ssh connection to a host other than those in an allowed group of hosts"
  tags = ["network", "mitre_remote_service"]

  block_inbound = true
  block_outbound = true

  tcp {
    matching = true // default
    ports = [22]
  }

  udp {
    matching = true // default
    ports = [22]
  }
}

```

### Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

#### Disallow incoming or outgoing connections

* `block_inbound` - (Required) Detect if there is an inbound connection.
* `block_outbound` - (Required) Detect if there is an outbound connection.

#### Detect TCP Connections

* `matching` - (Optional) Defines if the port matches or not with the provided list. Default is true.
* `ports` - (Required) List of ports to match.

#### Detect UDP Connections

* `matching` - (Optional) Defines if the port matches or not with the provided list. Default is true.
* `ports` - (Required) List of ports to match.

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_secure\_rule\_process

Creates a Sysdig Secure Process Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_rule_process" "sample" {
  name = "Launch Suspicious Network Tool in Container" // ID
  description = "Detect network tools launched inside container"

  matching = true // default
  processes = ["nc", "ncat", "nmap", "dig", "tcpdump", "tshark", "ngrep"]
}

```

### Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

#### Matching

* `matching` - (Optional) Defines if the process name matches or not with the provided list. Default is true.
* `processes` - (Required) List of processes to match.

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_secure\_rule\_syscall

Creates a Sysdig Secure Syscall Rule.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_secure_rule_syscall" "foo" {
  name = "Unexpected mount syscall" // ID
  description = "Syscall 'mount' detected"

  matching = true // default
  syscalls = ["mount"]
}
```

### Argument Reference

* `name` - (Required) The name of the Secure rule. It must be unique.
* `description` - (Required) The description of Secure rule.
* `tags` - (Optional) A list of tags for this rule.

#### Matching

* `matching` - (Optional) Defines if the syscall name matches or not with the provided list. Default is true.
* `processes` - (Required) List of syscalls to match.

### Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - Current version of the resource in Sysdig Secure.

---


## sysdig\_user

Creates a user in Sysdig.

~> **Note:** This resource is still experimental, and is subject of being changed.

### Example usage

```hcl
resource "sysdig_user" "foo-bar" {
  email = "foo.bar@sysdig.com"
  system_role = "ROLE_CUSTOMER"
  first_name = "foo"
  last_name = "bar"
}
```

### Argument Reference

* `email` - (Required) The email for the user to invite.

* `system_role` - (Optional) The privileges for the user. It can be either "ROLE_USER" or "ROLE_CUSTOMER".
    If set to "ROLE_CUSTOMER", the user will be known as an admin.

* `first_name` - (Optional) The name of the user.

* `last_name` - (Optional) The last name of the user.

---


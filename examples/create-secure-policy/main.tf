provider "sysdig" { }

resource "sysdig_secure_user_rules_file" "this" {
  content = "${file("${path.module}/rules-traefik.yaml")}"
}

resource "sysdig_secure_notification_channel" "sample-email" {
  name = "Example Channel - Email"
  enabled = true
  type = "EMAIL"
  recipients = "root@localhost.com"
  notify_when_ok = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-amazon-sns" {
  name = "Example Channel - Amazon SNS"
  enabled = true
  type = "SNS"
  topics = "arn:aws:sns:us-east-1:273107874544:my-alerts,arn:aws:sns:us-east-1:273107874544:my-alerts2"
  notify_when_ok = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-victorops" {
  name = "Example Channel - VictorOps"
  enabled = true
  type = "VICTOROPS"
  api_key = "1234342-4234243-4234-2"
  routing_key = "My team"
  notify_when_ok = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-opsgenie" {
  name = "Example Channel - OpsGenie"
  enabled = true
  type = "OPSGENIE"
  api_key = "2349324-342354353-5324-23"
  notify_when_ok = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-webhook" {
  name = "Example Channel - Webhook"
  enabled = true
  type = "WEBHOOK"
  url = "localhost:8080"
  notify_when_ok = false
  notify_when_resolved = false
}

resource "sysdig_secure_policy" "sample" {
  name = "Write apt stuff"
  description = "an attempt to write to the dpkg database by any non-dpkg related program"
  severity = 4
  enabled = true

  // Scope selection
  //filter = "host.ip.private = \"10.0.23.1\""
  container_scope = true
  host_scope = true

  notification_channels = [
    "${sysdig_secure_notification_channel.sample-email.id}",
    "${sysdig_secure_notification_channel.sample-opsgenie.id}"]

  //actions {
  //  container = "pause"

  //  capture {
  //    seconds_before_event = 60
  //    seconds_after_event = 60
  //  }
  //}

  // Falco rule selection
  falco_rule_name_regex = "Unexpected spawned process traefik"
}

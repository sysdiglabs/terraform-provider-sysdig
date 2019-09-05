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

resource "sysdig_secure_notification_channel" "sample-slack" {
  name = "Example Channel - Slack"
  enabled = true
  type = "SLACK"
  url = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
  channel = "#sysdig"
  notify_when_ok = true
  notify_when_resolved = true
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


resource "sysdig_secure_policy" "sample2" {
  name = "Other example of Policy"
  description = "this is other example of policy"
  severity = 4
  enabled = true

  container_scope = true
  host_scope = true


  processes {
    default = "accept"
    whitelist = [
      "mysql",
      "apache"]
    blacklist = [
      "ssh"]
  }

  containers {
    default = "none"
    whitelist = [
      "cassandra"]
    blacklist = [
      "mongo"]
  }

  syscalls {
    default = "accept"
    whitelist = [
      "accept",
      "close"]
    blacklist = [
      "bind",
      "bpf"]
  }

  network {
    inbound = "accept"

    outbound = "deny"

    listening_ports {
      default = "none"
      tcp {
        whitelist = [
          80,
          443]
        blacklist = [
          8080,
          5000]
      }
      udp {
        whitelist = [
          53,
          4000]
        blacklist = [
          3400,
          543]
      }
    }
  }

  filesystem {
    read {
      whitelist = [
        "/home"]
      blacklist = [
        "/etc"]
    }
    readwrite {
      whitelist = [
        "/home"]
      blacklist = [
        "/tmp"]
    }
    other_paths = "none"
  }


  notification_channels = [
    "${sysdig_secure_notification_channel.sample-victorops.id}"]

  falco_rule_name_regex = "Unexpected spawned process traefik"
}

resource "sysdig_secure_policies_priority" "priority" {
  policies = [
    "${sysdig_secure_policy.sample2.id}",
    "${sysdig_secure_policy.sample.id}"]
}

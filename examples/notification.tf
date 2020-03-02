
resource "sysdig_secure_notification_channel" "sample-email" {
  name                 = "Example Channel - Email"
  enabled              = true
  type                 = "EMAIL"
  recipients           = "root@localhost.com"
  notify_when_ok       = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-amazon-sns" {
  name                 = "Example Channel - Amazon SNS"
  enabled              = true
  type                 = "SNS"
  topics               = "arn:aws:sns:us-east-1:273107874544:my-alerts,arn:aws:sns:us-east-1:273107874544:my-alerts2"
  notify_when_ok       = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-victorops" {
  name                 = "Example Channel - VictorOps"
  enabled              = true
  type                 = "VICTOROPS"
  api_key              = "1234342-4234243-4234-2"
  routing_key          = "My team"
  notify_when_ok       = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-opsgenie" {
  name                 = "Example Channel - OpsGenie"
  enabled              = true
  type                 = "OPSGENIE"
  api_key              = "2349324-342354353-5324-23"
  notify_when_ok       = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-webhook" {
  name                 = "Example Channel - Webhook"
  enabled              = true
  type                 = "WEBHOOK"
  url                  = "localhost:8080"
  notify_when_ok       = false
  notify_when_resolved = false
}

resource "sysdig_secure_notification_channel" "sample-slack" {
  name                 = "Example Channel - Slack"
  enabled              = true
  type                 = "SLACK"
  url                  = "https://hooks.slack.cwom/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
  channel              = "#sysdig"
  notify_when_ok       = true
  notify_when_resolved = true
}

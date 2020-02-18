
resource "sysdig_secure_policy" "sample" {
  name = "Other example of Policy"
  description = "this is other example of policy"
  enabled = true
  severity = 4
  scope = "container.id != \"\""
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
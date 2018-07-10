provider "sysdig" { }

resource "sysdig_secure_user_rules_file" "this" {
  content = "${file("${path.module}/rules-traefik.yaml")}"
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

  //actions {
  //  container = "pause"

  //  capture {
  //    seconds_before_event = 60
  //    seconds_after_event = 60
  //  }
  //}

  // Falco rule selection
  depends_on = ["sysdig_secure_user_rules_file.this"]
  falco_rule_name_regex = "Unexpected spawned process traefik"
}

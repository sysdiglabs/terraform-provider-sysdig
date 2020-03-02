
resource "sysdig_monitor_alert_anomaly" "sample" {
  name        = "[Kubernetes] Anomaly Detection Alert"
  description = "Detects an anomaly in the cluster"
  severity    = 6

  monitor = ["cpu.used.percent", "memory.bytes.used"]

  trigger_after_minutes = 10

  multiple_alerts_by = ["kubernetes.cluster.name",
    "kubernetes.namespace.name",
    "kubernetes.deployment.name",
  "kubernetes.pod.name"]
}

resource "sysdig_monitor_alert_downtime" "sample" {
  name        = "[Kubernetes] Downtime Alert"
  description = "Detects a downtime in the Kubernetes cluster"
  severity    = 2

  entities_to_monitor = ["kubernetes.namespace.name"]

  trigger_after_minutes = 10
  trigger_after_pct     = 100
}

resource "sysdig_monitor_alert_event" "sample" {
  name        = "[Kubernetes] Failed to pull image"
  description = "A Kubernetes pod failed to pull an image from the registry"
  severity    = 4

  event_name  = "Failed to pull image"
  source      = "kubernetes"
  event_rel   = ">"
  event_count = 0

  multiple_alerts_by = ["kubernetes.pod.name"]

  trigger_after_minutes = 1
}

resource "sysdig_monitor_alert_group_outlier" "sample" {
  name        = "[Kubernetes] A node is using more CPU than the rest"
  description = "Monitors the cluster and checks when a node has more CPU usage than the others"
  severity    = 6

  monitor = ["cpu.used.percent"]

  trigger_after_minutes = 10

  capture {
    filename = "TERRAFORM_TEST"
    duration = 15
  }
}

resource "sysdig_monitor_alert_metric" "sample" {
  name        = "[Kubernetes] CrashLoopBackOff"
  description = "A Kubernetes pod failed to restart"
  severity    = 6

  metric                = "sum(timeAvg(kubernetes.pod.restart.count)) > 2"
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


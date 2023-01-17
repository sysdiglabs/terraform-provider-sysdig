package sysdig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sysdig_secure_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_API_TOKEN", nil),
			},
			"sysdig_secure_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_URL", "https://secure.sysdig.com"),
			},
			"sysdig_secure_insecure_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_INSECURE_TLS", false),
			},
			"sysdig_monitor_api_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_MONITOR_API_TOKEN", nil),
			},
			"sysdig_monitor_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_MONITOR_URL", "https://app.sysdigcloud.com"),
			},
			"sysdig_monitor_insecure_tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_MONITOR_INSECURE_TLS", false),
			},
			"extra_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sysdig_user": resourceSysdigUser(),

			"sysdig_secure_policy":                         resourceSysdigSecurePolicy(),
			"sysdig_secure_notification_channel_email":     resourceSysdigSecureNotificationChannelEmail(),
			"sysdig_secure_notification_channel_sns":       resourceSysdigSecureNotificationChannelSNS(),
			"sysdig_secure_notification_channel_opsgenie":  resourceSysdigSecureNotificationChannelOpsGenie(),
			"sysdig_secure_notification_channel_victorops": resourceSysdigSecureNotificationChannelVictorOps(),
			"sysdig_secure_notification_channel_webhook":   resourceSysdigSecureNotificationChannelWebhook(),
			"sysdig_secure_notification_channel_slack":     resourceSysdigSecureNotificationChannelSlack(),
			"sysdig_secure_notification_channel_pagerduty": resourceSysdigSecureNotificationChannelPagerduty(),
			"sysdig_secure_rule_container":                 resourceSysdigSecureRuleContainer(),
			"sysdig_secure_rule_filesystem":                resourceSysdigSecureRuleFilesystem(),
			"sysdig_secure_rule_network":                   resourceSysdigSecureRuleNetwork(),
			"sysdig_secure_rule_process":                   resourceSysdigSecureRuleProcess(),
			"sysdig_secure_rule_syscall":                   resourceSysdigSecureRuleSyscall(),
			"sysdig_secure_rule_falco":                     resourceSysdigSecureRuleFalco(),
			"sysdig_secure_team":                           resourceSysdigSecureTeam(),
			"sysdig_secure_list":                           resourceSysdigSecureList(),
			"sysdig_secure_macro":                          resourceSysdigSecureMacro(),
			"sysdig_secure_vulnerability_exception":        resourceSysdigSecureVulnerabilityException(),
			"sysdig_secure_vulnerability_exception_list":   resourceSysdigSecureVulnerabilityExceptionList(),
			"sysdig_secure_cloud_account":                  resourceSysdigSecureCloudAccount(),
			"sysdig_secure_benchmark_task":                 resourceSysdigSecureBenchmarkTask(),

			"sysdig_monitor_alert_downtime":                 resourceSysdigMonitorAlertDowntime(),
			"sysdig_monitor_alert_metric":                   resourceSysdigMonitorAlertMetric(),
			"sysdig_monitor_alert_event":                    resourceSysdigMonitorAlertEvent(),
			"sysdig_monitor_alert_anomaly":                  resourceSysdigMonitorAlertAnomaly(),
			"sysdig_monitor_alert_group_outlier":            resourceSysdigMonitorAlertGroupOutlier(),
			"sysdig_monitor_alert_promql":                   resourceSysdigMonitorAlertPromql(),
			"sysdig_monitor_alert_v2_event":                 resourceSysdigMonitorAlertV2Event(),
			"sysdig_monitor_alert_v2_metric":                resourceSysdigMonitorAlertV2Metric(),
			"sysdig_monitor_alert_v2_downtime":              resourceSysdigMonitorAlertV2Downtime(),
			"sysdig_monitor_alert_v2_prometheus":            resourceSysdigMonitorAlertV2Prometheus(),
			"sysdig_monitor_dashboard":                      resourceSysdigMonitorDashboard(),
			"sysdig_monitor_notification_channel_email":     resourceSysdigMonitorNotificationChannelEmail(),
			"sysdig_monitor_notification_channel_opsgenie":  resourceSysdigMonitorNotificationChannelOpsGenie(),
			"sysdig_monitor_notification_channel_pagerduty": resourceSysdigMonitorNotificationChannelPagerduty(),
			"sysdig_monitor_notification_channel_slack":     resourceSysdigMonitorNotificationChannelSlack(),
			"sysdig_monitor_notification_channel_sns":       resourceSysdigMonitorNotificationChannelSNS(),
			"sysdig_monitor_notification_channel_victorops": resourceSysdigMonitorNotificationChannelVictorOps(),
			"sysdig_monitor_notification_channel_webhook":   resourceSysdigMonitorNotificationChannelWebhook(),
			"sysdig_monitor_team":                           resourceSysdigMonitorTeam(),
			"sysdig_monitor_cloud_account":                  resourceSysdigMonitorCloudAccount(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sysdig_secure_trusted_cloud_identity": dataSourceSysdigSecureTrustedCloudIdentity(),
			"sysdig_secure_notification_channel":   dataSourceSysdigSecureNotificationChannel(),
			"sysdig_current_user":                  dataSourceSysdigCurrentUser(),
			"sysdig_user":                          dataSourceSysdigUser(),
			"sysdig_secure_connection":             dataSourceSysdigSecureConnection(),

			"sysdig_fargate_workload_agent":                 dataSourceSysdigFargateWorkloadAgent(),
			"sysdig_monitor_notification_channel_pagerduty": dataSourceSysdigMonitorNotificationChannelPagerduty(),
			"sysdig_monitor_notification_channel_email":     dataSourceSysdigMonitorNotificationChannelEmail(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	sysdigClient := &sysdigClients{d: d}
	return sysdigClient, nil
}

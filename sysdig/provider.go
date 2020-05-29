package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
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
		},
		ResourcesMap: map[string]*schema.Resource{
			"sysdig_user": resourceSysdigUser(),
			"sysdig_team": resourceSysdigTeam(),

			"sysdig_secure_policy":               resourceSysdigSecurePolicy(),
			"sysdig_secure_notification_channel": resourceSysdigSecureNotificationChannel(),
			"sysdig_secure_rule_container":       resourceSysdigSecureRuleContainer(),
			"sysdig_secure_rule_filesystem":      resourceSysdigSecureRuleFilesystem(),
			"sysdig_secure_rule_network":         resourceSysdigSecureRuleNetwork(),
			"sysdig_secure_rule_process":         resourceSysdigSecureRuleProcess(),
			"sysdig_secure_rule_syscall":         resourceSysdigSecureRuleSyscall(),
			"sysdig_secure_rule_falco":           resourceSysdigSecureRuleFalco(),
			"sysdig_secure_list":                 resourceSysdigSecureList(),
			"sysdig_secure_macro":                resourceSysdigSecureMacro(),

			"sysdig_monitor_alert_downtime":      resourceSysdigMonitorAlertDowntime(),
			"sysdig_monitor_alert_metric":        resourceSysdigMonitorAlertMetric(),
			"sysdig_monitor_alert_event":         resourceSysdigMonitorAlertEvent(),
			"sysdig_monitor_alert_anomaly":       resourceSysdigMonitorAlertAnomaly(),
			"sysdig_monitor_alert_group_outlier": resourceSysdigMonitorAlertGroupOutlier(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sysdig_secure_notification_channel": dataSourceSysdigSecureNotificationChannel(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	sysdigClient := &sysdigClients{d: d}
	return sysdigClient, nil
}

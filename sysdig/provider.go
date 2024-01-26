package sysdig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type SysdigProvider struct {
	SysdigClient SysdigClients
}

// Used by tests to get the provider
func Provider() *schema.Provider {
	sysdigClient := NewSysdigClients()
	provider := &SysdigProvider{SysdigClient: sysdigClient}
	return provider.Provider()
}

func (p *SysdigProvider) Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sysdig_secure_skip_policyv2msg": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_SKIP_POLICYV2MSG", true),
			},
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
			"sysdig_monitor_team_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_MONITOR_TEAM_ID", nil),
			},
			"sysdig_monitor_team_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_MONITOR_TEAM_NAME", nil),
			},
			"ibm_monitor_iam_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_MONITOR_IAM_URL", nil),
			},
			"ibm_monitor_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_MONITOR_INSTANCE_ID", nil),
			},
			"ibm_monitor_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_MONITOR_API_KEY", nil),
			},
			"sysdig_secure_team_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_TEAM_ID", nil),
			},
			"sysdig_secure_team_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_TEAM_NAME", nil),
			},
			"ibm_secure_iam_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_SECURE_IAM_URL", nil),
			},
			"ibm_secure_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_SECURE_INSTANCE_ID", nil),
			},
			"ibm_secure_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_IBM_SECURE_API_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sysdig_user":                 resourceSysdigUser(),
			"sysdig_group_mapping":        resourceSysdigGroupMapping(),
			"sysdig_group_mapping_config": resourceSysdigGroupMappingConfig(),
			"sysdig_custom_role":          resourceSysdigCustomRole(),
			"sysdig_team_service_account": resourceSysdigTeamServiceAccount(),

			"sysdig_secure_malware_policy":                                resourceSysdigSecureCompositePolicy(),
			"sysdig_secure_custom_policy":                                 resourceSysdigSecureCustomPolicy(),
			"sysdig_secure_managed_policy":                                resourceSysdigSecureManagedPolicy(),
			"sysdig_secure_managed_ruleset":                               resourceSysdigSecureManagedRuleset(),
			"sysdig_secure_policy":                                        resourceSysdigSecurePolicy(),
			"sysdig_secure_notification_channel_email":                    resourceSysdigSecureNotificationChannelEmail(),
			"sysdig_secure_notification_channel_sns":                      resourceSysdigSecureNotificationChannelSNS(),
			"sysdig_secure_notification_channel_opsgenie":                 resourceSysdigSecureNotificationChannelOpsGenie(),
			"sysdig_secure_notification_channel_victorops":                resourceSysdigSecureNotificationChannelVictorOps(),
			"sysdig_secure_notification_channel_webhook":                  resourceSysdigSecureNotificationChannelWebhook(),
			"sysdig_secure_notification_channel_slack":                    resourceSysdigSecureNotificationChannelSlack(),
			"sysdig_secure_notification_channel_pagerduty":                resourceSysdigSecureNotificationChannelPagerduty(),
			"sysdig_secure_notification_channel_prometheus_alert_manager": resourceSysdigSecureNotificationChannelPrometheusAlertManager(),
			"sysdig_secure_notification_channel_team_email":               resourceSysdigSecureNotificationChannelTeamEmail(),
			"sysdig_secure_notification_channel_msteams":                  resourceSysdigSecureNotificationChannelMSTeams(),
			"sysdig_secure_rule_container":                                resourceSysdigSecureRuleContainer(),
			"sysdig_secure_rule_filesystem":                               resourceSysdigSecureRuleFilesystem(),
			"sysdig_secure_rule_network":                                  resourceSysdigSecureRuleNetwork(),
			"sysdig_secure_rule_process":                                  resourceSysdigSecureRuleProcess(),
			"sysdig_secure_rule_syscall":                                  resourceSysdigSecureRuleSyscall(),
			"sysdig_secure_rule_falco":                                    resourceSysdigSecureRuleFalco(),
			"sysdig_secure_team":                                          resourceSysdigSecureTeam(),
			"sysdig_secure_list":                                          resourceSysdigSecureList(),
			"sysdig_secure_macro":                                         resourceSysdigSecureMacro(),
			"sysdig_secure_vulnerability_exception":                       resourceSysdigSecureVulnerabilityException(),
			"sysdig_secure_vulnerability_exception_list":                  resourceSysdigSecureVulnerabilityExceptionList(),
			"sysdig_secure_cloud_account":                                 resourceSysdigSecureCloudAccount(),
			"sysdig_secure_scanning_policy":                               resourceSysdigSecureScanningPolicy(),
			"sysdig_secure_scanning_policy_assignment":                    resourceSysdigSecureScanningPolicyAssignment(),
			"sysdig_secure_cloud_auth_account":                            resourceSysdigSecureCloudauthAccount(),

			"sysdig_monitor_silence_rule":                                  resourceSysdigMonitorSilenceRule(),
			"sysdig_monitor_alert_downtime":                                resourceSysdigMonitorAlertDowntime(),
			"sysdig_monitor_alert_metric":                                  resourceSysdigMonitorAlertMetric(),
			"sysdig_monitor_alert_event":                                   resourceSysdigMonitorAlertEvent(),
			"sysdig_monitor_alert_anomaly":                                 resourceSysdigMonitorAlertAnomaly(),
			"sysdig_monitor_alert_group_outlier":                           resourceSysdigMonitorAlertGroupOutlier(),
			"sysdig_monitor_alert_promql":                                  resourceSysdigMonitorAlertPromql(),
			"sysdig_monitor_alert_v2_event":                                resourceSysdigMonitorAlertV2Event(),
			"sysdig_monitor_alert_v2_metric":                               resourceSysdigMonitorAlertV2Metric(),
			"sysdig_monitor_alert_v2_downtime":                             resourceSysdigMonitorAlertV2Downtime(),
			"sysdig_monitor_alert_v2_prometheus":                           resourceSysdigMonitorAlertV2Prometheus(),
			"sysdig_monitor_alert_v2_change":                               resourceSysdigMonitorAlertV2Change(),
			"sysdig_monitor_alert_v2_form_based_prometheus":                resourceSysdigMonitorAlertV2FormBasedPrometheus(),
			"sysdig_monitor_alert_v2_group_outlier":                        resourceSysdigMonitorAlertV2GroupOutlier(),
			"sysdig_monitor_dashboard":                                     resourceSysdigMonitorDashboard(),
			"sysdig_monitor_notification_channel_email":                    resourceSysdigMonitorNotificationChannelEmail(),
			"sysdig_monitor_notification_channel_opsgenie":                 resourceSysdigMonitorNotificationChannelOpsGenie(),
			"sysdig_monitor_notification_channel_pagerduty":                resourceSysdigMonitorNotificationChannelPagerduty(),
			"sysdig_monitor_notification_channel_slack":                    resourceSysdigMonitorNotificationChannelSlack(),
			"sysdig_monitor_notification_channel_sns":                      resourceSysdigMonitorNotificationChannelSNS(),
			"sysdig_monitor_notification_channel_victorops":                resourceSysdigMonitorNotificationChannelVictorOps(),
			"sysdig_monitor_notification_channel_webhook":                  resourceSysdigMonitorNotificationChannelWebhook(),
			"sysdig_monitor_notification_channel_msteams":                  resourceSysdigMonitorNotificationChannelMSTeams(),
			"sysdig_monitor_notification_channel_google_chat":              resourceSysdigMonitorNotificationChannelGoogleChat(),
			"sysdig_monitor_notification_channel_prometheus_alert_manager": resourceSysdigMonitorNotificationChannelPrometheusAlertManager(),
			"sysdig_monitor_notification_channel_team_email":               resourceSysdigMonitorNotificationChannelTeamEmail(),
			"sysdig_monitor_notification_channel_custom_webhook":           resourceSysdigMonitorNotificationChannelCustomWebhook(),
			"sysdig_monitor_notification_channel_ibm_event_notification":   resourceSysdigMonitorNotificationChannelIBMEventNotification(),
			"sysdig_monitor_notification_channel_ibm_function":             resourceSysdigMonitorNotificationChannelIBMFunction(),
			"sysdig_monitor_team":                                          resourceSysdigMonitorTeam(),
			"sysdig_monitor_cloud_account":                                 resourceSysdigMonitorCloudAccount(),
			"sysdig_secure_posture_zone":                                   resourceSysdigSecurePostureZone(),
			"sysdig_secure_organization":                                   resourceSysdigSecureOrganization(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"sysdig_secure_trusted_cloud_identity":                        dataSourceSysdigSecureTrustedCloudIdentity(),
			"sysdig_secure_notification_channel":                          dataSourceSysdigSecureNotificationChannel(),
			"sysdig_secure_notification_channel_pagerduty":                dataSourceSysdigSecureNotificationChannelPagerduty(),
			"sysdig_secure_notification_channel_email":                    dataSourceSysdigSecureNotificationChannelEmail(),
			"sysdig_secure_notification_channel_opsgenie":                 dataSourceSysdigSecureNotificationChannelOpsGenie(),
			"sysdig_secure_notification_channel_slack":                    dataSourceSysdigSecureNotificationChannelSlack(),
			"sysdig_secure_notification_channel_sns":                      dataSourceSysdigSecureNotificationChannelSNS(),
			"sysdig_secure_notification_channel_victorops":                dataSourceSysdigSecureNotificationChannelVictorOps(),
			"sysdig_secure_notification_channel_webhook":                  dataSourceSysdigSecureNotificationChannelWebhook(),
			"sysdig_secure_notification_channel_msteams":                  dataSourceSysdigSecureNotificationChannelMSTeams(),
			"sysdig_secure_notification_channel_prometheus_alert_manager": dataSourceSysdigSecureNotificationChannelPrometheusAlertManager(),
			"sysdig_secure_notification_channel_team_email":               dataSourceSysdigSecureNotificationChannelTeamEmail(),
			"sysdig_secure_custom_policy":                                 dataSourceSysdigSecureCustomPolicy(),
			"sysdig_secure_managed_policy":                                dataSourceSysdigSecureManagedPolicy(),
			"sysdig_secure_managed_ruleset":                               dataSourceSysdigSecureManagedRuleset(),
			"sysdig_secure_rule_container":                                dataSourceSysdigSecureRuleContainer(),
			"sysdig_secure_rule_falco":                                    dataSourceSysdigSecureRuleFalco(),
			"sysdig_secure_rule_falco_count":                              dataSourceSysdigSecureRuleFalcoCount(),
			"sysdig_secure_rule_filesystem":                               dataSourceSysdigSecureRuleFilesystem(),
			"sysdig_secure_rule_network":                                  dataSourceSysdigSecureRuleNetwork(),
			"sysdig_secure_rule_process":                                  dataSourceSysdigSecureRuleProcess(),
			"sysdig_secure_rule_syscall":                                  dataSourceSysdigSecureRuleSyscall(),
			"sysdig_secure_posture_policies":                              dataSourceSysdigSecurePosturePolicies(),
			"sysdig_secure_custom_role_permissions":                       dataSourceSysdigSecureCustomRolePermissions(),

			"sysdig_current_user":      dataSourceSysdigCurrentUser(),
			"sysdig_user":              dataSourceSysdigUser(),
			"sysdig_secure_connection": dataSourceSysdigSecureConnection(),
			"sysdig_custom_role":       dataSourceSysdigCustomRole(),

			"sysdig_fargate_workload_agent":                                dataSourceSysdigFargateWorkloadAgent(),
			"sysdig_monitor_notification_channel_pagerduty":                dataSourceSysdigMonitorNotificationChannelPagerduty(),
			"sysdig_monitor_notification_channel_email":                    dataSourceSysdigMonitorNotificationChannelEmail(),
			"sysdig_monitor_notification_channel_opsgenie":                 dataSourceSysdigMonitorNotificationChannelOpsGenie(),
			"sysdig_monitor_notification_channel_slack":                    dataSourceSysdigMonitorNotificationChannelSlack(),
			"sysdig_monitor_notification_channel_sns":                      dataSourceSysdigMonitorNotificationChannelSNS(),
			"sysdig_monitor_notification_channel_victorops":                dataSourceSysdigMonitorNotificationChannelVictorOps(),
			"sysdig_monitor_notification_channel_webhook":                  dataSourceSysdigMonitorNotificationChannelWebhook(),
			"sysdig_monitor_notification_channel_msteams":                  dataSourceSysdigMonitorNotificationChannelMSTeams(),
			"sysdig_monitor_notification_channel_google_chat":              dataSourceSysdigMonitorNotificationChannelGoogleChat(),
			"sysdig_monitor_notification_channel_prometheus_alert_manager": dataSourceSysdigMonitorNotificationChannelPrometheusAlertManager(),
			"sysdig_monitor_notification_channel_team_email":               dataSourceSysdigMonitorNotificationChannelTeamEmail(),
			"sysdig_monitor_notification_channel_custom_webhook":           dataSourceSysdigMonitorNotificationChannelCustomWebhook(),
			"sysdig_monitor_notification_channel_ibm_event_notification":   dataSourceSysdigMonitorNotificationChannelIBMEventNotification(),
			"sysdig_monitor_notification_channel_ibm_function":             dataSourceSysdigMonitorNotificationChannelIBMFunction(),
			"sysdig_monitor_custom_role_permissions":                       dataSourceSysdigMonitorCustomRolePermissions(),
		},
		ConfigureContextFunc: p.providerConfigure,
	}
}

func (p *SysdigProvider) providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	p.SysdigClient.Configure(ctx, d)
	return p.SysdigClient, nil
}

package sysdig

const (
	notificationChannelTypeEmail                  = "EMAIL"
	notificationChannelTypeAmazonSNS              = "SNS"
	notificationChannelTypeOpsGenie               = "OPSGENIE"
	notificationChannelTypeVictorOps              = "VICTOROPS"
	notificationChannelTypeWebhook                = "WEBHOOK"
	notificationChannelTypeSlack                  = "SLACK"
	notificationChannelTypePagerduty              = "PAGER_DUTY"
	notificationChannelTypeMSTeams                = "MS_TEAMS"
	notificationChannelTypeGChat                  = "GCHAT"
	notificationChannelTypePrometheusAlertManager = "PROMETHEUS_ALERT_MANAGER"
	notificationChannelTypeTeamEmail              = "TEAM_EMAIL"
	notificationChannelTypeCustomWebhook          = "POWER_WEBHOOK"
	notificationChannelTypeIBMEventNotification   = "IBM_EVENT_NOTIFICATIONS"

	notificationChannelTypeSlackTemplateKeyV1   = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	notificationChannelTypeSlackTemplateKeyV2   = "SLACK_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"
	notificationChannelTypeMSTeamsTemplateKeyV1 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v1"
	notificationChannelTypeMSTeamsTemplateKeyV2 = "MS_TEAMS_SECURE_EVENT_NOTIFICATION_TEMPLATE_METADATA_v2"

	notificationChannelSecureEventNotificationContentSection = "SECURE_EVENT_NOTIFICATION_CONTENT"
)

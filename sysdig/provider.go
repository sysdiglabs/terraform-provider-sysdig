package sysdig

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"sysdig_secure_api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_API_TOKEN", nil),
			},
			"sysdig_secure_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SYSDIG_SECURE_URL", "https://secure.sysdig.com"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sysdig_secure_policy":               resourceSysdigSecurePolicy(),
			"sysdig_secure_notification_channel": resourceSysdigSecureNotificationChannel(),
			"sysdig_secure_rule_container":       resourceSysdigSecureRuleContainer(),
			"sysdig_secure_rule_filesystem":      resourceSysdigSecureRuleFilesystem(),
			"sysdig_secure_rule_network":         resourceSysdigSecureRuleNetwork(),
			"sysdig_secure_rule_process":         resourceSysdigSecureRuleProcess(),
			"sysdig_secure_rule_syscall":         resourceSysdigSecureRuleSyscall(),
			"sysdig_secure_rule_falco":           resourceSysdigSecureRuleFalco(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	sysdigSecureClient := secure.NewSysdigSecureClient(
		d.Get("sysdig_secure_api_token").(string),
		d.Get("sysdig_secure_url").(string),
	)
	return sysdigSecureClient, nil
}

package sysdig

import (
	"strconv"
	"strings"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createPolicyDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"type": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "falco",
			ValidateDiagFunc: validateDiagFunc(validatePolicyType),
		},
		"id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"severity": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"runbook": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"scope": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"enabled": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
			},
		},
		"notification_channels": {
			Type:     schema.TypeSet,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"actions": {
			Type:     schema.TypeList,
			Computed: true,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container": {
						Type:     schema.TypeString,
						Optional: true,
						Computed: true,
					},
					"capture": {
						Type:     schema.TypeList,
						Optional: true,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"seconds_after_event": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"seconds_before_event": {
									Type:     schema.TypeInt,
									Computed: true,
								},
								"name": {
									Type:     schema.TypeString,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func policyDataSourceToResourceData(policy v2.Policy, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(policy.ID))

	_ = d.Set("name", policy.Name)
	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")
	}

	_ = d.Set("description", policy.Description)
	_ = d.Set("severity", policy.Severity)
	_ = d.Set("enabled", policy.Enabled)
	_ = d.Set("scope", policy.Scope)
	_ = d.Set("version", policy.Version)
	_ = d.Set("notification_channels", policy.NotificationChannelIds)
	_ = d.Set("runbook", policy.Runbook)

	actions := []map[string]interface{}{{}}
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			actions[0]["container"] = strings.ToLower(action)
			//d.Set("actions.0.container", strings.ToLower(action))
		} else {
			actions[0]["capture"] = []map[string]interface{}{{
				"seconds_after_event":  action.AfterEventNs / 1000000000,
				"seconds_before_event": action.BeforeEventNs / 1000000000,
				"name":                 action.Name,
			}}
		}
	}

	_ = d.Set("actions", actions)

	rules := []map[string]interface{}{}

	for _, rule := range policy.Rules {
		rules = append(rules, map[string]interface{}{
			"name":    rule.Name,
			"enabled": rule.Enabled,
		})
	}

	_ = d.Set("rules", rules)
}

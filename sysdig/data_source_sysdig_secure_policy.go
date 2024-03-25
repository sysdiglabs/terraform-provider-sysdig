package sysdig

import (
	"context"
	"strconv"
	"strings"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
						Type:     schema.TypeSet,
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
								"filter": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"bucket_name": {
									Type:     schema.TypeString,
									Computed: true,
								},
								"folder": {
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
	_ = d.Set("notification_channels", policy.NotificationChannelIds)
	_ = d.Set("runbook", policy.Runbook)

	actions := []map[string]interface{}{{}}
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			actions[0]["container"] = strings.ToLower(action)
		} else {
			actions[0]["capture"] = []map[string]interface{}{{
				"seconds_after_event":  action.AfterEventNs / 1000000000,
				"seconds_before_event": action.BeforeEventNs / 1000000000,
				"name":                 action.Name,
				"filter":               action.Filter,
				"bucket_name":          action.BucketName,
				"folder":               action.Folder,
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

func commonDataSourceSecurePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}, resourceName string, isPolicyCorrectType func(v2.Policy) bool) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := d.Get("type").(string)

	policies, _, err := client.GetPolicies(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy v2.Policy
	for _, existingPolicy := range policies {
		if existingPolicy.Name == policyName && existingPolicy.Type == policyType {
			if !isPolicyCorrectType(existingPolicy) {
				return diag.Errorf("policy is not a %s", resourceName)
			}
			policy = existingPolicy
			break
		}
	}

	if policy.ID == 0 {
		return diag.Errorf("unable to find %s", resourceName)
	}

	loadedPolicy, _, err := client.GetPolicyByID(ctx, policy.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	policyDataSourceToResourceData(loadedPolicy, d)

	return nil
}

package sysdig

import (
	"context"
	"strconv"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureCompositePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureCompositePolicyRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createCompositePolicyDataSourceSchema(),
	}
}

func dataSourceSysdigSecureCompositePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return commonCompositePolicyDataSourceSecurePolicyRead(ctx, d, meta, "custom policy", isCustomCompositePolicy)
}

func isCustomCompositePolicy(policy v2.PolicyRulesComposite) bool {
	return !policy.Policy.IsDefault && policy.Policy.TemplateId == 0
}

func createCompositePolicyDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// IMPORTANT: Type is implicit: It's automatically added upon conversion to JSON
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name":                  NameSchema(),
		"description":           DescriptionComputedSchema(),
		"enabled":               EnabledComputedSchema(),
		"severity":              SeverityComputedSchema(),
		"scope":                 ScopeComputedSchema(),
		"version":               VersionSchema(),
		"notification_channels": NotificationChannelsComputedSchema(),
		"runbook":               RunbookComputedSchema(),
		"rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id":          ReadOnlyIntSchema(),
					"name":        ReadOnlyStringSchema(),
					"enabled":     EnabledComputedSchema(),
					"description": DescriptionComputedSchema(),
					"tags":        TagsComputedSchema(),
					"details": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"use_managed_hashes": BoolComputedSchema(),
								"additional_hashes":  HashesComputedSchema(),
								"ignore_hashes":      HashesComputedSchema(),
							},
						},
					},
				},
			},
		},
		"actions": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"prevent_malware": PreventMalwareActionComputedSchema(),
					"container":       ContainerActionComputedSchema(),
					"capture":         CaptureActionComputedSchema(),
				},
			},
		},
	}
}

// TODO: Move into common repo
func compositePolicyDataSourceActionsToResourceData(items []v2.Action) []map[string]interface{} {
	actions := []map[string]interface{}{{}}
	for _, action := range items {
		if action.Type == "POLICY_ACTION_PREVENT_MALWARE" {
			actions[0]["prevent_malware"] = true // TODO
		} else if action.Type == "POLICY_ACTION_PAUSE" || action.Type == "POLICY_ACTION_STOP" || action.Type == "POLICY_ACTION_KILL" { // TODO: Refactor
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
	return actions
}

func compositePolicyDataSourceToResourceData(policy v2.PolicyRulesComposite, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(policy.Policy.ID))

	_ = d.Set("name", policy.Policy.Name)
	if policy.Policy.Type != "" {
		_ = d.Set("type", policy.Policy.Type)
	} else {
		// _ = d.Set("type", "falco") // TODO
	}

	_ = d.Set("description", policy.Policy.Description)
	_ = d.Set("severity", policy.Policy.Severity)
	_ = d.Set("enabled", policy.Policy.Enabled)
	_ = d.Set("scope", policy.Policy.Scope)
	_ = d.Set("notification_channels", policy.Policy.NotificationChannelIds)
	_ = d.Set("runbook", policy.Policy.Runbook)

	actions := compositePolicyDataSourceActionsToResourceData(policy.Policy.Actions)
	_ = d.Set("actions", actions)

	if len(policy.Rules) == 0 {
		panic("policy.Rules is 0")
	}

	// TODO: Exract into a func
	enabledByRuleName := map[string]bool{}
	for _, rule := range policy.Policy.Rules {
		enabledByRuleName[rule.Name] = rule.Enabled
	}

	// TODO: Extract into a function and reuse in resource impl
	rules := []map[string]interface{}{}
	for _, rule := range policy.Rules {
		additionalHashes := []map[string]interface{}{}
		for k, v := range rule.Details.(*v2.MalwareRuleDetails).AdditionalHashes {
			additionalHashes = append(additionalHashes, map[string]interface{}{
				"hash":         k,
				"hash_aliases": v,
			})
		}

		if len(additionalHashes) == 0 {
			panic("additional hashes is 0")
		}

		// TODO: Refactor
		ignoreHashes := []map[string]interface{}{}
		for k, v := range rule.Details.(*v2.MalwareRuleDetails).IgnoreHashes {
			ignoreHashes = append(ignoreHashes, map[string]interface{}{
				"hash":         k,
				"hash_aliases": v,
			})
		}

		rules = append(rules, map[string]interface{}{
			"id":          rule.Id,
			"name":        rule.Name,
			"enabled":     enabledByRuleName[rule.Name],
			"description": rule.Description,
			"tags":        rule.Tags,
			"details": []map[string]interface{}{{
				"use_managed_hashes": rule.Details.(*v2.MalwareRuleDetails).UseManagedHashes,
				"additional_hashes":  additionalHashes,
				// "ignore_hashes":      ignoreHashes,
			}},
		})
	}

	// TODO
	_ = d.Set("rules", rules)
}

func commonCompositePolicyDataSourceSecurePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}, resourceName string, validationFunc func(v2.PolicyRulesComposite) bool) diag.Diagnostics {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := "malware" // d.Get("type").(string) // TODO: Okay to assume it's "malware" type

	policies, _, err := client.FilterCompositePoliciesByNameAndType(ctx, policyType, policyName)
	if err != nil {
		return diag.FromErr(err)
	}

	var policy v2.PolicyRulesComposite
	for _, existingPolicy := range policies {
		tflog.Info(ctx, "***", map[string]interface{}{"->": existingPolicy.Policy.Name})
		if existingPolicy.Policy.Name == policyName && existingPolicy.Policy.Type == policyType {
			if !validationFunc(existingPolicy) {
				return diag.Errorf("policy is not a %s", resourceName)
			}
			policy = existingPolicy
			break
		}
	}

	if policy.Policy == nil {
		return diag.Errorf("unable to find policy %s", resourceName)
	}

	if policy.Policy.ID == 0 {
		return diag.Errorf("unable to find %s", resourceName)
	}

	compositePolicyDataSourceToResourceData(policy, d)

	return nil
}

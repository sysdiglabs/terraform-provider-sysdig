package sysdig

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	policyTypeMalware = "malware"
)

type Target interface {
	*schema.ResourceData | *v2.PolicyRulesComposite
}

type Source interface {
	schema.ResourceData | v2.PolicyRulesComposite
}

func Reducer[T Target, S Source](reducers ...func(T, S) error) func(T, S) error {
	return func(target T, source S) error {
		return Reduce(target, source, reducers...)
	}
}

func Reduce[T Target, S Source](target T, source S, reducers ...func(T, S) error) error {
	for _, f := range reducers {
		if f != nil {
			err := f(target, source)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func setTFResourceBaseAttrs(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	d.SetId(strconv.Itoa(policy.Policy.ID))

	_ = d.Set("version", policy.Policy.Version)

	_ = d.Set("name", policy.Policy.Name)
	_ = d.Set("description", policy.Policy.Description)
	_ = d.Set("severity", policy.Policy.Severity)
	_ = d.Set("enabled", policy.Policy.Enabled)
	_ = d.Set("scope", policy.Policy.Scope)
	_ = d.Set("runbook", policy.Policy.Runbook)

	_ = d.Set("notification_channels", policy.Policy.NotificationChannelIds)

	return nil
}

func setTFResourceAdditionalAttrsMalware(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if policy.Policy.Type != "" {
		_ = d.Set("type", policy.Policy.Type)
	} else {
		_ = d.Set("type", policyTypeMalware)
	}

	return nil
}

func setTFResourcePolicyRulesMalware(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("The policy must have at least one rule attached to it")
	}

	rules := []map[string]interface{}{}
	for _, rule := range policy.Rules {
		additionalHashes := []map[string]interface{}{}
		for k, v := range rule.Details.(*v2.MalwareRuleDetails).AdditionalHashes {
			additionalHashes = append(additionalHashes, map[string]interface{}{
				"hash":         k,
				"hash_aliases": v,
			})
		}

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
			"description": rule.Description,
			"tags":        rule.Tags,
			"details": []map[string]interface{}{{
				"use_managed_hashes": rule.Details.(*v2.MalwareRuleDetails).UseManagedHashes,
				"additional_hashes":  additionalHashes,
				"ignore_hashes":      ignoreHashes,
			}},
		})
	}

	_ = d.Set("rules", rules)

	return nil
}

func setTFResourcePolicyActionsMalware(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	actions := []map[string]interface{}{{}}
	preventMalware := false
	for _, action := range policy.Policy.Actions {
		if action.Type == "POLICY_ACTION_PREVENT_MALWARE" {
			actions[0]["prevent_malware"] = true
			preventMalware = true
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

	// If prevent_malware was updated from true to false, ensure TF resource knows that
	if !preventMalware {
		actions[0]["prevent_malware"] = false
	}

	currentContainerAction := d.Get("actions.0.container").(string)
	currentCaptureAction := d.Get("actions.0.capture").([]interface{})
	// If the policy retrieved from service has no actions and the current state is default values,
	// then do not set the "actions" key as it may cause terraform to think there has been a state change
	if len(policy.Policy.Actions) > 0 || currentContainerAction != "" || len(currentCaptureAction) > 0 {
		_ = d.Set("actions", actions)
	}

	return nil
}

// var malwareTFResourceReducer interface{}// func(*schema.ResourceData, v2.PolicyRulesComposite)
var malwareTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourceAdditionalAttrsMalware,
	setTFResourcePolicyActionsMalware,
	setTFResourcePolicyRulesMalware,
)

func setPolicyBaseAttrs(policy *v2.PolicyRulesComposite, d schema.ResourceData) error {
	id, err := strconv.Atoi(d.Id())
	if err == nil && id != 0 {
		policy.Policy.ID = id
		policy.Policy.Version = d.Get("version").(int)
	}

	policy.Policy.Type = policyTypeMalware

	policy.Policy.Name = d.Get("name").(string)
	policy.Policy.Enabled = d.Get("enabled").(bool)

	policy.Policy.Description = d.Get("description").(string)
	policy.Policy.Severity = d.Get("severity").(int)

	policy.Policy.Runbook = d.Get("runbook").(string)
	policy.Policy.Scope = d.Get("scope").(string)

	policy.Policy.NotificationChannelIds = []int{}
	notificationChannelIdSet := d.Get("notification_channels").(*schema.Set)
	for _, id := range notificationChannelIdSet.List() {
		policy.Policy.NotificationChannelIds = append(policy.Policy.NotificationChannelIds, id.(int))
	}

	return nil
}

func setPolicyActionsMalware(policy *v2.PolicyRulesComposite, d schema.ResourceData) error {
	addActionsToPolicy(&d, policy.Policy)
	return nil
}

func setPolicyRulesMalware(policy *v2.PolicyRulesComposite, d schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rules"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		additionalHashes := map[string][]string{}
		if items, ok := d.GetOk("rules.0.details.0.additional_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]interface{}) {
				item := item.(map[string]interface{})
				k := item["hash"].(string)
				v := schemaSetToList(item["hash_aliases"])
				additionalHashes[k] = v
			}
		}

		// TODO: Extract into a function
		ignoreHashes := map[string][]string{}
		if items, ok := d.GetOk("rules.0.details.0.ignore_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]interface{}) {
				item := item.(map[string]interface{})
				k := item["hash"].(string)
				v := schemaSetToList(item["hash_aliases"])
				ignoreHashes[k] = v
			}
		}

		tags := schemaSetToList(d.Get("rules.0.tags"))
		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rules.0.name").(string),
			Description: d.Get("rules.0.description").(string),
			Tags:        tags,
			Details: v2.MalwareRuleDetails{
				RuleType:         v2.ElementType("MALWARE"), // TODO: Use const
				UseManagedHashes: d.Get("rules.0.details.0.use_managed_hashes").(bool),
				AdditionalHashes: additionalHashes,
				IgnoreHashes:     ignoreHashes,
			},
		}

		id := v2.FlexInt(d.Get("rules.0.id").(int))
		if int(id) != 0 {
			rule.Id = &id
		} else {
			return errors.New((fmt.Sprintf("id is nil: %s, %s", d.Get("rules.0.name"), d.Get("rules.0.id"))))
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

// var malwarePolicyReducer func(*v2.PolicyRulesComposite, schema.ResourceData)
var malwarePolicyReducer = Reducer(
	setPolicyBaseAttrs,
	setPolicyActionsMalware,
	setPolicyRulesMalware,
)

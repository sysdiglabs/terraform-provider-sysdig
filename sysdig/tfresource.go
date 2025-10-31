package sysdig

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	policyTypeMalware = "malware"
	policyTypeDrift   = "drift"
	policyTypeML      = "machine_learning"
	policyTypeAWSML   = "aws_machine_learning"
	policyTypeOktaML  = "okta_machine_learning"

	preventMalwareKey = "prevent_malware"
	preventDriftKey   = "prevent_drift"

	defaultMalwareTag = "malware"
	defaultDriftTag   = "drift"
	defaultMLTag      = "machine_learning"

	driftElementType = "DRIFT"
)

type Target interface {
	*schema.ResourceData | *v2.PolicyRulesComposite
}

type Source interface {
	// copylocks: Do not pass lock by value:
	// schema.ResourceData contains sync.Once contains sync.Mutex (govet)
	*schema.ResourceData | v2.PolicyRulesComposite
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

func schemaSetToList(values any) []string {
	v := values.(*schema.Set).List()

	x := make([]string, len(v))
	for i := range v {
		x[i] = v[i].(string)
	}
	return x
}

func toIntPtr(value any) *int {
	ptr := new(int)
	v, ok := value.(int)
	if ok {
		*ptr = v
	}
	return ptr
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

func setTFResourcePolicyType(policyType string) func(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	return func(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
		if policy.Policy.Type != "" {
			_ = d.Set("type", policy.Policy.Type)
		} else {
			_ = d.Set("type", policyType)
		}

		return nil
	}
}

func setTFResourcePolicyRulesMalware(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("the policy must have at least one rule attached to it")
	}

	rules := []map[string]any{}
	for _, rule := range policy.Rules {
		malwareRuleDetails := rule.Details.(*v2.MalwareRuleDetails)

		additionalHashes := []string{}
		for k := range malwareRuleDetails.AdditionalHashes {
			additionalHashes = append(additionalHashes, k)
		}
		slices.Sort(additionalHashes)

		ignoreHashes := []string{}
		for k := range malwareRuleDetails.IgnoreHashes {
			ignoreHashes = append(ignoreHashes, k)
		}
		slices.Sort(ignoreHashes)

		ignorePaths := []string{}
		for k := range malwareRuleDetails.IgnorePaths {
			ignorePaths = append(ignorePaths, k)
		}
		slices.Sort(ignorePaths)

		rules = append(rules, map[string]any{
			"id":                 rule.ID,
			"name":               rule.Name,
			"description":        rule.Description,
			"version":            rule.Version,
			"tags":               rule.Tags,
			"use_managed_hashes": malwareRuleDetails.UseManagedHashes,
			"use_yara_rules":     malwareRuleDetails.UseYaraRules,
			"additional_hashes":  additionalHashes,
			"ignore_hashes":      ignoreHashes,
			"use_regex":          malwareRuleDetails.UseRegex,
			"ignore_paths":       ignorePaths,
		})
	}

	_ = d.Set("rule", rules)

	return nil
}

func setTFResourcePolicyRulesDrift(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("the policy must have at least one rule attached to it")
	}

	var rules []map[string]any
	for _, rule := range policy.Rules {
		driftDetails, ok := rule.Details.(*v2.DriftRuleDetails)
		if !ok {
			return errors.New("unexpected rule details type, expected DriftRuleDetails")
		}

		// Directly use fields assuming backend returns zero values (not nil)
		exceptionsItems := driftDetails.Exceptions.Items
		exceptionsMatchItems := driftDetails.Exceptions.MatchItems

		var exceptionsBlock []map[string]any
		if len(exceptionsItems) > 0 || exceptionsMatchItems {
			exceptionsBlock = []map[string]any{
				{
					"items":       exceptionsItems,
					"match_items": exceptionsMatchItems,
				},
			}
		}

		prohibitedItems := driftDetails.ProhibitedBinaries.Items
		prohibitedMatchItems := driftDetails.ProhibitedBinaries.MatchItems

		var prohibitedBinariesBlock []map[string]any
		if len(prohibitedItems) > 0 || prohibitedMatchItems {
			prohibitedBinariesBlock = []map[string]any{
				{
					"items":       prohibitedItems,
					"match_items": prohibitedMatchItems,
				},
			}
		}

		processBasedExceptionsItems := driftDetails.ProcessBasedExceptions.Items
		processBasedExceptionMatchItems := driftDetails.ProcessBasedExceptions.MatchItems

		var processBasedExceptionsBlock []map[string]any
		if len(processBasedExceptionsItems) > 0 || processBasedExceptionMatchItems {
			processBasedExceptionsBlock = []map[string]any{
				{
					"items":       processBasedExceptionsItems,
					"match_items": processBasedExceptionMatchItems,
				},
			}
		}

		processBasedProhibitedBinariesItems := driftDetails.ProcessBasedDenylist.Items
		processBasedProhibitedBinariesMatchItems := driftDetails.ProcessBasedDenylist.MatchItems

		var processBasedProhibitedBinariesBlock []map[string]any
		if len(processBasedProhibitedBinariesItems) > 0 || processBasedProhibitedBinariesMatchItems {
			processBasedProhibitedBinariesBlock = []map[string]any{
				{
					"items":       processBasedProhibitedBinariesItems,
					"match_items": processBasedProhibitedBinariesMatchItems,
				},
			}
		}

		mode := driftDetails.Mode
		enabled := (mode != "disabled")

		ruleMap := map[string]any{
			"id":                           rule.ID,
			"name":                         rule.Name,
			"description":                  rule.Description,
			"version":                      rule.Version,
			"tags":                         rule.Tags,
			"enabled":                      enabled,
			"mounted_volume_drift_enabled": driftDetails.MountedVolumeDriftEnabled,
			"use_regex":                    driftDetails.UseRegex,
		}

		if exceptionsBlock != nil {
			ruleMap["exceptions"] = exceptionsBlock
		}
		if prohibitedBinariesBlock != nil {
			ruleMap["prohibited_binaries"] = prohibitedBinariesBlock
		}
		if processBasedExceptionsBlock != nil {
			ruleMap["process_based_exceptions"] = processBasedExceptionsBlock
		}
		if processBasedProhibitedBinariesBlock != nil {
			ruleMap["process_based_prohibited_binaries"] = processBasedProhibitedBinariesBlock
		}

		rules = append(rules, ruleMap)
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	return nil
}

func setTFResourcePolicyRulesML(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("the policy must have at least one rule attached to it")
	}

	rules := []map[string]any{}
	for _, rule := range policy.Rules {
		// Only a single block of anomaly detection trigger and cryptomining trigger is allowed
		// anomalyDetectionTrigger := []map[string]any{{
		// 	"enabled":   rule.Details.(*v2.MLRuleDetails).AnomalyDetectionTrigger.Enabled,
		// 	"threshold": rule.Details.(*v2.MLRuleDetails).AnomalyDetectionTrigger.Threshold,
		// 	"severity":  rule.Details.(*v2.MLRuleDetails).AnomalyDetectionTrigger.Severity,
		// }}

		cryptominingTrigger := []map[string]any{{
			"enabled":   rule.Details.(*v2.MLRuleDetails).CryptominingTrigger.Enabled,
			"threshold": rule.Details.(*v2.MLRuleDetails).CryptominingTrigger.Threshold,
		}}

		rules = append(rules, map[string]any{
			"id":                   rule.ID,
			"name":                 rule.Name,
			"description":          rule.Description,
			"version":              rule.Version,
			"tags":                 rule.Tags,
			"cryptomining_trigger": cryptominingTrigger,
		})
	}

	_ = d.Set("rule", rules)

	return nil
}

func setTFResourcePolicyRulesAWSML(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("the policy must have at least one rule attached to it")
	}

	rules := []map[string]any{}
	for _, rule := range policy.Rules {
		anomalousConsoleLogin := []map[string]any{{
			"enabled":   rule.Details.(*v2.AWSMLRuleDetails).AnomalousConsoleLogin.Enabled,
			"threshold": rule.Details.(*v2.AWSMLRuleDetails).AnomalousConsoleLogin.Threshold,
		}}

		rules = append(rules, map[string]any{
			"id":                      rule.ID,
			"name":                    rule.Name,
			"description":             rule.Description,
			"version":                 rule.Version,
			"tags":                    rule.Tags,
			"anomalous_console_login": anomalousConsoleLogin,
		})
	}

	_ = d.Set("rule", rules)

	return nil
}

func setTFResourcePolicyRulesOktaML(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	if len(policy.Rules) == 0 {
		return errors.New("the policy must have at least one rule attached to it")
	}

	rules := []map[string]any{}
	for _, rule := range policy.Rules {
		anomalousLogin := []map[string]any{{
			"enabled":   rule.Details.(*v2.OktaMLRuleDetails).AnomalousConsoleLogin.Enabled,
			"threshold": rule.Details.(*v2.OktaMLRuleDetails).AnomalousConsoleLogin.Threshold,
		}}

		rules = append(rules, map[string]any{
			"id":                      rule.ID,
			"name":                    rule.Name,
			"description":             rule.Description,
			"version":                 rule.Version,
			"tags":                    rule.Tags,
			"anomalous_console_login": anomalousLogin,
		})
	}

	if err := d.Set("rule", rules); err != nil {
		return err
	}

	return nil
}

// TODO: Split this func into smaller composable functions
func setTFResourcePolicyActions(key string) func(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
	return func(d *schema.ResourceData, policy v2.PolicyRulesComposite) error {
		actions := []map[string]any{{}}
		prevent := false
		for _, action := range policy.Policy.Actions {
			switch action.Type {
			case "POLICY_ACTION_PREVENT_MALWARE", "POLICY_ACTION_PREVENT_DRIFT":
				actions[0][key] = true
				prevent = true
			case "POLICY_ACTION_PAUSE", "POLICY_ACTION_STOP", "POLICY_ACTION_KILL": // TODO: Refactor
				action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
				actions[0]["container"] = strings.ToLower(action)
			default:
				actions[0]["capture"] = []map[string]any{{
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
		if !prevent {
			actions[0][key] = false
		}

		currentContainerAction := d.Get("actions.0.container").(string)
		currentCaptureAction := d.Get("actions.0.capture").([]any)
		// If the policy retrieved from service has no actions and the current state is default values,
		// then do not set the "actions" key as it may cause terraform to think there has been a state change
		if len(policy.Policy.Actions) > 0 || currentContainerAction != "" || len(currentCaptureAction) > 0 {
			_ = d.Set("actions", actions)
		}

		return nil
	}
}

var malwareTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourcePolicyType(policyTypeMalware),
	setTFResourcePolicyActions(preventMalwareKey),
	setTFResourcePolicyRulesMalware,
)

var driftTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourcePolicyType(policyTypeDrift),
	setTFResourcePolicyActions(preventDriftKey),
	setTFResourcePolicyRulesDrift,
)

var mlTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourcePolicyType(policyTypeML),
	setTFResourcePolicyRulesML,
)

var awsMLTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourcePolicyType(policyTypeAWSML),
	setTFResourcePolicyRulesAWSML,
)

var oktaMLTFResourceReducer = Reducer(
	setTFResourceBaseAttrs,
	setTFResourcePolicyType(policyTypeOktaML),
	setTFResourcePolicyRulesOktaML,
)

func setPolicyBaseAttrs(policyType string) func(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	return func(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
		id, err := strconv.Atoi(d.Id())
		if err == nil && id != 0 {
			policy.Policy.ID = id
		}

		v := d.Get("version").(int)
		if v != 0 {
			// Version can only be provided when updating existing policies
			policy.Policy.Version = v
		}

		policy.Policy.Type = policyType

		policy.Policy.Name = d.Get("name").(string)
		policy.Policy.Enabled = d.Get("enabled").(bool)

		policy.Policy.Description = d.Get("description").(string)
		policy.Policy.Severity = d.Get("severity").(int)

		policy.Policy.Runbook = d.Get("runbook").(string)
		policy.Policy.Scope = d.Get("scope").(string)

		policy.Policy.NotificationChannelIds = []int{}
		notificationChannelIDSet := d.Get("notification_channels").(*schema.Set)
		for _, id := range notificationChannelIDSet.List() {
			policy.Policy.NotificationChannelIds = append(policy.Policy.NotificationChannelIds, id.(int))
		}

		return nil
	}
}

func setPolicyActions(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	addActionsToPolicy(d, policy.Policy)
	return nil
}

func setPolicyRulesMalware(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rule"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		additionalHashes := map[string][]string{}
		if items, ok := d.GetOk("rule.0.additional_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]any) {
				hash := item.(string)
				additionalHashes[hash] = []string{}
			}
		}

		// TODO: Extract into a function
		ignoreHashes := map[string][]string{}
		if items, ok := d.GetOk("rule.0.ignore_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]any) {
				hash := item.(string)
				ignoreHashes[hash] = []string{}
			}
		}

		tags := schemaSetToList(d.Get("rule.0.tags"))
		// Set default tags as field tags must not be null
		if len(tags) == 0 {
			tags = []string{defaultMalwareTag}
		}

		ignorePaths := map[string][]string{}
		if items, ok := d.GetOk("rule.0.ignore_paths"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]any) {
				path := item.(string)
				ignorePaths[path] = []string{}
			}
		}

		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rule.0.name").(string),
			Description: d.Get("rule.0.description").(string),
			Tags:        tags,
			Details: v2.MalwareRuleDetails{
				RuleType:         v2.ElementType("MALWARE"), // TODO: Use const
				UseManagedHashes: d.Get("rule.0.use_managed_hashes").(bool),
				UseYaraRules:     d.Get("rule.0.use_yara_rules").(bool),
				AdditionalHashes: additionalHashes,
				IgnoreHashes:     ignoreHashes,
				UseRegex:         d.Get("rule.0.use_regex").(bool),
				IgnorePaths:      ignorePaths,
			},
		}

		id := v2.FlexInt(d.Get("rule.0.id").(int))
		if int(id) != 0 {
			rule.ID = &id
		}

		v := toIntPtr(d.Get("rule.0.version"))
		if *v != 0 {
			// Version can only be provided when updating existing rules
			rule.Version = v
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

func setPolicyRulesDrift(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rule"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		exceptions := extractIntoRuntimePolicyRuleList("rule.0.exceptions", d)

		prohibitedBinaries := extractIntoRuntimePolicyRuleList("rule.0.prohibited_binaries", d)

		processBasedExceptions := extractIntoRuntimePolicyRuleList("rule.0.process_based_exceptions", d)

		processBasedProhibitedBinaries := extractIntoRuntimePolicyRuleList("rule.0.process_based_prohibited_binaries", d)

		tags := schemaSetToList(d.Get("rule.0.tags"))
		// Set default tags as field tags must not be null
		if len(tags) == 0 {
			tags = []string{defaultDriftTag}
		}

		enabled := d.Get("rule.0.enabled").(bool)
		mode := "enabled"
		if !enabled {
			mode = "disabled"
		}

		mountedVolumeDriftEnabled := d.Get("rule.0.mounted_volume_drift_enabled").(bool)
		useRegex := d.Get("rule.0.use_regex").(bool)

		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rule.0.name").(string),
			Description: d.Get("rule.0.description").(string),
			Tags:        tags,
			Details: v2.DriftRuleDetails{
				RuleType:                  v2.ElementType(driftElementType), // TODO: Use const
				Mode:                      mode,
				Exceptions:                &exceptions,
				ProhibitedBinaries:        &prohibitedBinaries,
				ProcessBasedExceptions:    &processBasedExceptions,
				ProcessBasedDenylist:      &processBasedProhibitedBinaries,
				MountedVolumeDriftEnabled: mountedVolumeDriftEnabled,
				UseRegex:                  useRegex,
			},
		}

		id := v2.FlexInt(d.Get("rule.0.id").(int))
		if int(id) != 0 {
			rule.ID = &id
		}

		v := toIntPtr(d.Get("rule.0.version"))
		if *v != 0 {
			// Version can only be provided when updating existing rules
			rule.Version = v
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

func extractIntoRuntimePolicyRuleList(key string, d *schema.ResourceData) v2.RuntimePolicyRuleList {
	if _, ok := d.GetOk(key); ok {
		items := schemaSetToList(d.Get(key + ".0.items"))
		matchItems := d.Get(key + ".0.match_items").(bool)

		return v2.RuntimePolicyRuleList{
			Items:      items,
			MatchItems: matchItems,
		}
	}
	return v2.RuntimePolicyRuleList{
		Items:      []string{},
		MatchItems: false,
	}
}

func setPolicyRulesML(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rule"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		// TODO: Extract into a function
		cryptominingTrigger := &v2.MLRuleThresholdAndSeverity{}
		if _, ok := d.GetOk("rule.0.cryptomining_trigger"); ok { // TODO: Do not hardcode the indexes
			cryptominingTrigger.Enabled = d.Get("rule.0.cryptomining_trigger.0.enabled").(bool)
			cryptominingTrigger.Threshold = float64(d.Get("rule.0.cryptomining_trigger.0.threshold").(int))
		}
		anomalyDetectionTrigger := &v2.MLRuleThresholdAndSeverity{}

		tags := schemaSetToList(d.Get("rule.0.tags"))
		// Set default tags as field tags must not be null
		if len(tags) == 0 {
			tags = []string{defaultMLTag}
		}

		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rule.0.name").(string),
			Description: d.Get("rule.0.description").(string),
			// IMPORTANT: In order to update an ML policy,
			// correct version number must be provided
			Tags: tags,
			Details: v2.MLRuleDetails{
				RuleType:                v2.ElementType("MACHINE_LEARNING"), // TODO: Use const
				CryptominingTrigger:     cryptominingTrigger,
				AnomalyDetectionTrigger: anomalyDetectionTrigger,
			},
		}

		id := v2.FlexInt(d.Get("rule.0.id").(int))
		if int(id) != 0 {
			rule.ID = &id
		}

		v := toIntPtr(d.Get("rule.0.version"))
		if *v != 0 {
			// Version can only be provided when updating existing rules
			rule.Version = v
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

func setPolicyRulesAWSML(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rule"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		anomalousConsoleLogin := &v2.MLRuleThresholdAndSeverity{}
		if _, ok := d.GetOk("rule.0.anomalous_console_login"); ok { // TODO: Do not hardcode the indexes
			anomalousConsoleLogin.Enabled = d.Get("rule.0.anomalous_console_login.0.enabled").(bool)
			anomalousConsoleLogin.Threshold = float64(d.Get("rule.0.anomalous_console_login.0.threshold").(int))
		}

		tags := schemaSetToList(d.Get("rule.0.tags"))
		// Set default tags as field tags must not be null
		if len(tags) == 0 {
			tags = []string{defaultMLTag}
		}

		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rule.0.name").(string),
			Description: d.Get("rule.0.description").(string),
			// IMPORTANT: In order to update an ML policy,
			// correct version number must be provided
			Tags: tags,
			Details: v2.AWSMLRuleDetails{
				RuleType:              v2.ElementType("AWS_MACHINE_LEARNING"), // TODO: Use const
				AnomalousConsoleLogin: anomalousConsoleLogin,
			},
		}

		id := v2.FlexInt(d.Get("rule.0.id").(int))
		if int(id) != 0 {
			rule.ID = &id
		}

		v := toIntPtr(d.Get("rule.0.version"))
		if *v != 0 {
			// Version can only be provided when updating existing rules
			rule.Version = v
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

func setPolicyRulesOktaML(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rule"); ok {

		anomalousLogin := &v2.MLRuleThresholdAndSeverity{}
		if _, ok := d.GetOk("rule.0.anomalous_console_login"); ok { // TODO: Do not hardcode the indexes
			anomalousLogin.Enabled = d.Get("rule.0.anomalous_console_login.0.enabled").(bool)
			anomalousLogin.Threshold = float64(d.Get("rule.0.anomalous_console_login.0.threshold").(int))
		}

		tags := schemaSetToList(d.Get("rule.0.tags"))
		// Set default tags as field tags must not be null
		if len(tags) == 0 {
			tags = []string{defaultMLTag}
		}

		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rule.0.name").(string),
			Description: d.Get("rule.0.description").(string),
			// IMPORTANT: In order to update an ML policy,
			// correct version number must be provided
			Tags: tags,
			Details: v2.OktaMLRuleDetails{
				RuleType:              v2.ElementType("OKTA_MACHINE_LEARNING"), // TODO: Use const
				AnomalousConsoleLogin: anomalousLogin,
			},
		}

		id := v2.FlexInt(d.Get("rule.0.id").(int))
		if int(id) != 0 {
			rule.ID = &id
		}

		v := toIntPtr(d.Get("rule.0.version"))
		if *v != 0 {
			// Version can only be provided when updating existing rules
			rule.Version = v
		}

		policy.Rules = append(policy.Rules, rule)
	}
	return nil
}

var malwarePolicyReducer = Reducer(
	setPolicyBaseAttrs(policyTypeMalware),
	setPolicyActions,
	setPolicyRulesMalware,
)

var driftPolicyReducer = Reducer(
	setPolicyBaseAttrs(policyTypeDrift),
	setPolicyActions,
	setPolicyRulesDrift,
)

var mlPolicyReducer = Reducer(
	setPolicyBaseAttrs(policyTypeML),
	setPolicyRulesML,
)

var awsMLPolicyReducer = Reducer(
	setPolicyBaseAttrs(policyTypeAWSML),
	setPolicyRulesAWSML,
)

var oktaMLPolicyReducer = Reducer(
	setPolicyBaseAttrs(policyTypeOktaML),
	setPolicyRulesOktaML,
)

package sysdig

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ruleIdentifier represents a backend field identifier for zone scope filters.
type ruleIdentifier string

const (
	identAccount          ruleIdentifier = "account"
	identClusterId        ruleIdentifier = "clusterId"
	identDistribution     ruleIdentifier = "distribution"
	identNamespace        ruleIdentifier = "namespace"
	identLabels           ruleIdentifier = "labels"
	identLabelValues      ruleIdentifier = "labelValues"
	identAgentTags        ruleIdentifier = "agentTags"
	identLocation         ruleIdentifier = "location"
	identOrganization     ruleIdentifier = "organization"
	identSubscription     ruleIdentifier = "subscription"
	identName             ruleIdentifier = "name"
	identRegistry         ruleIdentifier = "registry"
	identRepository       ruleIdentifier = "repository"
	identGitIntegrationId ruleIdentifier = "gitIntegrationId"
	identGitSourceId      ruleIdentifier = "gitSourceId"
	identResourceGroupId  ruleIdentifier = "resourceGroupId"
	identAccountGroupId   ruleIdentifier = "accountGroupId"
	identAccountGroupName ruleIdentifier = "accountGroupName"
)

// allowedIdentifiers maps each target_type to the set of backend rule identifiers
// it accepts. This mirrors the backend's validIdentifiers map.
var allowedIdentifiers = map[string]map[ruleIdentifier]struct{}{
	"aws":        {identAccount: {}, identOrganization: {}, identLabels: {}, identLocation: {}},
	"gcp":        {identAccount: {}, identOrganization: {}, identLabels: {}, identLocation: {}},
	"azure":      {identAccount: {}, identOrganization: {}, identLabels: {}, identLocation: {}},
	"kubernetes": {identClusterId: {}, identNamespace: {}, identLabelValues: {}, identAgentTags: {}, identDistribution: {}},
	"host":       {identClusterId: {}, identName: {}, identAgentTags: {}},
	"image":      {identRegistry: {}, identRepository: {}},
	"git":        {identGitIntegrationId: {}, identGitSourceId: {}},
	"ibm":        {identAccount: {}, identOrganization: {}, identLabels: {}, identLocation: {}, identResourceGroupId: {}, identAccountGroupId: {}, identAccountGroupName: {}},
	"oci":        {identAccount: {}, identOrganization: {}, identLabels: {}, identLocation: {}},
}

// directFieldMap maps v2 expression field names (without prefixes) to their
// corresponding backend rule identifiers.
var directFieldMap = map[string]ruleIdentifier{
	"account":          identAccount,
	"organization":     identOrganization,
	"clusterId":        identClusterId,
	"namespace":        identNamespace,
	"distribution":     identDistribution,
	"location":         identLocation,
	"name":             identName,
	"registry":         identRegistry,
	"repository":       identRepository,
	"gitIntegrationId": identGitIntegrationId,
	"gitSourceId":      identGitSourceId,
	"resourceGroupId":  identResourceGroupId,
	"accountGroupId":   identAccountGroupId,
	"accountGroupName": identAccountGroupName,
	"subscription":     identSubscription,
}

// labelAsLabelsTargets lists target types where label.<key> maps to
// LabelsRuleIdentifier. All other target types use LabelValuesRuleIdentifier.
var labelAsLabelsTargets = map[string]struct{}{
	"aws":   {},
	"gcp":   {},
	"azure": {},
	"ibm":   {},
	"oci":   {},
}

// identToFieldPattern maps backend identifiers to user-facing field patterns
// for error messages.
var identToFieldPattern = map[ruleIdentifier]string{
	identAccount:          "account",
	identOrganization:     "organization",
	identClusterId:        "clusterId",
	identNamespace:        "namespace",
	identDistribution:     "distribution",
	identLabels:           "label.<key>",
	identLabelValues:      "label.<key>",
	identAgentTags:        "agent.tag.<key>",
	identLocation:         "location",
	identName:             "name",
	identRegistry:         "registry",
	identRepository:       "repository",
	identGitIntegrationId: "gitIntegrationId",
	identGitSourceId:      "gitSourceId",
	identResourceGroupId:  "resourceGroupId",
	identAccountGroupId:   "accountGroupId",
	identAccountGroupName: "accountGroupName",
	identSubscription:     "subscription",
}

// resolveIdentifier maps a v2 expression field to the backend's rule identifier,
// taking target_type into account for label.<key> disambiguation.
//
// Returns the identifier and true if the field is recognized, or ("", false)
// for fields the provider doesn't know about (forward compatibility).
func resolveIdentifier(targetType, field string) (ruleIdentifier, bool) {
	if strings.HasPrefix(field, "label.") && len(field) > len("label.") {
		if _, ok := labelAsLabelsTargets[targetType]; ok {
			return identLabels, true
		}
		return identLabelValues, true
	}

	if strings.HasPrefix(field, "agent.tag.") && len(field) > len("agent.tag.") {
		return identAgentTags, true
	}

	if id, ok := directFieldMap[field]; ok {
		return id, true
	}

	return "", false
}

// allowedFieldFamilies returns a sorted list of user-facing field patterns
// allowed for a given target_type.
func allowedFieldFamilies(targetType string) []string {
	allowed, ok := allowedIdentifiers[targetType]
	if !ok {
		return nil
	}

	seen := map[string]struct{}{}
	var fields []string
	for id := range allowed {
		f := identToFieldPattern[id]
		if _, dup := seen[f]; !dup {
			seen[f] = struct{}{}
			fields = append(fields, f)
		}
	}
	sort.Strings(fields)
	return fields
}

// validateExpressionsFromPlan validates expression fields using the raw plan's
// cty representation. This avoids type-assertion issues that occur when reading
// nested TypeList elements inside a TypeSet via d.Get() during CustomizeDiff.
func validateExpressionsFromPlan(d *schema.ResourceDiff) error {
	rawPlan := d.GetRawPlan()
	if rawPlan.IsNull() || !rawPlan.IsKnown() {
		return nil
	}

	scopeSet := rawPlan.GetAttr(SchemaScopeKey)
	if scopeSet.IsNull() || !scopeSet.IsKnown() {
		return nil
	}

	scopeIndex := 0
	for it := scopeSet.ElementIterator(); it.Next(); {
		_, scopeVal := it.Element()

		targetTypeVal := scopeVal.GetAttr(SchemaTargetTypeKey)
		if targetTypeVal.IsNull() || !targetTypeVal.IsKnown() {
			scopeIndex++
			continue
		}
		targetType := targetTypeVal.AsString()

		exprsVal := scopeVal.GetAttr(SchemaExpressionKey)
		if exprsVal.IsNull() || !exprsVal.IsKnown() || exprsVal.LengthInt() == 0 {
			scopeIndex++
			continue
		}

		allowed, ok := allowedIdentifiers[targetType]
		if !ok {
			known := make([]string, 0, len(allowedIdentifiers))
			for k := range allowedIdentifiers {
				known = append(known, k)
			}
			sort.Strings(known)
			return fmt.Errorf(
				"scope[%d]: unknown target_type %q; supported types: %s",
				scopeIndex, targetType, strings.Join(known, ", "),
			)
		}

		for i := 0; i < exprsVal.LengthInt(); i++ {
			exprVal := exprsVal.Index(cty.NumberIntVal(int64(i)))
			fieldVal := exprVal.GetAttr(SchemaFieldKey)
			if fieldVal.IsNull() || !fieldVal.IsKnown() {
				continue
			}
			field := fieldVal.AsString()
			if field == "" {
				continue
			}

			ident, recognized := resolveIdentifier(targetType, field)
			if !recognized {
				continue
			}

			if _, ok := allowed[ident]; !ok {
				return fmt.Errorf(
					"scope[%d].expression[%d]: field %q is not allowed for target_type %q; "+
						"allowed fields: %s",
					scopeIndex, i, field, targetType,
					strings.Join(allowedFieldFamilies(targetType), ", "),
				)
			}
		}

		scopeIndex++
	}

	return nil
}

// validateExpressionFields validates that all expression fields in a scope
// are allowed for the given target_type.
//
// Validation rules:
//   - Unknown target_type → error (target_type has a fixed set validated by the schema)
//   - Unrecognized field → skip (forward compatibility with new backend fields)
//   - Recognized field not in allowlist → error
func validateExpressionFields(targetType string, expressions []interface{}, scopeIndex int) error {
	allowed, ok := allowedIdentifiers[targetType]
	if !ok {
		known := make([]string, 0, len(allowedIdentifiers))
		for k := range allowedIdentifiers {
			known = append(known, k)
		}
		sort.Strings(known)
		return fmt.Errorf(
			"scope[%d]: unknown target_type %q; supported types: %s",
			scopeIndex, targetType, strings.Join(known, ", "),
		)
	}

	for i, raw := range expressions {
		expr, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		field, _ := expr[SchemaFieldKey].(string)
		if field == "" {
			continue
		}

		ident, recognized := resolveIdentifier(targetType, field)
		if !recognized {
			continue
		}

		if _, ok := allowed[ident]; !ok {
			return fmt.Errorf(
				"scope[%d].expression[%d]: field %q is not allowed for target_type %q; "+
					"allowed fields: %s",
				scopeIndex, i, field, targetType,
				strings.Join(allowedFieldFamilies(targetType), ", "),
			)
		}
	}

	return nil
}

// validateZoneExpressions runs expression field validation as a safety net
// for Create/Update, complementing the plan-time CustomizeDiff validation.
func validateZoneExpressions(d *schema.ResourceData) error {
	scopeSet, ok := d.Get(SchemaScopeKey).(*schema.Set)
	if !ok || scopeSet == nil {
		return nil
	}

	for i, raw := range scopeSet.List() {
		scope, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		expressions, ok := scope[SchemaExpressionKey].([]interface{})
		if !ok || len(expressions) == 0 {
			continue
		}

		targetType, _ := scope[SchemaTargetTypeKey].(string)
		if err := validateExpressionFields(targetType, expressions, i); err != nil {
			return err
		}
	}

	return nil
}

package sysdig

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// legacyAttributePattern matches legacy v1 attribute names that need migration.
// These are: labelValues, labels, agentTags (without dot notation).
// v2 attributes use dot notation like: label.key, agent.tag.key
var legacyAttributePattern = regexp.MustCompile(`\b(labelValues|labels|agentTags)\b`)

func resourceSysdigSecureZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSysdigSecureZoneCreate,
		ReadContext:   resourceSysdigSecureZoneRead,
		UpdateContext: resourceSysdigSecureZoneUpdate,
		DeleteContext: resourceSysdigSecureZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			scopeSet, ok := d.Get("scope").(*schema.Set)
			if !ok || scopeSet == nil {
				return nil
			}

			for i, raw := range scopeSet.List() {
				scope := raw.(map[string]interface{})

				hasRules := false
				if v, ok := scope["rules"]; ok && v != nil {
					hasRules = v.(string) != ""
				}

				hasExpr := false
				if v, ok := scope["expression"]; ok && v != nil {
					switch expr := v.(type) {
					case []interface{}:
						hasExpr = len(expr) > 0
					case *schema.Set:
						hasExpr = expr.Len() > 0
					}
				}

				if hasRules && hasExpr {
					return fmt.Errorf(
						"scope[%d]: 'rules' cannot be used together with 'expression'",
						i,
					)
				}

				if !hasRules && !hasExpr {
					return fmt.Errorf(
						"scope[%d]: either 'rules' or 'expression' must be specified",
						i,
					)
				}
			}

			// Validate expression fields against the target_type allowlist.
			// Uses GetRawPlan() for reliable type access — nested TypeList
			// elements inside a TypeSet may not materialize as
			// map[string]interface{} during diff computation.
			if err := validateExpressionsFromPlan(d); err != nil {
				return err
			}

			return nil
		},

		Schema: map[string]*schema.Schema{
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaIsSystemKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			SchemaAuthorKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaLastModifiedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaLastUpdated: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaScopeKey: {
				Type:     schema.TypeSet,
				MinItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaIDKey: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						SchemaTargetTypeKey: {
							Type:     schema.TypeString,
							Required: true,
						},
						SchemaRulesKey: {
							Type:     schema.TypeString,
							Optional: true,
							ValidateDiagFunc: func(v interface{}, path cty.Path) diag.Diagnostics {
								rules := v.(string)
								if rules != "" && legacyAttributePattern.MatchString(rules) {
									return diag.Diagnostics{
										diag.Diagnostic{
											Severity: diag.Warning,
											Summary:  "Deprecated legacy rules syntax",
											Detail:   "The 'rules' field with legacy attributes (labels, labelValues, agentTags) is deprecated. Use 'expression' blocks with v2 syntax (label.<key>, agent.tag.<key>) instead. See the documentation for migration guidance.",
										},
									}
								}
								return nil
							},
						},
						SchemaExpressionKey: {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									SchemaFieldKey: {
										Type:     schema.TypeString,
										Required: true,
									},
									SchemaOperatorKey: {
										Type:     schema.TypeString,
										Required: true,
									},
									SchemaValueKey: {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									SchemaValuesKey: {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSysdigSecureZoneCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientV1, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	clientV2, err := getZoneV2Client(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneId, diags := createZone(ctx, d, clientV1, clientV2)
	if diags.HasError() {
		return diags
	}

	d.SetId(fmt.Sprintf("%d", zoneId))
	return resourceSysdigSecureZoneRead(ctx, d, m)
}

func createZone(ctx context.Context, d *schema.ResourceData, clientV1 v2.ZoneInterface, clientV2 v2.ZoneV2Interface) (int, diag.Diagnostics) {
	legacyZone, e := categorizeZone(d)
	if e != nil {
		return 0, diag.FromErr(fmt.Errorf("error analyzing zone scope: %s", e))
	}
	if legacyZone {
		zoneRequest := zoneRequestFromResourceData(d)
		createdZone, err := clientV1.CreateZone(ctx, zoneRequest)
		if err != nil {
			return 0, diag.FromErr(fmt.Errorf("error creating Sysdig Zone: %s", err))
		}
		return createdZone.ID, nil
	}

	if err := validateZoneExpressions(d); err != nil {
		return 0, diag.FromErr(err)
	}

	zone := expandZoneV2(d)
	created, err := clientV2.CreateZoneV2(ctx, zone)
	if err != nil {
		return 0, diag.FromErr(fmt.Errorf("error creating zone: %w", err))
	}
	return created.ID, nil
}

func isNotFound(err error) bool {
	var apiErr *v2.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func categorizeZone(d *schema.ResourceData) (bool, error) {
	rawScopes := d.Get(SchemaScopeKey)
	if rawScopes == nil {
		return false, fmt.Errorf("scope is required and cannot be nil")
	}
	scopes, ok := rawScopes.(*schema.Set)
	if !ok || scopes == nil {
		return false, fmt.Errorf("expected scope to be a *schema.Set, got %T", rawScopes)
	}
	scopeList := scopes.List()

	var hasLegacyRules, hasV2Rules, hasExpr bool

	for _, raw := range scopeList {
		scope := raw.(map[string]any)

		if rules, ok := scope["rules"].(string); ok && rules != "" {
			// Check if rules contain legacy v1 attributes
			if legacyAttributePattern.MatchString(rules) {
				hasLegacyRules = true
			} else {
				hasV2Rules = true
			}
		}

		if expr, ok := scope["expression"].([]interface{}); ok && len(expr) > 0 {
			hasExpr = true
		}
	}

	// Expression blocks always mean v2
	if hasExpr {
		if hasLegacyRules {
			return false, fmt.Errorf("cannot mix expression blocks with legacy v1 rules syntax")
		}
		return false, nil // v2
	}

	// Legacy v1 rules cannot be mixed with v2 rules
	if hasLegacyRules && hasV2Rules {
		return false, fmt.Errorf("cannot mix legacy v1 rules (labelValues, labels, agentTags) with v2 rules (label.key, agent.tag.key)")
	}

	if hasLegacyRules {
		return true, nil // v1
	}

	// V2 rules or no rules (import/read) - treat as v2
	return false, nil
}

func resourceSysdigSecureZoneRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	clientv2, err := getZoneV2Client(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	legacyZone, e := categorizeZone(d)
	if e != nil {
		return diag.FromErr(fmt.Errorf("error analyzing zone scope: %s", e))
	}
	if legacyZone {
		zone, err := client.GetZoneByID(ctx, id)
		if err != nil {
			if isNotFound(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("error reading zone %d: %w", id, err))
		}

		_ = d.Set("name", zone.Name)
		_ = d.Set("description", zone.Description)
		_ = d.Set("is_system", zone.IsSystem)
		_ = d.Set("author", zone.Author)
		_ = d.Set("last_modified_by", zone.LastModifiedBy)
		_ = d.Set("last_updated", time.UnixMilli(zone.LastUpdated).Format(time.RFC3339))
		// For legacy zones, we need to set the rules field in the scope
		if err := d.Set(SchemaScopeKey, fromZoneScopesResponse(zone.Scopes)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting scope: %s", err))
		}
	} else {
		zone, err := clientv2.GetZoneV2(ctx, id)
		if err != nil {
			if isNotFound(err) {
				d.SetId("")
				return nil
			}
			return diag.FromErr(fmt.Errorf("error reading zone %d: %w", id, err))
		}
		_ = d.Set("name", zone.Name)
		_ = d.Set("description", zone.Description)
		_ = d.Set("is_system", zone.IsSystem)
		_ = d.Set("author", zone.Author)
		_ = d.Set("last_modified_by", zone.LastModifiedBy)
		_ = d.Set("last_updated", time.UnixMilli(zone.LastUpdated).Format(time.RFC3339))
		// "State follows config": if the user configured rules, write rules
		// into state; if they configured expressions, write expressions.
		// On the first Read (called from Create), d.Get returns the config
		// values. On subsequent Reads, state reflects the prior Read — which
		// already matched config — so the choice is self-reinforcing.
		preferRules := !stateHasExpressions(d)
		if err := d.Set(SchemaScopeKey, flattenZoneV2(zone, preferRules)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting scope: %s", err))
		}

	}

	return nil
}

func resourceSysdigSecureZoneUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	clientV2, err := getZoneV2Client(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	legacyZone, e := categorizeZone(d)
	if e != nil {
		return diag.FromErr(fmt.Errorf("error analyzing zone scope: %s", e))
	}
	if !legacyZone {
		if err := validateZoneExpressions(d); err != nil {
			return diag.FromErr(err)
		}

		zone := expandZoneV2(d)
		id, _ := strconv.Atoi(d.Id())
		zone.ID = id

		_, err = clientV2.UpdateZoneV2(ctx, zone)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating zone: %w", err))
		}

		return resourceSysdigSecureZoneRead(ctx, d, m)
	} else {
		zoneRequest := zoneRequestFromResourceData(d)

		_, err = client.UpdateZone(ctx, zoneRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating Sysdig Zone: %s", err))
		}
	}

	return resourceSysdigSecureZoneRead(ctx, d, m)
}

func resourceSysdigSecureZoneDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	clientV2, err := getZoneV2Client(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	legacyZone, e := categorizeZone(d)
	if e != nil {
		return diag.FromErr(fmt.Errorf("error analyzing zone scope: %s", e))
	}
	if !legacyZone {
		err = clientV2.DeleteZoneV2(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting Sysdig Zone: %s", err))
		}
	} else {
		err = client.DeleteZone(ctx, id)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error deleting Sysdig Zone: %s", err))
		}
	}

	d.SetId("")
	return nil
}

func zoneRequestFromResourceData(d *schema.ResourceData) *v2.ZoneRequest {
	zoneRequest := &v2.ZoneRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Scopes:      toZoneScopesRequest(d.Get(SchemaScopeKey).(*schema.Set)),
	}

	if d.Id() != "" {
		id, err := strconv.Atoi(d.Id())
		if err == nil {
			zoneRequest.ID = id
		}
	}

	return zoneRequest
}

func toZoneScopesRequest(scopes *schema.Set) []v2.ZoneScope {
	var zoneScopes []v2.ZoneScope
	for _, attr := range scopes.List() {
		s := attr.(map[string]any)
		zoneScopes = append(zoneScopes, v2.ZoneScope{
			ID:         s[SchemaIDKey].(int),
			TargetType: s[SchemaTargetTypeKey].(string),
			Rules:      s[SchemaRulesKey].(string),
		})
	}
	return zoneScopes
}

func fromZoneScopesResponse(scopes []v2.ZoneScope) []map[string]any {
	var flattenedScopes []map[string]any
	for _, scope := range scopes {
		flattenedScopes = append(flattenedScopes, map[string]any{
			SchemaIDKey:         scope.ID,
			SchemaTargetTypeKey: scope.TargetType,
			SchemaRulesKey:      scope.Rules,
		})
	}
	return flattenedScopes
}

func getZoneClient(clients SysdigClients) (v2.ZoneInterface, error) {
	var client v2.ZoneInterface
	var err error
	switch clients.GetClientType() {
	case IBMSecure:
		client, err = clients.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = clients.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func getZoneV2Client(clients SysdigClients) (v2.ZoneV2Interface, error) {
	var client v2.ZoneV2Interface
	var err error
	switch clients.GetClientType() {
	case IBMSecure:
		client, err = clients.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = clients.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

// stateHasExpressions reports whether any scope in the current Terraform state
// contains expression blocks. Used by Read to decide whether to flatten the
// API response as expressions or rules ("state follows config").
func stateHasExpressions(d *schema.ResourceData) bool {
	scopeSet, ok := d.Get(SchemaScopeKey).(*schema.Set)
	if !ok || scopeSet == nil {
		return false
	}
	for _, raw := range scopeSet.List() {
		m := raw.(map[string]interface{})
		if exprs, ok := m[SchemaExpressionKey].([]interface{}); ok && len(exprs) > 0 {
			return true
		}
	}
	return false
}

// flattenZoneV2 flattens the backend ZoneV2 representation into the Terraform schema.
//
// When preferRules is true, the rules string from the backend is written into state
// (matching a user config that uses rules). When false, structured expressions are
// preferred (matching a user config that uses expression blocks).
//
// NOTE: In the backend model, ScopeV2 is only a structural wrapper and has no semantic meaning.
// Each FilterV2 represents a logical scope and is mapped 1:1 to a Terraform "scope" block.
func flattenZoneV2(z *v2.ZoneV2, preferRules bool) []map[string]any {
	if z == nil {
		return nil
	}

	var allScopes []map[string]any

	for _, s := range z.Scopes {
		// ScopeV2 has no semantic meaning; it only groups filters in the backend model.
		for _, f := range s.Filters {
			allScopes = append(allScopes, flattenFilterV21(f, preferRules))
		}
	}
	return allScopes
}

// flattenFilterV21 converts a backend FilterV2 into a Terraform "scope" block.
// It handles both structured expressions and rules strings.
func flattenFilterV21(f v2.FilterV2, preferRules bool) map[string]any {
	out := map[string]any{
		SchemaIDKey:         f.ID,
		SchemaTargetTypeKey: f.ResourceType,
	}

	// When the user configured expressions, prefer them over rules.
	if !preferRules && len(f.Expressions) > 0 {
		var exps []interface{}
		for _, e := range f.Expressions {
			exps = append(exps, flattenExpressionV2(e))
		}
		out[SchemaExpressionKey] = exps
		return out
	}

	// When the user configured rules (or on import where there's no prior
	// state), write the rules string. Legacy zones without expressions
	// also land here.
	if f.Rules != "" {
		out[SchemaRulesKey] = f.Rules
	}
	return out
}

func flattenExpressionV2(e v2.ExpressionV2) map[string]any {
	m := map[string]any{
		"field":    e.Field,
		"operator": e.Operator,
		"value":    "",
		"values":   []string{},
	}

	if len(e.Values) > 0 {
		m["values"] = e.Values
	} else if e.Value != "" {
		m["value"] = e.Value
	}

	return m
}

// expandZoneV2 builds a ZoneV2 from Terraform data.
//
// NOTE: In the backend model, ScopeV2 is only a structural wrapper and has no semantic meaning.
// Each Terraform "scope" block is mapped 1:1 to a FilterV2.
func expandZoneV2(d *schema.ResourceData) *v2.ZoneV2 {
	zone := &v2.ZoneV2{
		Name:        d.Get(SchemaNameKey).(string),
		Description: d.Get(SchemaDescriptionKey).(string),
	}
	if rawScopes, ok := d.Get(SchemaScopeKey).(*schema.Set); ok {
		scope := v2.ScopeV2{}
		scopeList := rawScopes.List()
		for _, raw := range scopeList {
			scope.Filters = append(scope.Filters, expandFilterV2(raw))
		}
		if len(scope.Filters) > 0 {
			zone.Scopes = []v2.ScopeV2{scope}
		}
	}

	return zone
}

// expandFilterV2 converts the raw filter block into a v2.FilterV2.
// It handles both structured expressions and rules strings.
func expandFilterV2(raw interface{}) v2.FilterV2 {
	m := raw.(map[string]interface{})

	filter := v2.FilterV2{
		ID:           m[SchemaIDKey].(int),
		ResourceType: m[SchemaTargetTypeKey].(string),
	}

	// Check for rules string first (v2 rules syntax)
	if rules, ok := m[SchemaRulesKey].(string); ok && rules != "" {
		filter.Rules = rules
		return filter
	}

	// Otherwise, expand expressions
	if exprs, ok := m[SchemaExpressionKey].([]interface{}); ok {
		for _, e := range exprs {
			filter.Expressions = append(filter.Expressions, expandExpressionV2(e))
		}
	}

	return filter
}

// expandExpressionV2 converts the raw expression block into a v2.ExpressionV2.
func expandExpressionV2(raw interface{}) v2.ExpressionV2 {
	m := raw.(map[string]interface{})

	expr := v2.ExpressionV2{
		Field:    m["field"].(string),
		Operator: m["operator"].(string),
	}

	if vals, ok := m["values"].([]interface{}); ok && len(vals) > 0 {
		expr.Values = interfaceSliceToStrings(vals)
	} else if v, ok := m["value"].(string); ok && v != "" {
		expr.Value = v
	}

	return expr
}

// interfaceSliceToStrings converts a []interface{} (Terraform list) to []string.
func interfaceSliceToStrings(v []interface{}) []string {
	out := make([]string, len(v))
	for i, x := range v {
		out[i] = x.(string)
	}
	return out
}

package sysdig

import (
	"context"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func deprecatedResourceSysdigSecureScanningPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "The legacy scanning engine has been deprecated. This resource will be removed in future releases.",
		CreateContext:      deprecatedResourceSysdigScanningPolicyCreate,
		ReadContext:        deprecatedResourceSysdigScanningPolicyRead,
		UpdateContext:      deprecatedResourceSysdigScanningPolicyUpdate,
		DeleteContext:      deprecatedResourceSysdigScanningPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Required: true,
			},
			"isdefault": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1_0",
			},
			"policy_bundle_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"rules": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gate": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"always", "dockerfile", "files", "licenses", "metadata", "npms", "packages", "passwd_file", "retrieved_files", "vulnerabilities", "secret_scans", "ruby_gems"}, false)),
						},
						"trigger": {
							Type:     schema.TypeString,
							Required: true,
							// ValidateDiagFunc: TODO: create inline func to validate each trigger options depending on gate https://developer.hashicorp.com/terraform/plugin/sdkv2/schemas/schema-behaviors
						},
						"action": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"WARN", "STOP"}, false)),
						},
						"params": {
							Type:     schema.TypeSet,
							Required: true,
							// ValidateDiagFunc: TODO: function to validate name is valid for the given trigger,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
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

func getDeprecatedSecureScanningPolicyClient(c SysdigClients) (v2.DeprecatedScanningPolicyInterface, error) {
	return c.sysdigSecureClientV2()
}

func deprecatedResourceSysdigScanningPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicy := deprecatedScanningPolicyFromResourceData(d)
	scanningPolicy, err = client.CreateDeprecatedScanningPolicy(ctx, scanningPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	deprecatedScanningPolicyToResourceData(&scanningPolicy, d)

	return nil
}

func deprecatedResourceSysdigScanningPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicy := deprecatedScanningPolicyFromResourceData(d)
	id := d.Get("id").(string)
	scanningPolicy.ID = id
	_, err = client.UpdateDeprecatedScanningPolicyByID(ctx, scanningPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedResourceSysdigScanningPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	scanningPolicy, err := client.GetDeprecatedScanningPolicyByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	deprecatedScanningPolicyToResourceData(&scanningPolicy, d)

	return nil
}

func deprecatedResourceSysdigScanningPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	err = client.DeleteDeprecatedScanningPolicyByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedScanningPolicyToResourceData(scanningPolicy *v2.DeprecatedScanningPolicy, d *schema.ResourceData) {
	d.SetId(scanningPolicy.ID)
	_ = d.Set("name", scanningPolicy.Name)
	_ = d.Set("version", scanningPolicy.Version)
	_ = d.Set("comment", scanningPolicy.Comment)
	_ = d.Set("isdefault", scanningPolicy.IsDefault)
	_ = d.Set("policy_bundle_id", scanningPolicy.PolicyBundleId)

	var rules []map[string]interface{}
	for _, rule := range scanningPolicy.Rules {
		ruleInfo := deprecatedScanningPolicyRulesToResourceData(rule)

		rules = append(rules, ruleInfo)
	}

	_ = d.Set("rules", rules)
}

func deprecatedScanningPolicyRulesToResourceData(scanningPolicyRule v2.DeprecatedScanningGate) map[string]interface{} {
	rule := map[string]interface{}{
		"id":      scanningPolicyRule.ID,
		"gate":    scanningPolicyRule.Gate,
		"trigger": scanningPolicyRule.Trigger,
		"action":  scanningPolicyRule.Action,
	}

	var params []map[string]interface{}
	for _, param := range scanningPolicyRule.Params {
		params = append(params, map[string]interface{}{
			"name":  param.Name,
			"value": param.Value,
		})
	}
	rule["params"] = params

	return rule
}

func deprecatedScanningPolicyFromResourceData(d *schema.ResourceData) v2.DeprecatedScanningPolicy {
	scanningPolicy := v2.DeprecatedScanningPolicy{
		Name:           d.Get("name").(string),
		ID:             d.Get("id").(string),
		Comment:        d.Get("comment").(string),
		Version:        d.Get("version").(string),
		IsDefault:      d.Get("isdefault").(bool),
		PolicyBundleId: d.Get("policy_bundle_id").(string),
	}
	scanningPolicy.Rules = deprecatedScanningPolicyRulesFromResourceData(d)

	return scanningPolicy
}

func deprecatedScanningPolicyRulesFromResourceData(d *schema.ResourceData) (rules []v2.DeprecatedScanningGate) {
	for _, ruleItr := range d.Get("rules").(*schema.Set).List() {
		ruleInfo := ruleItr.(map[string]interface{})
		rule := v2.DeprecatedScanningGate{
			Gate:    ruleInfo["gate"].(string),
			ID:      ruleInfo["id"].(string),
			Trigger: ruleInfo["trigger"].(string),
			Action:  ruleInfo["action"].(string),
		}
		var params []v2.DeprecatedScanningGateParam
		for _, paramsItr := range ruleInfo["params"].(*schema.Set).List() {
			paramsInfo := paramsItr.(map[string]interface{})
			param := v2.DeprecatedScanningGateParam{
				Name:  paramsInfo["name"].(string),
				Value: paramsInfo["value"].(string),
			}
			params = append(params, param)
		}
		rule.Params = params
		rules = append(rules, rule)
	}
	return rules
}

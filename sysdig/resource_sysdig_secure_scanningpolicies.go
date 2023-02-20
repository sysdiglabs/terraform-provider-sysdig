package sysdig

import (
	"context"
	"time"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureScanningPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigScanningPolicyCreate,
		ReadContext:   resourceSysdigScanningPolicyRead,
		UpdateContext: resourceSysdigScanningPolicyUpdate,
		DeleteContext: resourceSysdigScanningPolicyDelete,
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

func resourceSysdigScanningPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicy := scanningPolicyFromResourceData(d)
	scanningPolicy, err = client.CreateScanningPolicy(ctx, scanningPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyToResourceData(&scanningPolicy, d)

	return nil
}

func resourceSysdigScanningPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicy := scanningPolicyFromResourceData(d)
	id := d.Get("id").(string)
	scanningPolicy.ID = id
	_, err = client.UpdateScanningPolicyById(ctx, scanningPolicy)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigScanningPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	scanningPolicy, err := client.GetScanningPolicyById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyToResourceData(&scanningPolicy, d)

	return nil
}

func resourceSysdigScanningPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Get("id").(string)
	err = client.DeleteScanningPolicyById(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func scanningPolicyToResourceData(scanningPolicy *secure.ScanningPolicy, d *schema.ResourceData) {

	d.SetId(scanningPolicy.ID)
	_ = d.Set("name", scanningPolicy.Name)
	_ = d.Set("version", scanningPolicy.Version)
	_ = d.Set("comment", scanningPolicy.Comment)
	_ = d.Set("isdefault", scanningPolicy.IsDefault)
	_ = d.Set("policy_bundle_id", scanningPolicy.PolicyBundleId)

	var rules []map[string]interface{}
	for _, rule := range scanningPolicy.Rules {
		ruleInfo := scanningPolicyRulesToResourceData(rule)

		rules = append(rules, ruleInfo)
	}

	_ = d.Set("rules", rules)

}

func scanningPolicyRulesToResourceData(scanningPolicyRule secure.ScanningGate) map[string]interface{} {
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

func scanningPolicyFromResourceData(d *schema.ResourceData) secure.ScanningPolicy {
	scanningPolicy := secure.ScanningPolicy{
		Name:           d.Get("name").(string),
		ID:             d.Get("id").(string),
		Comment:        d.Get("comment").(string),
		Version:        d.Get("version").(string),
		IsDefault:      d.Get("isdefault").(bool),
		PolicyBundleId: d.Get("policy_bundle_id").(string),
	}
	scanningPolicy.Rules = scanningPolicyRulesFromResourceData(d)

	return scanningPolicy
}

func scanningPolicyRulesFromResourceData(d *schema.ResourceData) (rules []secure.ScanningGate) {
	for _, ruleItr := range d.Get("rules").(*schema.Set).List() {
		ruleInfo := ruleItr.(map[string]interface{})
		rule := secure.ScanningGate{
			Gate:    ruleInfo["gate"].(string),
			ID:      ruleInfo["id"].(string),
			Trigger: ruleInfo["trigger"].(string),
			Action:  ruleInfo["action"].(string),
		}
		var params []secure.ScanningGateParam
		for _, paramsItr := range ruleInfo["params"].(*schema.Set).List() {
			paramsInfo := paramsItr.(map[string]interface{})
			param := secure.ScanningGateParam{
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

package sysdig

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureManagedRuleset() *schema.Resource {
	timeout := 5 * time.Minute
	return &schema.Resource{
		CreateContext: resourceSysdigManagedRulesetCreate,
		ReadContext:   resourceSysdigManagedRulesetRead,
		UpdateContext: resourceSysdigManagedRulesetUpdate,
		DeleteContext: resourceSysdigManagedRulesetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSysdigSecureManagedRulesetImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: createPolicySchema(map[string]*schema.Schema{
			"inherited_from": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"template_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"severity": {
				Type:             schema.TypeInt,
				Default:          4,
				Optional:         true,
				ValidateDiagFunc: validateDiagFunc(validation.IntBetween(0, 7)),
			},
			"disabled_rules": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func resourceSysdigManagedRulesetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("inherited_from.0.name").(string)
	policyType := d.Get("inherited_from.0.type").(string)

	managedPolicy, err := getManagedPolicy(ctx, client, policyName, policyType)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := v2.Policy{}

	policy.Rules = managedPolicy.Rules
	updateManagedRulesetFromResourceData(&policy, d)
	policy.TemplateId = managedPolicy.TemplateId
	policy.TemplateVersion = managedPolicy.TemplateVersion

	createdPolicy, err := client.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	managedRulesetToResourceData(&createdPolicy, d)

	return nil
}

func managedRulesetToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	commonPolicyToResourceData(policy, d)

	_ = d.Set("description", policy.Description)
	_ = d.Set("severity", policy.Severity)
	_ = d.Set("template_id", policy.TemplateId)

	disabledRules := []string{}
	for _, rule := range policy.Rules {
		if !rule.Enabled {
			disabledRules = append(disabledRules, rule.Name)
		}
	}
	_ = d.Set("disabled_rules", disabledRules)
}

func updateManagedRulesetFromResourceData(policy *v2.Policy, d *schema.ResourceData) {
	commonPolicyFromResourceData(policy, d)
	policy.Description = d.Get("description").(string)
	policy.Severity = d.Get("severity").(int)
	policy.TemplateId = d.Get("template_id").(int)

	disabledRules := d.Get("disabled_rules").(*schema.Set)
	for _, rule := range policy.Rules {
		if disabledRules.Contains(rule.Name) {
			rule.Enabled = false
		} else {
			rule.Enabled = true
		}
	}
}

func resourceSysdigManagedRulesetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	policy, statusCode, err := client.GetPolicyByID(ctx, id)

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	managedRulesetToResourceData(&policy, d)

	return nil
}

func resourceSysdigManagedRulesetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeletePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigManagedRulesetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	policy, statusCode, err := client.GetPolicyByID(ctx, id)

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	updateManagedRulesetFromResourceData(&policy, d)

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSysdigSecureManagedRulesetImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}

	managedRuleset, _, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if managedRuleset.TemplateId == 0 || managedRuleset.IsDefault {
		return nil, errors.New("unable to import policy that is not a managed ruleset")
	}

	policies, _, err := client.GetPolicies(ctx)
	if err != nil {
		return nil, err
	}

	for _, policy := range policies {
		if policy.IsDefault && policy.TemplateId == managedRuleset.TemplateId {
			inheritedFrom := map[string]string{
				"name": policy.Name,
				"type": policy.Type,
			}
			_ = d.Set("inherited_from", []map[string]string{inheritedFrom})

			break
		}
	}

	return []*schema.ResourceData{d}, nil
}

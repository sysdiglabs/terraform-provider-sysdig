package sysdig

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCustomPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigCustomPolicyCreate,
		ReadContext:   resourceSysdigCustomPolicyRead,
		UpdateContext: resourceSysdigCustomPolicyUpdate,
		DeleteContext: resourceSysdigCustomPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSysdigSecureCustomPolicyImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: createPolicySchema(map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "falco",
				ValidateDiagFunc: validateDiagFunc(validatePolicyType),
			},
			"severity": {
				Type:             schema.TypeInt,
				Default:          4,
				Optional:         true,
				ValidateDiagFunc: validateDiagFunc(validation.IntBetween(0, 7)),
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		}),
	}
}

func resourceSysdigCustomPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := customPolicyFromResourceData(d)
	policy, err = client.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	customPolicyToResourceData(&policy, d)

	return nil
}

func customPolicyFromResourceData(d *schema.ResourceData) v2.Policy {
	policy := &v2.Policy{}
	commonPolicyFromResourceData(policy, d)

	policy.Description = d.Get("description").(string)
	policy.Severity = d.Get("severity").(int)
	policy.Type = d.Get("type").(string)

	policy.Rules = []*v2.PolicyRule{}

	rules := d.Get("rules").([]interface{})
	for index := range rules {
		rule := &v2.PolicyRule{
			Name:    d.Get(fmt.Sprintf("rules.%d.name", index)).(string),
			Enabled: d.Get(fmt.Sprintf("rules.%d.enabled", index)).(bool),
		}
		policy.Rules = append(policy.Rules, rule)
	}

	return *policy
}

func customPolicyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	commonPolicyToResourceData(policy, d)

	_ = d.Set("description", policy.Description)
	_ = d.Set("severity", policy.Severity)
	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")
	}

	currentRules := getPolicyRulesFromResourceData(d)
	newRules := make([]map[string]interface{}, len(policy.Rules))
	for _, rule := range policy.Rules {
		newRules = append(newRules, map[string]interface{}{
			"name":    rule.Name,
			"enabled": rule.Enabled,
		})
	}

	areRulesSame := reflect.DeepEqual(currentRules, newRules)
	if !areRulesSame {
		_ = d.Set("rules", newRules)
	} else {
		_ = d.Set("rules", currentRules)
	}
}

func getPolicyRulesFromResourceData(d *schema.ResourceData) []*v2.PolicyRule {
	rules := d.Get("rules").([]interface{})
	policyRules := make([]*v2.PolicyRule, len(rules))

	for i, rule := range rules {
		policyRules[i] = &v2.PolicyRule{
			Name:    rule.(map[string]interface{})["name"].(string),
			Enabled: rule.(map[string]interface{})["enabled"].(bool),
		}
	}

	return policyRules
}

func resourceSysdigCustomPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	customPolicyToResourceData(&policy, d)

	return nil
}

func resourceSysdigCustomPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeletePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigCustomPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := customPolicyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigSecureCustomPolicyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}

	policy, _, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if policy.IsDefault || policy.TemplateId != 0 {
		return nil, errors.New("unable to import policy that is not a custom policy")
	}

	return []*schema.ResourceData{d}, nil
}

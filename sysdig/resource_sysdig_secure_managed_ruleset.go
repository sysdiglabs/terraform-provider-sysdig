package sysdig

import (
	"context"
	"net/http"
	"strconv"
	"strings"
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
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"notification_channels": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"runbook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"actions": policyActionBlockSchema,
		},
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

	updateManagedRulesetFromResourceData(&policy, d)
	policy.TemplateId = managedPolicy.TemplateId
	policy.TemplateVersion = managedPolicy.TemplateVersion
	policy.Rules = managedPolicy.Rules

	createdPolicy, err := client.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	managedRulesetToResourceData(&createdPolicy, d)

	return nil
}

func managedRulesetToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	if policy.ID != 0 {
		d.SetId(strconv.Itoa(policy.ID))
	}

	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("enabled", policy.Enabled)
	_ = d.Set("severity", policy.Severity)
	_ = d.Set("scope", policy.Scope)
	_ = d.Set("version", policy.Version)
	_ = d.Set("notification_channels", policy.NotificationChannelIds)
	_ = d.Set("runbook", policy.Runbook)
	_ = d.Set("template_id", policy.TemplateId)

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
			}}
		}
	}

	currentContainerAction := d.Get("actions.0.container").(string)
	currentCaptureAction := d.Get("actions.0.capture").([]interface{})
	// If the policy retrieved from service has no actions and the current state is default values,
	// then do not set the "actions" key as it may cause terraform to think there has been a state change
	if len(policy.Actions) > 0 || currentContainerAction != "" || len(currentCaptureAction) > 0 {
		_ = d.Set("actions", actions)
	}

	disabledRules := []string{}
	for _, rule := range policy.Rules {
		if !rule.Enabled {
			disabledRules = append(disabledRules, rule.Name)
		}
	}
	_ = d.Set("disabled_rules", disabledRules)
}

func updateManagedRulesetFromResourceData(policy *v2.Policy, d *schema.ResourceData) {
	policy.Name = d.Get("name").(string)
	policy.Description = d.Get("description").(string)
	policy.Enabled = d.Get("enabled").(bool)
	policy.Runbook = d.Get("runbook").(string)
	policy.Severity = d.Get("severity").(int)
	policy.Scope = d.Get("scope").(string)
	policy.TemplateId = d.Get("template_id").(int)

	addActionsToPolicy(d, policy)

	disabledRules := d.Get("disabled_rules").(*schema.Set)
	for _, rule := range policy.Rules {
		if disabledRules.Contains(rule.Name) {
			rule.Enabled = false
		} else {
			rule.Enabled = true
		}
	}

	policy.NotificationChannelIds = []int{}
	notificationChannelIdSet := d.Get("notification_channels").(*schema.Set)
	for _, id := range notificationChannelIdSet.List() {
		policy.NotificationChannelIds = append(policy.NotificationChannelIds, id.(int))
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
		d.SetId("")
		if statusCode == http.StatusNotFound {
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
		d.SetId("")
		if statusCode == http.StatusNotFound {
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

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
)

func resourceSysdigSecureManagedPolicy() *schema.Resource {
	timeout := 5 * time.Minute
	return &schema.Resource{
		CreateContext: resourceSysdigManagedPolicyCreate,
		ReadContext:   resourceSysdigManagedPolicyRead,
		UpdateContext: resourceSysdigManagedPolicyUpdate,
		DeleteContext: resourceSysdigManagedPolicyDelete,
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
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "falco",
				ValidateDiagFunc: validateDiagFunc(validatePolicyType),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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

func resourceSysdigManagedPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policies, _, err := client.GetPolicies(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := d.Get("type").(string)
	var policy v2.Policy
	for _, existingPolicy := range policies {
		if existingPolicy.Name == policyName && existingPolicy.Type == policyType {
			if !existingPolicy.IsDefault {
				return diag.Errorf("policy is not a managed policy - use `resource_sysdig_secure_policy`")
			}
			policy = existingPolicy
		}
	}

	if policy.ID == 0 {
		return diag.Errorf("unable to find managed policy")
	}

	updateManagedPolicyFromResourceData(&policy, d)

	policy, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	managedPolicyToResourceData(&policy, d)

	return nil
}

func managedPolicyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	if policy.ID != 0 {
		d.SetId(strconv.Itoa(policy.ID))
	}

	_ = d.Set("name", policy.Name)
	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")
	}
	_ = d.Set("enabled", policy.Enabled)
	_ = d.Set("scope", policy.Scope)
	_ = d.Set("version", policy.Version)
	_ = d.Set("notification_channels", policy.NotificationChannelIds)
	_ = d.Set("runbook", policy.Runbook)

	actions := []map[string]interface{}{{}}
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			actions[0]["container"] = strings.ToLower(action)
			//d.Set("actions.0.container", strings.ToLower(action))
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

func updateManagedPolicyFromResourceData(policy *v2.Policy, d *schema.ResourceData) {
	policy.Enabled = d.Get("enabled").(bool)
	policy.Runbook = d.Get("runbook").(string)
	policy.Scope = d.Get("scope").(string)

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

func resourceSysdigManagedPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	managedPolicyToResourceData(&policy, d)

	return nil
}

func resourceSysdigManagedPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	// Reset everything back to default values for managed policy
	policy, statusCode, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		d.SetId("")
		if statusCode == http.StatusNotFound {
			return diag.FromErr(err)
		}
	}

	// Disable the policy as the managed policy is no longer going to be managed by Terraform
	policy.Enabled = false
	policy.Runbook = ""
	policy.Scope = ""
	policy.Actions = []v2.Action{}
	policy.NotificationChannelIds = []int{}
	for _, rule := range policy.Rules {
		rule.Enabled = true
	}

	policy, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigManagedPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updateManagedPolicyFromResourceData(&policy, d)

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

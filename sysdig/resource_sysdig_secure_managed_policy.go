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
)

func resourceSysdigSecureManagedPolicy() *schema.Resource {
	timeout := 5 * time.Minute
	return &schema.Resource{
		CreateContext: resourceSysdigManagedPolicyCreate,
		ReadContext:   resourceSysdigManagedPolicyRead,
		UpdateContext: resourceSysdigManagedPolicyUpdate,
		DeleteContext: resourceSysdigManagedPolicyDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: createPolicySchema(map[string]*schema.Schema{
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "falco",
				ValidateDiagFunc: validateDiagFunc(validatePolicyType),
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

func resourceSysdigManagedPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policyName := d.Get("name").(string)
	policyType := d.Get("type").(string)

	policy, err := getManagedPolicy(ctx, client, policyName, policyType)
	if err != nil {
		return diag.FromErr(err)
	}

	updateManagedPolicyFromResourceData(policy, d)

	updatedPolicy, err := client.UpdatePolicy(ctx, *policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	managedPolicyToResourceData(&updatedPolicy, d)

	return nil
}

func managedPolicyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	commonPolicyToResourceData(policy, d)

	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")
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
	commonPolicyFromResourceData(policy, d)

	disabledRules := d.Get("disabled_rules").(*schema.Set)
	for _, rule := range policy.Rules {
		if disabledRules.Contains(rule.Name) {
			rule.Enabled = false
		} else {
			rule.Enabled = true
		}
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
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	managedPolicyToResourceData(&policy, d)

	return nil
}

func resourceSysdigManagedPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	// Reset everything back to default values for managed policy
	policy, statusCode, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
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
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigManagedPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
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

	updateManagedPolicyFromResourceData(&policy, d)

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func getManagedPolicy(ctx context.Context, client v2.PolicyInterface, policyName string, policyType string) (*v2.Policy, error) {
	policies, _, err := client.GetPolicies(ctx)
	if err != nil {
		return nil, err
	}

	var policy v2.Policy
	for _, existingPolicy := range policies {
		if existingPolicy.Name == policyName && existingPolicy.Type == policyType {
			if !existingPolicy.IsDefault {
				return nil, errors.New("policy is not a managed policy")
			}
			policy = existingPolicy
		}
	}

	if policy.ID != 0 {
		return &policy, nil
	}

	return nil, errors.New("unable to find managed policy")
}

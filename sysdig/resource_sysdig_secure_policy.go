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

var validatePolicyType = validation.StringInSlice([]string{"falco", "list_matching", "k8s_audit", "aws_cloudtrail", "gcp_auditlog", "azure_platformlogs"}, false)
var policyActionBlockSchema = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"container": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"stop", "pause", "kill"}, false),
			},
			"capture": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"seconds_after_event": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(0)),
						},
						"seconds_before_event": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(0)),
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	},
}

func resourceSysdigSecurePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigPolicyCreate,
		ReadContext:   resourceSysdigPolicyRead,
		UpdateContext: resourceSysdigPolicyUpdate,
		DeleteContext: resourceSysdigPolicyDelete,

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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"rule_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func getSecurePolicyClient(c SysdigClients) (v2.PolicyInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policy := policyFromResourceData(d)
	policy, err = client.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	policyToResourceData(&policy, d)

	return nil
}

func policyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	if policy.ID != 0 {
		d.SetId(strconv.Itoa(policy.ID))
	}

	_ = d.Set("name", policy.Name)
	_ = d.Set("description", policy.Description)
	_ = d.Set("scope", policy.Scope)
	_ = d.Set("enabled", policy.Enabled)
	_ = d.Set("version", policy.Version)
	_ = d.Set("severity", policy.Severity)
	_ = d.Set("runbook", policy.Runbook)
	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")

	}

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

	_ = d.Set("notification_channels", policy.NotificationChannelIds)

	_ = d.Set("rule_names", policy.RuleNames)

}

func policyFromResourceData(d *schema.ResourceData) v2.Policy {
	policy := v2.Policy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Severity:    d.Get("severity").(int),
		Enabled:     d.Get("enabled").(bool),
		Type:        d.Get("type").(string),
		Runbook:     d.Get("runbook").(string),
	}

	scope := d.Get("scope").(string)
	if scope != "" {
		policy.Scope = scope
	}

	addActionsToPolicy(d, &policy)

	policy.RuleNames = []string{}
	rule_names := d.Get("rule_names").(*schema.Set)
	for _, name := range rule_names.List() {
		if rule_name, ok := name.(string); ok {
			rule_name = strings.TrimSpace(rule_name)
			policy.RuleNames = append(policy.RuleNames, rule_name)
		}
	}

	policy.NotificationChannelIds = []int{}
	notificationChannelIdSet := d.Get("notification_channels").(*schema.Set)
	for _, id := range notificationChannelIdSet.List() {
		policy.NotificationChannelIds = append(policy.NotificationChannelIds, id.(int))
	}

	return policy
}

func addActionsToPolicy(d *schema.ResourceData, policy *v2.Policy) {
	policy.Actions = []v2.Action{}
	actions := d.Get("actions").([]interface{})
	if len(actions) == 0 {
		return
	}

	containerAction := d.Get("actions.0.container").(string)
	if containerAction != "" {
		containerAction = strings.ToUpper("POLICY_ACTION_" + containerAction)

		policy.Actions = append(policy.Actions, v2.Action{Type: containerAction})
	}

	if captureAction := d.Get("actions.0.capture").([]interface{}); len(captureAction) > 0 {
		afterEventNs := d.Get("actions.0.capture.0.seconds_after_event").(int) * 1000000000
		beforeEventNs := d.Get("actions.0.capture.0.seconds_before_event").(int) * 1000000000
		name := d.Get("actions.0.capture.0.name").(string)
		policy.Actions = append(policy.Actions, v2.Action{
			Type:                 "POLICY_ACTION_CAPTURE",
			IsLimitedToContainer: false,
			AfterEventNs:         afterEventNs,
			BeforeEventNs:        beforeEventNs,
			Name:                 name,
		})
	}
}

func resourceSysdigPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	policyToResourceData(&policy, d)

	return nil
}

func resourceSysdigPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceSysdigPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	policy := policyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

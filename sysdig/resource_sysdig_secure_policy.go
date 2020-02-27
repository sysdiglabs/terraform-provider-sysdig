package sysdig

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

var defaultMatchActions = map[string]string{
	"accept": "DEFAULT_MATCH_EFFECT_ACCEPT",
	"deny":   "DEFAULT_MATCH_EFFECT_DENY",
	"none":   "DEFAULT_MATCH_EFFECT_NEXT",
}

var matchActions = map[string]string{
	"accept": "MATCH_EFFECT_ACCEPT",
	"deny":   "MATCH_EFFECT_DENY",
	"none":   "MATCH_EFFECT_NEXT",
}

func resourceSysdigSecurePolicy() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigPolicyCreate,
		Read:   resourceSysdigPolicyRead,
		Update: resourceSysdigPolicyUpdate,
		Delete: resourceSysdigPolicyDelete,

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
			"severity": {
				Type:         schema.TypeInt,
				Default:      4,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 4, 6, 7}),
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
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"stop", "pause"}, false),
						},
						"capture": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds_after_event": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"seconds_before_event": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntAtLeast(0),
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

func resourceSysdigPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	policy := policyFromResourceData(d)
	policy, err := client.CreatePolicy(policy)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(policy.ID))
	d.Set("version", policy.Version)

	return nil
}

func policyFromResourceData(d *schema.ResourceData) secure.Policy {
	policy := secure.Policy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Severity:    d.Get("severity").(int),
		Enabled:     d.Get("enabled").(bool),
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

func addActionsToPolicy(d *schema.ResourceData, policy *secure.Policy) {
	policy.Actions = []secure.Action{}
	actions := d.Get("actions").([]interface{})
	if len(actions) == 0 {
		return
	}

	containerAction := d.Get("actions.0.container").(string)
	if containerAction != "" {
		containerAction = strings.ToUpper("POLICY_ACTION_" + containerAction)

		policy.Actions = append(policy.Actions, secure.Action{Type: containerAction})
	}

	if captureAction := d.Get("actions.0.capture").([]interface{}); len(captureAction) > 0 {
		afterEventNs := d.Get("actions.0.capture.0.seconds_after_event").(int) * 1000000000
		beforeEventNs := d.Get("actions.0.capture.0.seconds_before_event").(int) * 1000000000
		policy.Actions = append(policy.Actions, secure.Action{
			Type:                 "POLICY_ACTION_CAPTURE",
			IsLimitedToContainer: false,
			AfterEventNs:         afterEventNs,
			BeforeEventNs:        beforeEventNs,
		})
	}
}

func resourceSysdigPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, _ := strconv.Atoi(d.Id())
	policy, err := client.GetPolicyById(id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("scope", policy.Scope)
	d.Set("enabled", policy.Enabled)
	d.Set("version", policy.Version)
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			d.Set("actions.0.container", strings.ToLower(action))
		} else {
			d.Set("actions.0.capture.seconds_after_event", action.AfterEventNs/1000000000)
			d.Set("actions.0.capture.seconds_before_event", action.BeforeEventNs/100000000)
		}
	}

	d.Set("notification_channels", policy.NotificationChannelIds)
	d.Set("rule_names", policy.RuleNames)

	return nil
}

func resourceSysdigPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	id, _ := strconv.Atoi(d.Id())

	return client.DeletePolicy(id)
}

func resourceSysdigPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*SysdigClients).sysdigSecureClient

	policy := policyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err := client.UpdatePolicy(policy)
	return err
}

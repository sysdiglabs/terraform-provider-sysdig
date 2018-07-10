package sysdig

import (
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

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
				Type:     schema.TypeInt,
				Required: true,
			},
			"container_scope": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"host_scope": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"falco_rule_name_regex": {
				Type:     schema.TypeString,
				Required: true,
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
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"seconds_after_event": &schema.Schema{
										Type:     schema.TypeInt,
										Required: true,
									},
									"seconds_before_event": &schema.Schema{
										Type:     schema.TypeInt,
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

func resourceSysdigPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	policy := policyFromResourceData(d)
	policy, err := client.CreatePolicy(policy)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(policy.ID))
	d.Set("version", policy.Version)

	return nil
}

func policyFromResourceData(d *schema.ResourceData) Policy {
	policy := Policy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Severity:       d.Get("severity").(int),
		ContainerScope: d.Get("container_scope").(bool),
		HostScope:      d.Get("host_scope").(bool),
		Enabled:        d.Get("enabled").(bool),
		Scope:          d.Get("filter").(string),
	}

	addActionsToPolicy(d, &policy)

	falcoRuleNameRegex := d.Get("falco_rule_name_regex")
	if falcoRuleNameRegex != nil {
		policy.FalcoConfiguration = FalcoConfiguration{
			RuleNameRegEx: falcoRuleNameRegex.(string),
		}
	}

	return policy
}

func addActionsToPolicy(d *schema.ResourceData, policy *Policy) {
	actions := d.Get("actions").([]interface{})
	if len(actions) == 0 {
		return
	}

	policy.Actions = []Action{}

	containerAction := d.Get("actions.0.container").(string)
	if containerAction != "" {
		containerAction = strings.ToUpper("POLICY_ACTION_" + containerAction)

		policy.Actions = append(policy.Actions, Action{Type: containerAction})
	}

	captureAction := d.Get("actions.0.capture").(map[string]interface{})
	if len(captureAction) != 0 {
		afterEventNs, _ := strconv.Atoi(d.Get("actions.0.capture.seconds_after_event").(string) + "000000000")
		beforeEventNs, _ := strconv.Atoi(d.Get("actions.0.capture.seconds_before_event").(string) + "000000000")
		policy.Actions = append(policy.Actions, Action{
			Type:                 "POLICY_ACTION_CAPTURE",
			IsLimitedToContainer: false,
			AfterEventNs:         afterEventNs,
			BeforeEventNs:        beforeEventNs,
		})
	}
}

func resourceSysdigPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())
	policy, err := client.GetPolicyById(id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("container_scope", policy.ContainerScope)
	d.Set("host_scope", policy.HostScope)
	d.Set("enabled", policy.Enabled)
	d.Set("filter", policy.Scope)
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

	if policy.FalcoConfiguration.RuleNameRegEx != "" {
		d.Set("falco_rule_name_regex", policy.FalcoConfiguration.RuleNameRegEx)
	}

	return nil
}

func resourceSysdigPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	id, _ := strconv.Atoi(d.Id())

	return client.DeletePolicy(id)
}

func resourceSysdigPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SysdigSecureClient)

	policy := policyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err := client.UpdatePolicy(policy)
	return err
}

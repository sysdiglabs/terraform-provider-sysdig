package sysdig

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

var validateStatefulRuleSource = validation.StringInSlice([]string{"awscloudtrail_stateful"}, false)

func resourceSysdigSecureStatefulRule() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleStatefulCreate,
		UpdateContext: resourceSysdigRuleStatefulUpdate,
		ReadContext:   resourceSysdigRuleStatefulRead,
		DeleteContext: resourceSysdigRuleStatefulDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"source": {
				Type:             schema.TypeString,
				Optional:         false,
				Required:         true,
				ValidateDiagFunc: validateDiagFunc(validateStatefulRuleSource),
			},
			"ruletype": {
				Type:     schema.TypeString,
				Optional: false,
				Required: true,
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
					v2.RuleTypeStatefulSequence,
					v2.StatefulCountRuleType,
					v2.StatefulUniqPercentRuleType,
				}, false)),
			},
			"append": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"exceptions": {
				Type:     schema.TypeList,
				Optional: false,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comps": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		}),
	}
}

func resourceSysdigRuleStatefulCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleStatefulFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err = client.CreateStatefulRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	d.SetId(strconv.Itoa(rule.ID))
	_ = d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleStatefulRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// for stateful rules, we'll need to get the rule group

	nameObj, ok := d.GetOk("name")
	if !ok {
		return diag.FromErr(errors.New("name is required"))
	}

	name := nameObj.(string)

	sourceObj, ok := d.GetOk("source")
	if !ok {
		return diag.FromErr(errors.New("source is required"))
	}

	source := sourceObj.(string)

	rules, err := client.GetStatefulRuleGroup(ctx, name, source)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(rules) == 0 {
		d.SetId("")
	}

	var rule v2.Rule

	for _, r := range rules {
		if r.ID == id {
			rule = r
			break
		}
	}

	if rule.Details.Append != nil && !(*(rule.Details.Append)) {
		if rule.Details.Condition == nil {
			return diag.Errorf("no condition data for a Stateful rule")
		}
	}

	_ = d.Set("name", rule.Name)
	_ = d.Set("source", rule.Details.Source)

	if rule.Details.Append != nil {
		_ = d.Set("append", *rule.Details.Append)
	}
	if err := updateResourceDataExceptions(d, rule.Details.Exceptions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigRuleStatefulUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleStatefulFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateStatefulRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleStatefulDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteStatefulRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleStatefulFromResourceData(d *schema.ResourceData) (v2.Rule, error) {
	rule := v2.Rule{
		Name: d.Get("name").(string),
	}

	ruleType := d.Get("ruletype").(string)
	rule.Details.RuleType = ruleType

	appendMode, appendModeIsSet := d.GetOk("append")
	if appendModeIsSet {
		ptr := appendMode.(bool)
		rule.Details.Append = &ptr
	}

	if source, ok := d.GetOk("source"); ok && source.(string) != "" {
		rule.Details.Source = source.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return v2.Rule{}, errors.New("source must be set when append = false")
	}

	if exceptionsField, ok := d.GetOk("exceptions"); ok {
		StatefulExceptions := []*v2.Exception{}
		for _, exception := range exceptionsField.([]interface{}) {
			exceptionMap := exception.(map[string]interface{})
			newStatefulException := &v2.Exception{
				Name: exceptionMap["name"].(string),
			}

			fields := cast.ToStringSlice(exceptionMap["fields"])
			if len(fields) >= 1 {
				newStatefulException.Fields = fields
			}

			comps := cast.ToStringSlice(exceptionMap["comps"])
			if len(comps) >= 1 {
				newStatefulException.Comps = comps
			}

			values := cast.ToString(exceptionMap["values"])
			err := json.Unmarshal([]byte(values), &newStatefulException.Values)
			if err != nil {
				return v2.Rule{}, err
			}

			StatefulExceptions = append(StatefulExceptions, newStatefulException)
		}
		rule.Details.Exceptions = StatefulExceptions
	}

	return rule, nil
}

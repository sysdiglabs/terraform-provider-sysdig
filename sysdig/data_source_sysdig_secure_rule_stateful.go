package sysdig

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleStateful() *schema.Resource {
	timeout := 1 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleStatefulRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"source": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateDiagFunc(validateStatefulRuleSource),
			},
			"ruletype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"append": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"exceptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSysdigRuleStatefulRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

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

	ruleIndexObj, ok := d.GetOk("index")
	ruleIndex := 0
	if ok {
		ruleIndex = ruleIndexObj.(int)
	}

	rule := rules[ruleIndex]

	if len(rules) == 0 {
		d.SetId("")
	} else {
		d.SetId(strconv.Itoa(rule.ID))
	}

	_ = d.Set("name", rule.Name)
	_ = d.Set("source", source)

	if rule.Details.Append != nil {
		_ = d.Set("append", *rule.Details.Append)
	}

	exceptions := make([]any, 0, len(rule.Details.Exceptions))
	for _, exception := range rule.Details.Exceptions {
		if exception == nil {
			return diag.Errorf("exception is nil")
		}
		valuesData, err := json.Marshal(exception.Values)
		if err != nil {
			return diag.Errorf("error marshalling exception values '%+v': %s", exception.Values, err)
		}

		exceptions = append(exceptions, map[string]any{
			"name":   exception.Name,
			"values": string(valuesData),
		})
	}

	if err := d.Set("exceptions", exceptions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

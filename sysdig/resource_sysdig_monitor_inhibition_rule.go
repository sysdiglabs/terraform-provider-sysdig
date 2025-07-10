package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigMonitorInhibitionRule() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorInhibitionRuleCreate,
		UpdateContext: resourceSysdigMonitorInhibitionRuleUpdate,
		ReadContext:   resourceSysdigMonitorInhibitionRuleRead,
		DeleteContext: resourceSysdigMonitorInhibitionRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"source_matchers": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"EQUALS", "NOT_EQUALS", "REGEXP_MATCHES", "NOT_REGEXP_MATCHES"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"target_matchers": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"EQUALS", "NOT_EQUALS", "REGEXP_MATCHES", "NOT_REGEXP_MATCHES"}, false),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"equal": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getMonitorInhibitionRuleClient(c SysdigClients) (v2.InhibitionRuleInterface, error) {
	var client v2.InhibitionRuleInterface
	var err error
	switch c.GetClientType() {
	case IBMMonitor:
		client, err = c.ibmMonitorClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigMonitorClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigMonitorInhibitionRuleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorInhibitionRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	inhibitionRule, err := monitorInhibitionRuleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	inhibitionRule, err = client.CreateInhibitionRule(ctx, inhibitionRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(inhibitionRule.ID))

	return resourceSysdigMonitorInhibitionRuleRead(ctx, d, meta)
}

func resourceSysdigMonitorInhibitionRuleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorInhibitionRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	inhibitionRule, err := client.GetInhibitionRuleByID(ctx, id)
	if err != nil {
		if err == v2.ErrAlertV2NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = monitorInhibitionRuleToResourceData(inhibitionRule, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorInhibitionRuleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorInhibitionRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	inhibitionRule, err := monitorInhibitionRuleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	inhibitionRule.Version = d.Get("version").(int)
	inhibitionRule.ID, err = strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateInhibitionRule(ctx, inhibitionRule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorInhibitionRuleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getMonitorInhibitionRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteInhibitionRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func monitorInhibitionRuleFromResourceData(d *schema.ResourceData) (v2.InhibitionRule, error) {
	inhibitionRule := v2.InhibitionRule{}

	inhibitionRule.Name = d.Get("name").(string)
	inhibitionRule.Description = d.Get("description").(string)
	inhibitionRule.Enabled = d.Get("enabled").(bool)

	for _, sourceMatcher := range d.Get("source_matchers").([]any) {
		sourceMatcherMap := sourceMatcher.(map[string]any)
		labelName := sourceMatcherMap["label_name"].(string)
		operator := sourceMatcherMap["operator"].(string)
		value := sourceMatcherMap["value"].(string)
		inhibitionRule.SourceMatchers = append(inhibitionRule.SourceMatchers, v2.LabelMatchers{
			LabelName: labelName,
			Operator:  operator,
			Value:     value,
		})
	}

	for _, targetMatcher := range d.Get("target_matchers").([]any) {
		targetMatcherMap := targetMatcher.(map[string]any)
		labelName := targetMatcherMap["label_name"].(string)
		operator := targetMatcherMap["operator"].(string)
		value := targetMatcherMap["value"].(string)
		inhibitionRule.TargetMatchers = append(inhibitionRule.TargetMatchers, v2.LabelMatchers{
			LabelName: labelName,
			Operator:  operator,
			Value:     value,
		})
	}

	for _, equalItemRaw := range d.Get("equal").([]any) {
		if equalItem, ok := equalItemRaw.(string); ok {
			inhibitionRule.Equal = append(inhibitionRule.Equal, equalItem)
		}
	}

	return inhibitionRule, nil
}

func monitorInhibitionRuleToResourceData(inhibitionRule v2.InhibitionRule, d *schema.ResourceData) error {
	_ = d.Set("version", inhibitionRule.Version)
	_ = d.Set("name", inhibitionRule.Name)
	_ = d.Set("description", inhibitionRule.Description)
	_ = d.Set("enabled", inhibitionRule.Enabled)

	var sourceMatchers []any
	for _, m := range inhibitionRule.SourceMatchers {
		sourceMatchers = append(sourceMatchers, map[string]any{
			"label_name": m.LabelName,
			"operator":   m.Operator,
			"value":      m.Value,
		})
	}
	_ = d.Set("source_matchers", sourceMatchers)

	var targetMatchers []any
	for _, m := range inhibitionRule.TargetMatchers {
		targetMatchers = append(targetMatchers, map[string]any{
			"label_name": m.LabelName,
			"operator":   m.Operator,
			"value":      m.Value,
		})
	}
	_ = d.Set("target_matchers", targetMatchers)

	_ = d.Set("equal", inhibitionRule.Equal)

	return nil
}

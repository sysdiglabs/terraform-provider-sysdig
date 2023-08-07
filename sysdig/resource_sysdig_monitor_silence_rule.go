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

func resourceSysdigMonitorSilenceRule() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorSilenceRuleCreate,
		UpdateContext: resourceSysdigMonitorSilenceRuleUpdate,
		ReadContext:   resourceSysdigMonitorSilenceRuleRead,
		DeleteContext: resourceSysdigMonitorSilenceRuleDelete,
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
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"start_ts": {
				Type:     schema.TypeString,
				Required: true,
			},
			"duration_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(60),
			},
			"alert_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_channel_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func getMonitorSilenceRuleClient(c SysdigClients) (v2.SilenceRuleInterface, error) {
	var client v2.SilenceRuleInterface
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

func resourceSysdigMonitorSilenceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorSilenceRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	silenceRule, err := monitorSilenceRuleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	silenceRule, err = client.CreateSilenceRule(ctx, silenceRule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(silenceRule.ID))

	return resourceSysdigMonitorSilenceRuleRead(ctx, d, meta)
}

func resourceSysdigMonitorSilenceRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorSilenceRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	silenceRule, err := client.GetSilenceRule(ctx, id)

	if err != nil {
		if err == v2.SilenceRuleNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// suppress diff of "enabled" field if the silence interval is over: it will always be false from the api
	// any update of an ended silence rule results in a 422 Unprocessable Entity error from the api
	silenceRuleEnd := time.Unix(silenceRule.StartTs/1000+int64(silenceRule.DurationInSec), 0)
	if time.Now().After(silenceRuleEnd) {
		silenceRule.Enabled = d.Get("enabled").(bool)
	}

	err = monitorSilenceRuleToResourceData(silenceRule, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorSilenceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorSilenceRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	silenceRule, err := monitorSilenceRuleFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	silenceRule.Version = d.Get("version").(int)
	silenceRule.ID, err = strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateSilenceRule(ctx, silenceRule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorSilenceRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorSilenceRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteSilenceRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func monitorSilenceRuleFromResourceData(d *schema.ResourceData) (v2.SilenceRule, error) {
	silenceRule := v2.SilenceRule{}

	silenceRule.Name = d.Get("name").(string)
	silenceRule.Enabled = d.Get("enabled").(bool)
	startTs, err := strconv.ParseInt(d.Get("start_ts").(string), 10, 64)
	if err != nil {
		return silenceRule, err
	}
	silenceRule.StartTs = startTs
	silenceRule.DurationInSec = d.Get("duration_seconds").(int)
	silenceRule.Scope = d.Get("scope").(string)
	alertIds := d.Get("alert_ids").(*schema.Set)
	for _, rawAlertId := range alertIds.List() {
		if alertId, ok := rawAlertId.(int); ok {
			silenceRule.AlertIds = append(silenceRule.AlertIds, alertId)
		}
	}
	notificationChannelIds := d.Get("notification_channel_ids").(*schema.Set)
	for _, rawNotificationChannelId := range notificationChannelIds.List() {
		if notificationChannelId, ok := rawNotificationChannelId.(int); ok {
			silenceRule.NotificationChannelIds = append(silenceRule.NotificationChannelIds, notificationChannelId)
		}
	}
	return silenceRule, nil
}

func monitorSilenceRuleToResourceData(silenceRule v2.SilenceRule, d *schema.ResourceData) (err error) {
	_ = d.Set("name", silenceRule.Name)
	_ = d.Set("enabled", silenceRule.Enabled)
	_ = d.Set("start_ts", strconv.FormatInt(silenceRule.StartTs, 10))
	_ = d.Set("duration_seconds", silenceRule.DurationInSec)
	_ = d.Set("alert_ids", silenceRule.AlertIds)
	_ = d.Set("scope", silenceRule.Scope)
	_ = d.Set("notification_channel_ids", silenceRule.NotificationChannelIds)
	_ = d.Set("version", silenceRule.Version)
	return nil
}

package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigAgentAccessKey() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext:   resourceSysdigAgentAccessKeyRead,
		CreateContext: resourceSysdigAgentAccessKeyCreate,
		DeleteContext: resourceSysdigAgentAccessKeyDelete,
		UpdateContext: resourceSysdigAgentAccessKeyUpdate,
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
			"reservation": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"date_disabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSysdigAgentAccessKeyDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAgentAccessKey(ctx, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAgentAccessKeyCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}
	agentAccessKey, err := agentAccessKeyFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	agentAccessKey, err = client.CreateAgentAccessKey(ctx, agentAccessKey)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(agentAccessKey.ID))
	resourceSysdigAgentAccessKeyRead(ctx, data, meta)

	return nil
}

func resourceSysdigAgentAccessKeyUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}
	agentAccessKey, err := agentAccessKeyFromResourceData(data)
	if err != nil {
		return diag.FromErr(err)
	}

	agentAccessKey, err = client.UpdateAgentAccessKey(ctx, agentAccessKey, data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(agentAccessKey.ID))

	resourceSysdigAgentAccessKeyRead(ctx, data, meta)

	return nil
}

func agentAccessKeyFromResourceData(data *schema.ResourceData) (*v2.AgentAccessKey, error) {
	metadataFromResourceData := data.Get("metadata").(map[string]interface{})
	metadata := make(map[string]string)

	for key, val := range metadataFromResourceData {
		// Convert each value to a string, using fmt.Sprintf
		strVal := fmt.Sprintf("%v", val)
		metadata[key] = strVal
	}

	var enabled bool
	if data.Get("enabled") != nil {
		enabled = data.Get("enabled").(bool)
	}

	return &v2.AgentAccessKey{
		Reservation:  data.Get("reservation").(int),
		Limit:        data.Get("limit").(int),
		TeamID:       data.Get("team_id").(int),
		Enabled:      enabled,
		DateDisabled: data.Get("date_disabled").(string),
		DateCreated:  data.Get("date_created").(string),
		Metadata:     metadata,
	}, nil
}

func resourceSysdigAgentAccessKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	agentKeyId := d.Id()

	agentAccessKey, err := client.GetAgentAccessKeyByID(ctx, agentKeyId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(agentAccessKey.ID))
	_ = d.Set("reservation", agentAccessKey.Reservation)
	_ = d.Set("limit", agentAccessKey.Limit)
	_ = d.Set("team_id", agentAccessKey.TeamID)
	_ = d.Set("metadata", agentAccessKey.Metadata)
	_ = d.Set("enabled", agentAccessKey.Enabled)
	_ = d.Set("date_created", agentAccessKey.DateCreated)
	_ = d.Set("date_disabled", agentAccessKey.DateDisabled)
	_ = d.Set("access_key", agentAccessKey.AgentAccessKey)
	return nil
}

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

func resourceSysdigMonitorNotificationChannelIBMFunction() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorNotificationChannelIBMFunctionCreate,
		UpdateContext: resourceSysdigMonitorNotificationChannelIBMFunctionUpdate,
		ReadContext:   resourceSysdigMonitorNotificationChannelIBMFunctionRead,
		DeleteContext: resourceSysdigMonitorNotificationChannelIBMFunctionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createMonitorNotificationChannelSchema(map[string]*schema.Schema{
			"ibm_function_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"WEB_ACTION", "CLOUD_FUNCTION"}, false),
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_data": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"iam_api_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"whisk_auth_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceSysdigMonitorNotificationChannelIBMFunctionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelIBMFunctionFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigMonitorNotificationChannelIBMFunctionRead(ctx, d, meta)
}

func resourceSysdigMonitorNotificationChannelIBMFunctionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	nc, err := client.GetNotificationChannelById(ctx, id)
	if err != nil {
		if err == v2.NotificationChannelNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = secureNotificationChannelIBMFunctionToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelIBMFunctionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelIBMFunctionFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	nc.Version = d.Get("version").(int)
	nc.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateNotificationChannel(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorNotificationChannelIBMFunctionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteNotificationChannel(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func secureNotificationChannelIBMFunctionFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_IBM_FUNCTION
	nc.Options.IbmFunctionType = d.Get("ibm_function_type").(string)
	nc.Options.Url = d.Get("url").(string)
	nc.Options.CustomData = d.Get("custom_data").(map[string]interface{})
	if nc.Options.IbmFunctionType == "CLOUD_FUNCTION" {
		nc.Options.APIKey = d.Get("iam_api_key").(string)
	} else {
		nc.Options.APIKey = ""
	}
	if nc.Options.IbmFunctionType == "WEB_ACTION" {
		nc.Options.AdditionalHeaders = map[string]interface{}{
			"X-Require-Whisk-Auth": d.Get("whisk_auth_token").(string),
		}
	} else {
		nc.Options.AdditionalHeaders = map[string]interface{}{}
	}

	return
}

func secureNotificationChannelIBMFunctionToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("ibm_function_type", nc.Options.IbmFunctionType)
	_ = d.Set("url", nc.Options.Url)
	_ = d.Set("custom_data", nc.Options.CustomData)
	_ = d.Set("iam_api_key", nc.Options.APIKey)
	if nc.Options.AdditionalHeaders != nil {
		whishAuthToken, ok := nc.Options.AdditionalHeaders["X-Require-Whisk-Auth"]
		if ok {
			_ = d.Set("whisk_auth_token", whishAuthToken)
		}
	}

	return
}

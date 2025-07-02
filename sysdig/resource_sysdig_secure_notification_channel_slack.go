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

func resourceSysdigSecureNotificationChannelSlack() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelSlackCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelSlackUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelSlackRead,
		DeleteContext: resourceSysdigSecureNotificationChannelSlackDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createSecureNotificationChannelSchema(map[string]*schema.Schema{
			"url": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"channel": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_private_channel": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"private_channel_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelSlackCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelSlackFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigSecureNotificationChannelSlackRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelSlackRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelByID(ctx, id)
	if err != nil {
		if err == v2.ErrNotificationChannelNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = secureNotificationChannelSlackToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelSlackUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelSlackFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	nc.Version = d.Get("version").(int)
	nc.ID, err = strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateNotificationChannel(ctx, nc)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigSecureNotificationChannelSlackRead(ctx, d, meta)

	return nil
}

func resourceSysdigSecureNotificationChannelSlackDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteNotificationChannel(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func secureNotificationChannelSlackFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = notificationChannelTypeSlack
	nc.Options.URL = d.Get("url").(string)
	nc.Options.Channel = d.Get("channel").(string)
	nc.Options.PrivateChannel = d.Get("is_private_channel").(bool)
	nc.Options.PrivateChannelURL = d.Get("private_channel_url").(string)

	setNotificationChannelSlackTemplateConfig(&nc, d)

	return
}

func setNotificationChannelSlackTemplateConfig(nc *v2.NotificationChannel, d *schema.ResourceData) {
	templateVersion := d.Get("template_version").(string)

	switch templateVersion {
	case "v1":
		nc.Options.TemplateConfiguration = []v2.NotificationChannelTemplateConfiguration{
			{
				TemplateKey: notificationChannelTypeSlackTemplateKeyV1,
				TemplateConfigurationSections: []v2.NotificationChannelTemplateConfigurationSection{
					{
						SectionName: notificationChannelSecureEventNotificationContentSection,
						ShouldShow:  true,
					},
				},
			},
		}
	case "v2":
		nc.Options.TemplateConfiguration = []v2.NotificationChannelTemplateConfiguration{
			{
				TemplateKey: notificationChannelTypeSlackTemplateKeyV2,
				TemplateConfigurationSections: []v2.NotificationChannelTemplateConfigurationSection{
					{
						SectionName: notificationChannelSecureEventNotificationContentSection,
						ShouldShow:  true,
					},
				},
			},
		}
	}
}

func secureNotificationChannelSlackToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("url", nc.Options.URL)
	_ = d.Set("channel", nc.Options.Channel)
	_ = d.Set("is_private_channel", nc.Options.PrivateChannel)
	_ = d.Set("private_channel_url", nc.Options.PrivateChannelURL)

	err = getTemplateVersionFromNotificationChannelSlack(nc, d)

	return
}

func getTemplateVersionFromNotificationChannelSlack(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	if len(nc.Options.TemplateConfiguration) == 0 {
		return
	}

	if len(nc.Options.TemplateConfiguration) > 1 {
		return fmt.Errorf("expected slack notification templates to have only one configuration, found %d", len(nc.Options.TemplateConfiguration))
	}

	switch nc.Options.TemplateConfiguration[0].TemplateKey {
	case notificationChannelTypeSlackTemplateKeyV2:
		_ = d.Set("template_version", "v2")
	default:
		_ = d.Set("template_version", "v1")
	}

	return nil
}

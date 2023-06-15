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

func resourceSysdigSecureNotificationChannelMSTeams() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureNotificationChannelMSTeamsCreate,
		UpdateContext: resourceSysdigSecureNotificationChannelMSTeamsUpdate,
		ReadContext:   resourceSysdigSecureNotificationChannelMSTeamsRead,
		DeleteContext: resourceSysdigSecureNotificationChannelMSTeamsDelete,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"template_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
		}),
	}
}

func resourceSysdigSecureNotificationChannelMSTeamsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err := secureNotificationChannelMSTeamsFromResourceData(d, teamID)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationChannel, err = client.CreateNotificationChannel(ctx, notificationChannel)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(notificationChannel.ID))

	return resourceSysdigSecureNotificationChannelMSTeamsRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelMSTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := client.GetNotificationChannelById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	err = secureNotificationChannelMSTeamsToResourceData(&nc, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureNotificationChannelMSTeamsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureNotificationChannelClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	teamID, err := client.CurrentTeamID(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	nc, err := secureNotificationChannelMSTeamsFromResourceData(d, teamID)
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

	return resourceSysdigSecureNotificationChannelMSTeamsRead(ctx, d, meta)
}

func resourceSysdigSecureNotificationChannelMSTeamsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func secureNotificationChannelMSTeamsFromResourceData(d *schema.ResourceData, teamID int) (nc v2.NotificationChannel, err error) {
	nc, err = secureNotificationChannelFromResourceData(d, teamID)
	if err != nil {
		return
	}

	nc.Type = NOTIFICATION_CHANNEL_TYPE_MS_TEAMS
	nc.Options.Url = d.Get("url").(string)

	setNotificationChannelMSTeamsTemplateConfig(&nc, d)

	return
}

func setNotificationChannelMSTeamsTemplateConfig(nc *v2.NotificationChannel, d *schema.ResourceData) {
	templateVersion := d.Get("template_version").(string)

	switch templateVersion {
	case "v1":
		nc.Options.TemplateConfiguration = []v2.NotificationChannelTemplateConfiguration{
			{
				TemplateKey: NOTIFICATION_CHANNEL_TYPE_MS_TEAMS_TEMPLATE_KEY_V1,
				TemplateConfigurationSections: []v2.NotificationChannelTemplateConfigurationSection{
					{
						SectionName: NOTIFICATION_CHANNEL_SECURE_EVENT_NOTIFICATION_CONTENT_SECTION,
						ShouldShow:  true,
					},
				},
			},
		}
	case "v2":
		nc.Options.TemplateConfiguration = []v2.NotificationChannelTemplateConfiguration{
			{
				TemplateKey: NOTIFICATION_CHANNEL_TYPE_MS_TEAMS_TEMPLATE_KEY_V2,
				TemplateConfigurationSections: []v2.NotificationChannelTemplateConfigurationSection{
					{
						SectionName: NOTIFICATION_CHANNEL_SECURE_EVENT_NOTIFICATION_CONTENT_SECTION,
						ShouldShow:  true,
					},
				},
			},
		}
	}
}

func secureNotificationChannelMSTeamsToResourceData(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	err = secureNotificationChannelToResourceData(nc, d)
	if err != nil {
		return
	}

	_ = d.Set("url", nc.Options.Url)

	err = getTemplateVersionFromNotificationChannelMSTeams(nc, d)

	return
}

func getTemplateVersionFromNotificationChannelMSTeams(nc *v2.NotificationChannel, d *schema.ResourceData) (err error) {
	if len(nc.Options.TemplateConfiguration) == 0 {
		return
	}

	if len(nc.Options.TemplateConfiguration) > 1 {
		return fmt.Errorf("expected ms teams notification templates to have only one configuration, found %d", len(nc.Options.TemplateConfiguration))
	}

	switch nc.Options.TemplateConfiguration[0].TemplateKey {
	case NOTIFICATION_CHANNEL_TYPE_MS_TEAMS_TEMPLATE_KEY_V2:
		_ = d.Set("template_version", "v2")
	default:
		_ = d.Set("template_version", "v1")
	}

	return nil
}

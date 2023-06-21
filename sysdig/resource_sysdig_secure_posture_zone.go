package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecurePostureZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSysdigSecurePostureZoneCreate,
		UpdateContext: resourceSysdigSecurePostureZoneUpdate,
		DeleteContext: resourceSysdigSecurePostureZoneDelete,
		ReadContext:   resourceSysdigSecurePostureZoneRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"policies": {
				Required: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getPostureZoneClient(c SysdigClients) (v2.PostureZoneInterface, error) {
	var client v2.PostureZoneInterface
	var err error
	switch c.GetClientType() {
	case IBMSecure:
		client, err = c.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigSecurePostureZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	commonClient, err := meta.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}
	u, err := commonClient.GetCurrentUser(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	policiesTmp := d.Get("policies").([]interface{})
	policies := make([]string, len(policiesTmp))
	for i, p := range policiesTmp {
		policies[i] = p.(string)
	}

	id := "0"
	if d.Id() != "" {
		id = d.Id()
	}
	req := &v2.PostureZoneRequest{
		ID:          id,
		Name:        d.Get("name").(string),
		Description: "",
		PolicyIDs:   policies,
		Scopes:      make([]v2.PostureZoneScope, 0),
		Username:    u.Email,
	}

	client, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := client.CreateOrUpdatePostureZone(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zone.ID)

	resourceSysdigSecurePostureZoneRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecurePostureZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceSysdigSecurePostureZoneCreate(ctx, d, meta)
}

func resourceSysdigSecurePostureZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := client.GetPostureZone(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("name", zone.Name)
	pids := make([]string, len(zone.Policies))
	for i, p := range zone.Policies {
		pids[i] = p.ID
	}
	_ = d.Set("policies", pids)
	return nil
}

func resourceSysdigSecurePostureZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeletePostureZone(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

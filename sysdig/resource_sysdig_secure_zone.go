package sysdig

import (
	"context"
	"fmt"
	"strconv"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSysdigZoneCreate,
		ReadContext:   resourceSysdigZoneRead,
		UpdateContext: resourceSysdigZoneUpdate,
		DeleteContext: resourceSysdigZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"author": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rules": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceSysdigZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneRequest := zoneRequestFromResourceData(d)

	createdZone, err := client.CreateZone(ctx, zoneRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Sysdig Zone: %s", err))
	}

	d.SetId(fmt.Sprintf("%d", createdZone.ID))
	return resourceSysdigZoneRead(ctx, d, m)
}

func resourceSysdigZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	zone, err := client.GetZoneById(ctx, id)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_ = d.Set("name", zone.Name)
	_ = d.Set("description", zone.Description)
	_ = d.Set("scopes", flattenZoneScopes(zone.Scopes))
	_ = d.Set("is_system", zone.IsSystem)
	_ = d.Set("author", zone.Author)
	_ = d.Set("last_modified_by", zone.LastModifiedBy)
	_ = d.Set("last_updated", zone.LastUpdated)

	return nil
}

func resourceSysdigZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneRequest := zoneRequestFromResourceData(d)

	_, err = client.UpdateZone(ctx, zoneRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Sysdig Zone: %s", err))
	}

	return resourceSysdigZoneRead(ctx, d, m)
}

func resourceSysdigZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	err = client.DeleteZone(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Sysdig Zone: %s", err))
	}

	d.SetId("")
	return nil
}

func zoneRequestFromResourceData(d *schema.ResourceData) *v2.ZoneRequest {
	return &v2.ZoneRequest{
		ID:          d.Get("id").(int),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Scopes:      expandZoneScopes(d.Get("scopes").([]interface{})),
	}
}

func expandZoneScopes(scopes []interface{}) []v2.ZoneScope {
	var zoneScopes []v2.ZoneScope
	for _, scope := range scopes {
		scopeMap := scope.(map[string]interface{})
		zoneScopes = append(zoneScopes, v2.ZoneScope{
			TargetType: scopeMap["target_type"].(string),
			Rules:      scopeMap["rules"].(string),
		})
	}
	return zoneScopes
}

func flattenZoneScopes(scopes []v2.ZoneScope) []interface{} {
	var flattenedScopes []interface{}
	for _, scope := range scopes {
		flattenedScopes = append(flattenedScopes, map[string]interface{}{
			"target_type": scope.TargetType,
			"rules":       scope.Rules,
		})
	}
	return flattenedScopes
}

func getZoneClient(clients SysdigClients) (v2.ZoneInterface, error) {
	return clients.sysdigSecureClientV2()
}

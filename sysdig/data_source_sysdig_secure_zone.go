package sysdig

import (
	"context"
	"fmt"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureZoneRead,

		Schema: map[string]*schema.Schema{
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaIsSystemKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			SchemaAuthorKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaLastModifiedBy: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaLastUpdated: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaScopeKey: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaIDKey: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						SchemaTargetTypeKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaRulesKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The ID of the zone to retrieve.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"id", "name"},
				Description:  "The name of the zone to retrieve.",
			},
		},
	}
}

func dataSourceSysdigSecureZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	var zone *v2.Zone
	zoneIDRaw, hasZoneID := d.GetOk("id")
	if hasZoneID {
		zoneID, err := strconv.Atoi(zoneIDRaw.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid zone id: %s", err))
		}
		zone, err = client.GetZoneById(ctx, zoneID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zone by ID: %s", err))
		}
	} else if nameRaw, hasName := d.GetOk("name"); hasName {
		name := nameRaw.(string)
		zones, err := client.GetZones(ctx, name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zones: %s", err))
		}
		for _, z := range zones {
			if z.Name == name {
				zone = &z
				break
			}
		}
		if zone == nil {
			return diag.FromErr(fmt.Errorf("zone with name '%s' not found", name))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either id or name must be specified"))
	}

	d.SetId(fmt.Sprintf("%d", zone.ID))
	_ = d.Set(SchemaNameKey, zone.Name)
	_ = d.Set(SchemaDescriptionKey, zone.Description)
	_ = d.Set(SchemaIsSystemKey, zone.IsSystem)
	_ = d.Set(SchemaAuthorKey, zone.Author)
	_ = d.Set(SchemaLastModifiedBy, zone.LastModifiedBy)
	_ = d.Set(SchemaLastUpdated, time.UnixMilli(zone.LastUpdated).Format(time.RFC3339))

	if err := d.Set(SchemaScopeKey, fromZoneScopesResponse(zone.Scopes)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting scope: %s", err))
	}

	return nil
}

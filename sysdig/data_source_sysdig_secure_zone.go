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
						// Not marked Deprecated: rules with v2-compatible syntax are fully supported.
						// Only v1 syntax (labels, labelValues, agentTags) is deprecated, but since
						// this is a Computed field, SDK v2 has no mechanism for conditional deprecation.
						// The resource-side ValidateDiagFunc handles the v1-only warning.
						SchemaRulesKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaExpressionKey: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									SchemaFieldKey:    {Type: schema.TypeString, Computed: true},
									SchemaOperatorKey: {Type: schema.TypeString, Computed: true},
									SchemaValueKey:    {Type: schema.TypeString, Computed: true},
									SchemaValuesKey: {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
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

func dataSourceSysdigSecureZoneRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	clientV2, err := getZoneV2Client(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	var zoneV2 *v2.ZoneV2
	zoneIDRaw, hasZoneID := d.GetOk("id")
	if hasZoneID {
		zoneID, err := strconv.Atoi(zoneIDRaw.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zone by ID: %s", err))
		}
		zoneV2, err = clientV2.GetZoneV2(ctx, zoneID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zone v2 by ID: %s", err))
		}
	} else if nameRaw, hasName := d.GetOk("name"); hasName {
		name := nameRaw.(string)
		zones, err := clientV2.GetZonesV2(ctx, name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zones: %s", err))
		}
		for _, z := range zones {
			if z.Name == name {
				zoneV2 = &z
				break
			}
		}
		if zoneV2 == nil {
			return diag.FromErr(fmt.Errorf("zone with name '%s' not found", name))
		}
		zoneV2, err = clientV2.GetZoneV2(ctx, zoneV2.ID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error fetching zones: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either id or name must be specified"))
	}

	d.SetId(fmt.Sprintf("%d", zoneV2.ID))
	_ = d.Set(SchemaNameKey, zoneV2.Name)
	_ = d.Set(SchemaDescriptionKey, zoneV2.Description)
	_ = d.Set(SchemaIsSystemKey, zoneV2.IsSystem)
	_ = d.Set(SchemaAuthorKey, zoneV2.Author)
	_ = d.Set(SchemaLastModifiedBy, zoneV2.LastModifiedBy)
	_ = d.Set(SchemaLastUpdated, time.UnixMilli(zoneV2.LastUpdated).Format(time.RFC3339))

	if err := d.Set(SchemaScopeKey, getZoneScopes(zoneV2)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting scope: %s", err))
	}

	return nil
}

func getZoneScopes(zoneV2 *v2.ZoneV2) []any {
	// Build expression lookup by filter ID from the v2 response.
	out := make([]any, 0)
	if zoneV2 != nil {
		for _, s := range zoneV2.Scopes {
			for _, f := range s.Filters {
				if f.ID != 0 && len(f.Expressions) > 0 {
					var exprs []any
					for _, e := range f.Expressions {
						exprs = append(exprs, flattenExpressionV2(e))
					}
					m := map[string]any{
						SchemaIDKey:         f.ID,
						SchemaTargetTypeKey: f.ResourceType,
						SchemaRulesKey:      f.Rules,
					}
					m[SchemaExpressionKey] = exprs
					out = append(out, m)
				}
			}
		}
	}
	return out
}

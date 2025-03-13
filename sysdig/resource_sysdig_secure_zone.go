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

func resourceSysdigSecureZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSysdigSecureZoneCreate,
		ReadContext:   resourceSysdigSecureZoneRead,
		UpdateContext: resourceSysdigSecureZoneUpdate,
		DeleteContext: resourceSysdigSecureZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Optional: true,
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
			SchemaScopesKey: {
				Required: true,
				Type:     schema.TypeSet,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaScopeKey: {
							Type:     schema.TypeSet,
							MinItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									SchemaIDKey: {
										Type:     schema.TypeInt,
										Computed: true,
									},
									SchemaTargetTypeKey: {
										Type:     schema.TypeString,
										Required: true,
									},
									SchemaRulesKey: {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSysdigSecureZoneCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	return resourceSysdigSecureZoneRead(ctx, d, m)
}

func resourceSysdigSecureZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	_ = d.Set("scopes", fromZoneScopesResponse(zone.Scopes))
	_ = d.Set("is_system", zone.IsSystem)
	_ = d.Set("author", zone.Author)
	_ = d.Set("last_modified_by", zone.LastModifiedBy)
	_ = d.Set("last_updated", time.UnixMilli(zone.LastUpdated).Format(time.RFC3339))

	return nil
}

func resourceSysdigSecureZoneUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := getZoneClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneRequest := zoneRequestFromResourceData(d)

	_, err = client.UpdateZone(ctx, zoneRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Sysdig Zone: %s", err))
	}

	return resourceSysdigSecureZoneRead(ctx, d, m)
}

func resourceSysdigSecureZoneDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	zoneRequest := &v2.ZoneRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Scopes:      toZoneScopesRequest(d.Get("scopes").(*schema.Set)),
	}

	if d.Id() != "" {
		id, err := strconv.Atoi(d.Id())
		if err == nil {
			zoneRequest.ID = id
		}
	}

	return zoneRequest
}

func toZoneScopesRequest(scopes *schema.Set) []v2.ZoneScope {
	var zoneScopes []v2.ZoneScope
	for _, scopeData := range scopes.List() {
		scopeMap := scopeData.(map[string]interface{})
		scopeSet := scopeMap[SchemaScopeKey].(*schema.Set)
		for _, attr := range scopeSet.List() {
			s := attr.(map[string]interface{})
			zoneScopes = append(zoneScopes, v2.ZoneScope{
				ID:         s[SchemaIDKey].(int),
				TargetType: s[SchemaTargetTypeKey].(string),
				Rules:      s[SchemaRulesKey].(string),
			})
		}
	}
	return zoneScopes
}

func fromZoneScopesResponse(scopes []v2.ZoneScope) []interface{} {
	var flattenedScopes []interface{}
	for _, scope := range scopes {
		flattenedScopes = append(flattenedScopes, map[string]interface{}{
			SchemaIDKey:         scope.ID,
			SchemaTargetTypeKey: scope.TargetType,
			SchemaRulesKey:      scope.Rules,
		})
	}
	response := []interface{}{
		map[string]interface{}{
			SchemaScopeKey: schema.NewSet(schema.HashResource(&schema.Resource{
				Schema: map[string]*schema.Schema{
					SchemaIDKey: {
						Type:     schema.TypeInt,
						Computed: true,
					},
					SchemaTargetTypeKey: {
						Type:     schema.TypeString,
						Required: true,
					},
					SchemaRulesKey: {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			}), flattenedScopes),
		},
	}

	return response
}

func getZoneClient(clients SysdigClients) (v2.ZoneInterface, error) {
	return clients.sysdigSecureClientV2()
}

package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourceSysdigSecurePostureZone() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceCreateOrUpdatePostureZone,
		UpdateContext: resourceCreateOrUpdatePostureZone,
		DeleteContext: resourceSysdigSecurePostureZoneDelete,
		ReadContext:   resourceSysdigSecurePostureZoneRead,
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
			SchemaNameKey: {
				Required: true,
				Type:     schema.TypeString,
			},
			SchemaDescriptionKey: {
				Optional: true,
				Type:     schema.TypeString,
			},
			SchemaPolicyIDsKey: {
				Optional: true,
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			SchemaAuthorKey: {
				Computed: true,
				Type:     schema.TypeString,
			},
			SchemaLastModifiedBy: {
				Computed: true,
				Type:     schema.TypeString,
			},
			SchemaLastUpdated: {
				Computed: true,
				Type:     schema.TypeString,
			},
			SchemaScopesKey: {
				Optional: true,
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

func resourceCreateOrUpdatePostureZone(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	commonClient, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	identityCtx, err := commonClient.GetIdentityContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	policiesData := d.Get(SchemaPolicyIDsKey).([]interface{})
	policies := make([]string, len(policiesData))
	for i, p := range policiesData {
		policies[i] = p.(string)
	}

	scopesList := d.Get(SchemaScopesKey).(*schema.Set).List()
	scopes := make([]v2.PostureZoneScope, 0)
	if len(scopesList) > 0 {
		scopeList := scopesList[0].(map[string]interface{})[SchemaScopeKey].(*schema.Set).List()
		for _, attr := range scopeList {
			s := attr.(map[string]interface{})
			scopes = append(scopes, v2.PostureZoneScope{
				TargetType: s[SchemaTargetTypeKey].(string),
				Rules:      s[SchemaRulesKey].(string),
			})
		}
	}

	username := identityCtx.Username
	if identityCtx.ServiceAccountID != 0 {
		username = identityCtx.ServiceAccountName
	}

	req := &v2.PostureZoneRequest{
		ID:          d.Id(),
		Name:        d.Get(SchemaNameKey).(string),
		Description: d.Get(SchemaDescriptionKey).(string),
		PolicyIDs:   policies,
		Scopes:      scopes,
		Username:    username,
	}

	zoneClient, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := zoneClient.CreateOrUpdatePostureZone(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(zone.ID)

	resourceSysdigSecurePostureZoneRead(ctx, d, meta)
	return nil
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

	err = d.Set(SchemaNameKey, zone.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaDescriptionKey, zone.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaAuthorKey, zone.Author)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaLastModifiedBy, zone.LastModifiedBy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaLastUpdated, zone.LastUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	pIDs := make([]string, len(zone.Policies))
	for i, p := range zone.Policies {
		pIDs[i] = p.ID
	}
	err = d.Set(SchemaPolicyIDsKey, pIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	scopes := make([]map[string]interface{}, len(zone.Scopes))
	for i, s := range zone.Scopes {
		scopes[i] = map[string]interface{}{
			SchemaTargetTypeKey: s.TargetType,
			SchemaRulesKey:      s.Rules,
		}
	}
	if len(scopes) > 0 {
		err = d.Set(SchemaScopesKey, []interface{}{
			map[string]interface{}{
				SchemaScopeKey: scopes,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

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

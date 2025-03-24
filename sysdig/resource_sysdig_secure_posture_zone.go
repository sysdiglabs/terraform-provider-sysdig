package sysdig

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecurePostureZone() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceCreatePostureZone,
		UpdateContext: resourceUpdatePostureZone,
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
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
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

func resourceCreatePostureZone(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policies, err := getPolicies(d)
	if err != nil {
		return diag.FromErr(err)
	}

	scopes, err := getScopes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneClient, err := getZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	postureZoneClient, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneRequest := &v2.ZoneRequest{
		Name:        d.Get(SchemaNameKey).(string),
		Description: d.Get(SchemaDescriptionKey).(string),
		Scopes:      scopes,
	}

	zone, err := zoneClient.CreateZone(ctx, zoneRequest)
	if err != nil {
		return diag.Errorf("Error creating resource: %s", err)
	}

	policyIDs, err := convertPoliciesToInt(policies)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &v2.ZonePoliciesRequest{
		ZoneID:    zone.ID,
		PolicyIDs: policyIDs,
	}

	err = postureZoneClient.BindZoneToPolicies(ctx, req)
	if err != nil {
		log.Err(err).Int("zone_id", zone.ID).Msg("Error attaching zone to policies... deleting created zone")
		err2 := zoneClient.DeleteZone(ctx, zone.ID)
		if err2 != nil {
			return diag.Errorf("Error deleting zone [zone ID = %d] after failed attaching to policies: %s", zone.ID, err2)
		}
		return diag.Errorf("Error attaching zone to policies: %s", err)
	}

	d.SetId(strconv.Itoa(zone.ID))
	return resourceSysdigSecurePostureZoneRead(ctx, d, meta)
}

func resourceUpdatePostureZone(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	policies, err := getPolicies(d)
	if err != nil {
		return diag.FromErr(err)
	}

	scopes, err := getScopes(d)
	if err != nil {
		return diag.FromErr(err)
	}

	zoneClient, err := getZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	postureZoneClient, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Error updating posture zone resource, ID is not integer: %s", d.Id())
	}

	zoneRequest := &v2.ZoneRequest{
		ID:          id,
		Name:        d.Get(SchemaNameKey).(string),
		Description: d.Get(SchemaDescriptionKey).(string),
		Scopes:      scopes,
	}

	if d.HasChange(SchemaNameKey) || d.HasChange(SchemaDescriptionKey) || d.HasChange(SchemaScopesKey) {
		_, err = zoneClient.UpdateZone(ctx, zoneRequest)
		if err != nil {
			return diag.Errorf("Error updating resource: %s", err)
		}
	}

	policyIDs, err := convertPoliciesToInt(policies)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange(SchemaPolicyIDsKey) {
		req := &v2.ZonePoliciesRequest{
			ZoneID:    id,
			PolicyIDs: policyIDs,
		}

		err = postureZoneClient.BindZoneToPolicies(ctx, req)
		if err != nil {
			log.Err(err).Int("zone_id", id).Msg("Error attaching zone to policies")
			return diag.Errorf("Error attaching zone to policies: %s", err)
		}
	}

	return resourceSysdigSecurePostureZoneRead(ctx, d, meta)
}

func resourceSysdigSecurePostureZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	zone, err := client.GetPostureZone(ctx, id)
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

	pIDs := make([]int, len(zone.Policies))
	for i, p := range zone.Policies {
		id, err := strconv.Atoi(p.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		pIDs[i] = id
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
	postureClient, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	zoneClient, err := getZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = removeZoneFromPolicies(ctx, postureClient, id)
	if err != nil {
		return diag.Errorf("Error removing zone from policies: %s", err)
	}

	err = zoneClient.DeleteZone(ctx, id)
	if err != nil {

		return diag.Errorf("Error deleting zone: %s", err)
	}

	return nil
}

func removeZoneFromPolicies(ctx context.Context, client v2.PostureZoneInterface, zoneID int) error {
	req := &v2.ZonePoliciesRequest{
		ZoneID:    zoneID,
		PolicyIDs: []int{},
	}

	return client.BindZoneToPolicies(ctx, req)
}

func getPolicies(d *schema.ResourceData) ([]string, error) {
	policiesData := d.Get(SchemaPolicyIDsKey).(*schema.Set).List()
	policies := make([]string, len(policiesData))
	for i, p := range policiesData {
		policies[i] = strconv.Itoa(p.(int))
	}
	return policies, nil
}

func getScopes(d *schema.ResourceData) ([]v2.ZoneScope, error) {
	scopesList := d.Get(SchemaScopesKey).(*schema.Set).List()
	scopes := make([]v2.ZoneScope, 0)
	if len(scopesList) > 0 {
		scopeList := scopesList[0].(map[string]interface{})[SchemaScopeKey].(*schema.Set).List()
		for _, attr := range scopeList {
			s := attr.(map[string]interface{})
			scopes = append(scopes, v2.ZoneScope{
				TargetType: s[SchemaTargetTypeKey].(string),
				Rules:      s[SchemaRulesKey].(string),
			})
		}
	}
	return scopes, nil
}

func convertPoliciesToInt(policies []string) ([]int, error) {
	policyIDs := make([]int, len(policies))
	for i, p := range policies {
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("error converting policy ID to int: %s", err)
		}
		policyIDs[i] = id
	}
	return policyIDs, nil
}

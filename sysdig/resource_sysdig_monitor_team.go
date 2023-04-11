package sysdig

import (
	"context"
	"fmt"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	TeamSchemaEnableIBMPlatformMetricsKey   = "enable_ibm_platform_metrics"
	TeamSchemaIBMPlatformMetricsKey         = "ibm_platform_metrics"
	TeamSchemaCanUseSysdigCaptureKey        = "can_use_sysdig_capture"
	TeamSchemaCanSeeInfrastructureEventsKey = "can_see_infrastructure_events"
	TeamSchemaCanUseAWSDataKey              = "can_use_aws_data"
	TeamSchemaEntrypointKey                 = "entrypoint"
	TeamSchemaEntrypointTypeKey             = "type"
	TeamSchemaEntrypointSelectionKey        = "selection"
)

func createBaseMonitorTeamSchema() map[string]*schema.Schema {
	s := createBaseTeamSchema()
	MergeMap(s, map[string]*schema.Schema{
		TeamSchemaEnableIBMPlatformMetricsKey: {
			Type: schema.TypeBool,
		},
		TeamSchemaIBMPlatformMetricsKey: {
			Type: schema.TypeString,
		},
		TeamSchemaCanUseSysdigCaptureKey: {
			Type: schema.TypeBool,
		},
		TeamSchemaCanSeeInfrastructureEventsKey: {
			Type: schema.TypeBool,
		},
		TeamSchemaCanUseAWSDataKey: {
			Type: schema.TypeBool,
		},
		TeamSchemaEntrypointKey: {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					TeamSchemaEntrypointTypeKey: {
						Type: schema.TypeString,
					},
					TeamSchemaEntrypointSelectionKey: {
						Type: schema.TypeString,
					},
				},
			},
		},
	})

	return s
}

func createMonitorTeamSchema() map[string]*schema.Schema {
	s := createBaseMonitorTeamSchema()

	s[TeamSchemaThemeKey].Default = "#05C391"
	s[TeamSchemaThemeKey].Optional = true

	s[TeamSchemaNameKey].Required = true

	s[TeamSchemaDescriptionKey].Optional = true

	s[TeamSchemaScopeByKey].Default = "host"
	s[TeamSchemaScopeByKey].Optional = true

	s[TeamSchemaFilterKey].Optional = true

	s[TeamSchemaEnableIBMPlatformMetricsKey].Optional = true
	s[TeamSchemaIBMPlatformMetricsKey].Optional = true

	s[TeamSchemaCanUseSysdigCaptureKey].Default = false
	s[TeamSchemaCanUseSysdigCaptureKey].Optional = true

	s[TeamSchemaCanSeeInfrastructureEventsKey].Default = false
	s[TeamSchemaCanSeeInfrastructureEventsKey].Optional = true

	s[TeamSchemaCanUseAWSDataKey].Default = false
	s[TeamSchemaCanUseAWSDataKey].Optional = true

	s[TeamSchemaUserRolesKey].Optional = true
	userRolesSchema := s[TeamSchemaUserRolesKey].Elem.(*schema.Resource).Schema
	userRolesSchema[TeamSchemaUserRolesEmailKey].Required = true
	userRolesSchema[TeamSchemaUserRolesRoleKey].Default = "ROLE_TEAM_STANDARD"
	userRolesSchema[TeamSchemaUserRolesRoleKey].Optional = true
	userRolesSchema[TeamSchemaUserRolesRoleKey].ValidateFunc = validation.StringInSlice([]string{
		"ROLE_TEAM_STANDARD", "ROLE_TEAM_EDIT", "ROLE_TEAM_READ", "ROLE_TEAM_MANAGER",
	}, false)

	s[TeamSchemaEntrypointKey].Required = true
	entrypointSchema := s[TeamSchemaEntrypointKey].Elem.(*schema.Resource).Schema
	entrypointSchema[TeamSchemaEntrypointTypeKey].Required = true
	entrypointSchema[TeamSchemaEntrypointTypeKey].ValidateFunc = validation.StringInSlice([]string{
		"Explore", "Dashboards", "Events", "Alerts", "Settings",
	}, false)
	entrypointSchema[TeamSchemaEntrypointSelectionKey].Optional = true

	s[TeamSchemaDefaultTeamKey].Default = false
	s[TeamSchemaDefaultTeamKey].Optional = true

	s[TeamSchemaVersionKey].Computed = true

	return s
}

func resourceSysdigMonitorTeam() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorTeamCreate,
		UpdateContext: resourceSysdigMonitorTeamUpdate,
		ReadContext:   resourceSysdigMonitorTeamRead,
		DeleteContext: resourceSysdigMonitorTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(5 * time.Minute), // Removing the team is for some reason slower.
		},

		Schema: createMonitorTeamSchema(),
	}
}

func getMonitorTeamClient(c SysdigClients) (v2.TeamInterface, error) {
	var client v2.TeamInterface
	var err error
	switch c.GetClientType() {
	case IBMMonitor:
		client, err = c.ibmMonitorClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigMonitorClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func resourceSysdigMonitorTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	team := teamFromResourceData(d, clients.GetClientType())
	team.Products = []string{"SDC"}

	team, err = client.CreateTeam(ctx, team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	resourceSysdigMonitorTeamRead(ctx, d, meta)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamById(ctx, id)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = teamMonitorToResourceData(d, clients, t)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	return nil
}

func teamMonitorToResourceData(d *schema.ResourceData, c SysdigClients, t v2.Team) error {
	d.SetId(strconv.Itoa(t.ID))

	err := d.Set(TeamSchemaVersionKey, t.Version)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaThemeKey, t.Theme)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaNameKey, t.Name)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaDescriptionKey, t.Description)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaScopeByKey, t.Show)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaFilterKey, t.Filter)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaCanUseSysdigCaptureKey, t.CanUseSysdigCapture)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaCanSeeInfrastructureEventsKey, t.CanUseCustomEvents)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaCanUseAWSDataKey, t.CanUseAwsMetrics)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaDefaultTeamKey, t.DefaultTeam)
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaUserRolesKey, userMonitorRolesToSet(t.UserRoles))
	if err != nil {
		return err
	}
	err = d.Set(TeamSchemaEntrypointKey, entrypointToSet(t.EntryPoint))
	if err != nil {
		return err
	}

	if c.GetClientType() == IBMMonitor {
		err = resourceSysdigMonitorTeamReadIBM(d, &t)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceSysdigMonitorTeamReadIBM(d *schema.ResourceData, t *v2.Team) error {
	var ibmPlatformMetrics *string
	if t.NamespaceFilters != nil {
		ibmPlatformMetrics = t.NamespaceFilters.IBMPlatformMetrics
	}
	err := d.Set(TeamSchemaEnableIBMPlatformMetricsKey, t.CanUseBeaconMetrics)
	if err != nil {
		return err
	}
	return d.Set(TeamSchemaIBMPlatformMetricsKey, ibmPlatformMetrics)
}

func userMonitorRolesToSet(userRoles []v2.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		if role.Admin { // Admins are added by default, so skip them
			continue
		}

		roleMap := map[string]interface{}{
			TeamSchemaUserRolesEmailKey: role.Email,
			TeamSchemaUserRolesRoleKey:  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func entrypointToSet(entrypoint *v2.EntryPoint) (res []map[string]interface{}) {
	if entrypoint == nil {
		return
	}

	entrypointMap := map[string]interface{}{
		TeamSchemaEntrypointTypeKey:      entrypoint.Module,
		TeamSchemaEntrypointSelectionKey: entrypoint.Selection,
	}
	return append(res, entrypointMap)
}

func resourceSysdigMonitorTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clients := meta.(SysdigClients)
	client, err := getMonitorTeamClient(clients)
	if err != nil {
		return diag.FromErr(err)
	}

	t := teamFromResourceData(d, clients.GetClientType())
	t.Products = []string{"SDC"}

	t.Version = d.Get(TeamSchemaVersionKey).(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(ctx, t)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigMonitorTeamRead(ctx, d, meta)
	return nil
}

func resourceSysdigMonitorTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getMonitorTeamClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteTeam(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func updateNamespaceFilters(filters *v2.NamespaceFilters, update v2.NamespaceFilters) *v2.NamespaceFilters {
	if filters == nil {
		filters = &v2.NamespaceFilters{}
	}

	if update.IBMPlatformMetrics != nil {
		filters.IBMPlatformMetrics = update.IBMPlatformMetrics
	}

	return filters
}

func teamFromResourceData(d *schema.ResourceData, clientType ClientType) v2.Team {
	canUseSysdigCapture := d.Get(TeamSchemaCanUseSysdigCaptureKey).(bool)
	canUseCustomEvents := d.Get(TeamSchemaCanSeeInfrastructureEventsKey).(bool)
	canUseAwsMetrics := d.Get(TeamSchemaCanUseAWSDataKey).(bool)
	canUseBeaconMetrics := false
	t := v2.Team{
		Theme:               d.Get(TeamSchemaThemeKey).(string),
		Name:                d.Get(TeamSchemaNameKey).(string),
		Description:         d.Get(TeamSchemaDescriptionKey).(string),
		Show:                d.Get(TeamSchemaScopeByKey).(string),
		Filter:              d.Get(TeamSchemaFilterKey).(string),
		CanUseSysdigCapture: &canUseSysdigCapture,
		CanUseCustomEvents:  &canUseCustomEvents,
		CanUseAwsMetrics:    &canUseAwsMetrics,
		CanUseBeaconMetrics: &canUseBeaconMetrics,
		DefaultTeam:         d.Get(TeamSchemaDefaultTeamKey).(bool),
	}

	userRoles := make([]v2.UserRoles, 0)
	for _, userRole := range d.Get(TeamSchemaUserRolesKey).(*schema.Set).List() {
		ur := userRole.(map[string]interface{})
		userRoles = append(userRoles, v2.UserRoles{
			Email: ur[TeamSchemaUserRolesEmailKey].(string),
			Role:  ur[TeamSchemaUserRolesRoleKey].(string),
		})
	}
	t.UserRoles = userRoles

	t.EntryPoint = &v2.EntryPoint{}
	t.EntryPoint.Module = d.Get(fmt.Sprintf("%s.0.%s", TeamSchemaEntrypointKey, TeamSchemaEntrypointTypeKey)).(string)
	if val, ok := d.GetOk(fmt.Sprintf("%s.0.%s", TeamSchemaEntrypointKey, TeamSchemaEntrypointSelectionKey)); ok {
		t.EntryPoint.Selection = val.(string)
	}

	if clientType == IBMMonitor {
		teamFromResourceDataIBM(d, &t)
	}

	return t
}

func teamFromResourceDataIBM(d *schema.ResourceData, t *v2.Team) {
	canUseBeaconMetrics := d.Get(TeamSchemaEnableIBMPlatformMetricsKey).(bool)
	t.CanUseBeaconMetrics = &canUseBeaconMetrics

	if v, ok := d.GetOk(TeamSchemaIBMPlatformMetricsKey); ok {
		metrics := v.(string)
		t.NamespaceFilters = updateNamespaceFilters(t.NamespaceFilters, v2.NamespaceFilters{
			IBMPlatformMetrics: &metrics,
		})
	}
}

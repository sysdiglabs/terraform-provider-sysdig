package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/monitor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

func resourceSysdigMonitorTeam() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigMonitorTeamCreate,
		Update: resourceSysdigMonitorTeamUpdate,
		Read:   resourceSysdigMonitorTeamRead,
		Delete: resourceSysdigMonitorTeamDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"theme": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "#05C391",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_by": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "host",
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"can_use_sysdig_capture": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_see_infrastructure_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_use_aws_data": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"user_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"role": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "ROLE_TEAM_STANDARD",
							ValidateFunc: validation.StringInSlice([]string{"ROLE_TEAM_STANDARD", "ROLE_TEAM_EDIT", "ROLE_TEAM_READ", "ROLE_TEAM_MANAGER"}, false),
						},
					},
				},
			},
			"entrypoint": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Explore", "Dashboards", "Events", "Alerts", "Settings"}, false),
						},

						"selection": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"default_team": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigMonitorTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	team := teamFromResourceData(d)

	team, err = client.CreateTeam(team)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(team.ID))
	d.Set("version", team.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigMonitorTeamRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamById(id)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("version", t.Version)
	d.Set("theme", t.Theme)
	d.Set("name", t.Name)
	d.Set("description", t.Description)
	d.Set("scope_by", t.Show)
	d.Set("filter", t.Filter)
	d.Set("can_use_sysdig_capture", t.CanUseSysdigCapture)
	d.Set("default_team", t.DefaultTeam)
	d.Set("user_roles", userMonitorRolesToSet(t.UserRoles))
	d.Set("entrypoint", entrypointToSet(t.EntryPoint))

	return nil
}

func userMonitorRolesToSet(userRoles []monitor.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		if role.Admin { // Admins are added by default, so skip them
			continue
		}

		roleMap := map[string]interface{}{
			"email": role.Email,
			"role":  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func entrypointToSet(entrypoint monitor.EntryPoint) (res []map[string]interface{}) {
	entrypointMap := map[string]interface{}{
		"type":      entrypoint.Module,
		"selection": entrypoint.Selection,
	}
	return append(res, entrypointMap)
}

func resourceSysdigMonitorTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	t := teamFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(t)

	return err
}

func resourceSysdigMonitorTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigMonitorClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteTeam(id)
}

func teamFromResourceData(d *schema.ResourceData) monitor.Team {
	t := monitor.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Show:                d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: d.Get("can_use_sysdig_capture").(bool),
		CanUseCustomEvents:  d.Get("can_see_infrastructure_events").(bool),
		CanUseAwsMetrics:    d.Get("can_use_aws_data").(bool),
		DefaultTeam:         d.Get("default_team").(bool),
	}

	userRoles := []monitor.UserRoles{}
	for _, userRole := range d.Get("user_roles").(*schema.Set).List() {
		ur := userRole.(map[string]interface{})
		userRoles = append(userRoles, monitor.UserRoles{
			Email: ur["email"].(string),
			Role:  ur["role"].(string),
		})
	}
	t.UserRoles = userRoles

	t.EntryPoint.Module = d.Get("entrypoint.0.type").(string)
	if val, ok := d.GetOk("entrypoint.0.selection"); ok {
		t.EntryPoint.Selection = val.(string)
	}

	return t
}

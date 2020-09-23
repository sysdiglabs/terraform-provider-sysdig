package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"time"
)

func resourceSysdigSecureTeam() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigSecureTeamCreate,
		Update: resourceSysdigSecureTeamUpdate,
		Read:   resourceSysdigSecureTeamRead,
		Delete: resourceSysdigSecureTeamDelete,

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
				Default:  "#73A1F7",
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
				Default:  "container",
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_sysdig_capture": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"user_roles": {
				Type:     schema.TypeSet,
				Required: true,
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

func resourceSysdigSecureTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	team := secureTeamFromResourceData(d)

	team, err = client.CreateTeam(team)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(team.ID))
	d.Set("version", team.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigSecureTeamRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
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
	d.Set("scope_by", t.ScopeBy)
	d.Set("filter", t.Filter)
	d.Set("use_sysdig_capture", t.CanUseSysdigCapture)
	d.Set("default_team", t.DefaultTeam)
	d.Set("user_roles", userSecureRolesToSet(t.UserRoles))

	return nil
}

func userSecureRolesToSet(userRoles []secure.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		roleMap := map[string]interface{}{
			"email": role.Email,
			"role":  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func resourceSysdigSecureTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	t := secureTeamFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(t)

	return err
}

func resourceSysdigSecureTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteTeam(id)
}

func secureTeamFromResourceData(d *schema.ResourceData) secure.Team {
	t := secure.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		ScopeBy:             d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: d.Get("use_sysdig_capture").(bool),
		DefaultTeam:         d.Get("default_team").(bool),
	}

	userRoles := []secure.UserRoles{}
	for _, userRole := range d.Get("user_roles").(*schema.Set).List() {
		ur := userRole.(map[string]interface{})
		userRoles = append(userRoles, secure.UserRoles{
			Email: ur["email"].(string),
			Role:  ur["role"].(string),
		})
	}
	t.UserRoles = userRoles

	return t
}

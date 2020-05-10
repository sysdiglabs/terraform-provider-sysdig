package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigSecureTeam() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigTeamCreate,
		Update: resourceSysdigTeamUpdate,
		Read:   resourceSysdigTeamRead,
		Delete: resourceSysdigTeamDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
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
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},

						"role": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ROLE_TEAM_STANDARD",
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

func resourceSysdigTeamCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
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
func resourceSysdigTeamRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("canUseSysdigCapture", t.CanUseSysdigCapture)
	d.Set("default_team", t.DefaultTeam)
	d.Set("user_roles", t.UserRoles)

	return nil
}

func resourceSysdigTeamUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	t := teamFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(t)

	return err
}

func resourceSysdigTeamDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteTeam(id)
}

func teamFromResourceData(d *schema.ResourceData) secure.Team {
	t := secure.Team{
		Theme:               d.Get("theme").(string),
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		ScopeBy:             d.Get("scope_by").(string),
		Filter:              d.Get("filter").(string),
		CanUseSysdigCapture: d.Get("use_sysdig_capture").(bool),
		DefaultTeam:         d.Get("default_team").(bool),
		Products:            []string{"SDS"},
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

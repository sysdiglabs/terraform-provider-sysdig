package sysdig

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

func resourceSysdigSecureTeam() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureTeamCreate,
		UpdateContext: resourceSysdigSecureTeamUpdate,
		ReadContext:   resourceSysdigSecureTeamRead,
		DeleteContext: resourceSysdigSecureTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(5 * time.Minute), // Removing the team is for some reason slower.
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

func resourceSysdigSecureTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	team := secureTeamFromResourceData(d)

	team, err = client.CreateTeam(ctx, team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(team.ID))
	err = d.Set("version", team.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}
	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigSecureTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	t, err := client.GetTeamById(ctx, id)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	err = d.Set("version", t.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}

	err = d.Set("theme", t.Theme)
	if err != nil {
		log.Println("error assigning 'theme'")
	}

	err = d.Set("name", t.Name)
	if err != nil {
		log.Println("error assigning 'name'")
	}

	err = d.Set("description", t.Description)
	if err != nil {
		log.Println("error assigning 'description'")
	}

	err = d.Set("scope_by", t.ScopeBy)
	if err != nil {
		log.Println("error assigning 'scope_by'")
	}

	err = d.Set("filter", t.Filter)
	if err != nil {
		log.Println("error assigning 'filter'")
	}

	err = d.Set("use_sysdig_capture", t.CanUseSysdigCapture)
	if err != nil {
		log.Println("error assigning 'use_sysdig_capture'")
	}

	err = d.Set("default_team", t.DefaultTeam)
	if err != nil {
		log.Println("error assigning 'default_team'")
	}

	err = d.Set("user_roles", userSecureRolesToSet(t.UserRoles))
	if err != nil {
		log.Println("error assigning 'user_roles'")
	}

	return nil
}

func userSecureRolesToSet(userRoles []secure.UserRoles) (res []map[string]interface{}) {
	for _, role := range userRoles {
		if role.Admin {
			continue // Admins are added by default, so skip them
		}
		roleMap := map[string]interface{}{
			"email": role.Email,
			"role":  role.Role,
		}
		res = append(res, roleMap)
	}
	return
}

func resourceSysdigSecureTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	t := secureTeamFromResourceData(d)

	t.Version = d.Get("version").(int)
	t.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateTeam(ctx, t)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
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

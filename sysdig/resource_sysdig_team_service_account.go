package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigTeamServiceAccount() *schema.Resource {
	timeout := 5 * time.Minute
	return &schema.Resource{
		ReadContext:   resourceSysdigTeamServiceAccountRead,
		CreateContext: resourceSysdigTeamServiceAccountCreate,
		UpdateContext: resourceSysdigTeamServiceAccountUpdate,
		DeleteContext: resourceSysdigTeamServiceAccountDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ROLE_TEAM_READ",
			},
			"expiration_date": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"team_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"system_role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ROLE_SERVICE_ACCOUNT",
			},
			"date_created": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSysdigTeamServiceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		diag.FromErr(err)
	}

	teamServiceAccount, err := client.GetTeamServiceAccountById(ctx, id)
	if err != nil {
		if err == v2.TeamServiceAccountNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = teamServiceAccountToResourceData(teamServiceAccount, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigTeamServiceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	teamServiceAccount := teamServiceAccountFromResourceData(d)
	teamServiceAccount, err = client.CreateTeamServiceAccount(ctx, teamServiceAccount)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(teamServiceAccount.ID))

	resourceSysdigTeamServiceAccountRead(ctx, d, m)

	return nil
}

func resourceSysdigTeamServiceAccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var err error

	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	teamServiceAccount := teamServiceAccountFromResourceData(d)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	teamServiceAccount.ID = id
	_, err = client.UpdateTeamServiceAccount(ctx, teamServiceAccount, id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigTeamServiceAccountRead(ctx, d, m)

	return nil
}

func resourceSysdigTeamServiceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.DeleteTeamServiceAccount(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func teamServiceAccountFromResourceData(d *schema.ResourceData) *v2.TeamServiceAccount {
	return &v2.TeamServiceAccount{
		Name:           d.Get("name").(string),
		TeamRole:       d.Get("role").(string),
		ExpirationDate: int64(d.Get("expiration_date").(int) * 1000),
		TeamId:         d.Get("team_id").(int),
		SystemRole:     d.Get("system_role").(string),
		ApiKey:         d.Get("api_key").(string),
	}
}

func teamServiceAccountToResourceData(teamServiceAccount *v2.TeamServiceAccount, d *schema.ResourceData) error {
	err := d.Set("name", teamServiceAccount.Name)
	if err != nil {
		return err
	}
	err = d.Set("role", teamServiceAccount.TeamRole)
	if err != nil {
		return err
	}
	err = d.Set("expiration_date", teamServiceAccount.ExpirationDate/1000)
	if err != nil {
		return err
	}
	err = d.Set("team_id", teamServiceAccount.TeamId)
	if err != nil {
		return err
	}
	err = d.Set("system_role", teamServiceAccount.SystemRole)
	if err != nil {
		return err
	}
	err = d.Set("date_created", teamServiceAccount.DateCreated)
	if err != nil {
		return err
	}
	err = d.Set("api_key", teamServiceAccount.ApiKey)
	if err != nil {
		return err
	}

	return nil
}

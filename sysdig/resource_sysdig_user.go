package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
)

func resourceSysdigUser() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigUserCreate,
		UpdateContext: resourceSysdigUserUpdate,
		ReadContext:   resourceSysdigUserRead,
		DeleteContext: resourceSysdigUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ROLE_USER",
			},
			"first_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user := userFromResourceData(d)

	user, err = client.CreateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(user.ID))
	d.Set("version", user.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	u, err := client.GetUserById(ctx, id)

	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.Set("version", u.Version)
	d.Set("system_role", u.SystemRole)
	d.Set("email", u.Email)
	d.Set("first_name", u.FirstName)
	d.Set("last_name", u.LastName)

	return nil
}

func resourceSysdigUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user := userFromResourceData(d)

	user.Version = d.Get("version").(int)
	user.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateUser(ctx, user)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteUser(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func userFromResourceData(d *schema.ResourceData) (u *common.User) {
	u = &common.User{
		SystemRole: d.Get("system_role").(string),
		Email:      d.Get("email").(string),
		FirstName:  d.Get("first_name").(string),
		LastName:   d.Get("last_name").(string),
	}
	return u
}

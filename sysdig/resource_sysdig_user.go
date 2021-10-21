package sysdig

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/common"
)

func resourceSysdigUser() *schema.Resource {
	timeout := 5 * time.Minute

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
	err = d.Set("version", user.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}

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

	err = d.Set("version", u.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}
	err = d.Set("system_role", u.SystemRole)
	if err != nil {
		log.Println("error assigning 'system_role'")
	}
	err = d.Set("email", u.Email)
	if err != nil {
		log.Println("error assigning 'email'")
	}
	err = d.Set("first_name", u.FirstName)
	if err != nil {
		log.Println("error assigning 'first_name'")
	}
	err = d.Set("last_name", u.LastName)
	if err != nil {
		log.Println("error assigning 'last_name'")
	}
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

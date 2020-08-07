package sysdig

import (
	"github.com/draios/terraform-provider-sysdig/sysdig/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceSysdigUser() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigUserCreate,
		Update: resourceSysdigUserUpdate,
		Read:   resourceSysdigUserRead,
		Delete: resourceSysdigUserDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				Sensitive: true,
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

func resourceSysdigUserCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return err
	}

	user := userFromResourceData(d)

	user, err = client.CreateUser(user)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(user.ID))
	d.Set("version", user.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigUserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())
	u, err := client.GetUserById(id)

	if err != nil {
		d.SetId("")
		return err
	}

	d.Set("version", u.Version)
	d.Set("system_role", u.SystemRole)
	d.Set("email", u.Email)
	d.Set("first_name", u.FirstName)
	d.Set("last_name", u.LastName)

	return nil
}

func resourceSysdigUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return err
	}

	u := userFromResourceData(d)

	u.Version = d.Get("version").(int)
	u.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateUser(u)

	return err
}

func resourceSysdigUserDelete(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return err
	}

	id, _ := strconv.Atoi(d.Id())

	return client.DeleteUser(id)
}

func userFromResourceData(d *schema.ResourceData) (u common.User) {
	u = common.User{
		SystemRole: d.Get("system_role").(string),
		Password:   d.Get("password").(string),
		Email:      d.Get("email").(string),
		FirstName:  d.Get("first_name").(string),
		LastName:   d.Get("last_name").(string),
	}
	return u
}

package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigCurrentUser() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Read: dataSourceSysdigCurrentUserRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"system_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigCurrentUserRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return err
	}

	user, err := client.GetCurrentUser()
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(user.ID))
	d.Set("email", user.Email)
	d.Set("name", user.FirstName)
	d.Set("last_name", user.LastName)
	d.Set("system_role", user.SystemRole)

	return nil
}

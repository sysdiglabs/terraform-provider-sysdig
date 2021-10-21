package sysdig

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigCurrentUser() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigCurrentUserRead,

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
func dataSourceSysdigCurrentUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(user.ID))
	err = d.Set("email", user.Email)
	if err != nil {
		log.Println("error assigning 'email'")
	}
	err = d.Set("name", user.FirstName)
	if err != nil {
		log.Println("error assigning 'name'")
	}
	err = d.Set("last_name", user.LastName)
	if err != nil {
		log.Println("error assigning 'last_name'")
	}
	err = d.Set("system_role", user.SystemRole)
	if err != nil {
		log.Println("error assigning 'system_role'")
	}

	return nil
}

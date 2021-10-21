package sysdig

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigUser() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigUserRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSysdigUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClient()
	if err != nil {
		return diag.FromErr(err)
	}

	u, err := client.GetUserByEmail(ctx, d.Get("email").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(u.ID))
	err = d.Set("version", u.Version)
	if err != nil {
		log.Println("error asigning 'version' to alert")
	}

	err = d.Set("system_role", u.SystemRole)
	if err != nil {
		log.Println("error asigning 'system_role' to alert")
	}

	err = d.Set("first_name", u.FirstName)
	if err != nil {
		log.Println("error asigning 'first_name' to alert")
	}

	err = d.Set("last_name", u.LastName)
	if err != nil {
		log.Println("error asigning 'last_name' to alert")
	}

	return nil

}

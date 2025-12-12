package sysdig

import (
	"context"
	"net/http"
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

func dataSourceSysdigUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	u, statusCode, err := client.GetUserByEmail(ctx, d.Get("email").(string))
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		} else {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(u.ID))
	_ = d.Set("version", u.Version)
	_ = d.Set("system_role", u.SystemRole)
	_ = d.Set("first_name", u.FirstName)
	_ = d.Set("last_name", u.LastName)

	return nil
}

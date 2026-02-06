package sysdig

import (
	"context"
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
			"customer_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"customer_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigCurrentUserRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := meta.(SysdigClients).commonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.GetCurrentUser(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(user.ID))
	_ = d.Set("email", user.Email)
	_ = d.Set("name", user.FirstName)
	_ = d.Set("last_name", user.LastName)
	_ = d.Set("system_role", user.SystemRole)

	if user.Customer != nil {
		_ = d.Set("customer_id", user.Customer.ID)
		_ = d.Set("customer_name", user.Customer.Name)
		_ = d.Set("customer_external_id", user.Customer.ExternalID)
	} else {
		_ = d.Set("customer_id", nil)
		_ = d.Set("customer_name", nil)
		_ = d.Set("customer_external_id", nil)
	}

	return nil
}

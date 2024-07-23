package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceSysdigAllowedIpRange() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSysdigAllowedIpRangeRead,
		CreateContext: resourceSysdigAllowedIpRangeCreate,
		UpdateContext: resourceSysdigAllowedIpRangeUpdate,
		DeleteContext: resourceSysdigAllowedIpRangeDelete,
		Schema: map[string]*schema.Schema{
			"ip_range": {
				Type:     schema.TypeString,
				Required: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceSysdigAllowedIpRangeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	allowedIpRange, err := client.GetAllowedIpRangeById(ctx, id)
	if err != nil {
		if err == v2.AllowedIpRangeNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = ipRangeToResourceData(allowedIpRange, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigAllowedIpRangeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	allowedIpRange, err := ipRangeFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAllowedIpRange, err := client.CreateAllowedIpRange(ctx, allowedIpRange)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdAllowedIpRange.ID))

	resourceSysdigAllowedIpRangeRead(ctx, d, m)

	return nil

}

func resourceSysdigAllowedIpRangeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	allowedIpRange, err := ipRangeFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)

	}

	allowedIpRange.ID = id
	_, err = client.UpdateAllowedIpRange(ctx, allowedIpRange, id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigAllowedIpRangeRead(ctx, d, m)

	return nil
}

func resourceSysdigAllowedIpRangeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteAllowedIpRange(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ipRangeToResourceData(allowedIpRange *v2.AllowedIpRange, d *schema.ResourceData) error {
	err := d.Set("ip_range", allowedIpRange.IpRange)
	if err != nil {
		return err
	}
	err = d.Set("note", allowedIpRange.Note)
	if err != nil {
		return err
	}
	err = d.Set("enabled", allowedIpRange.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func ipRangeFromResourceData(d *schema.ResourceData) (*v2.AllowedIpRange, error) {
	return &v2.AllowedIpRange{
		IpRange: d.Get("ip_range").(string),
		Note:    d.Get("note").(string),
		Enabled: d.Get("enabled").(bool),
	}, nil
}

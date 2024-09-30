package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceSysdigIPFilter() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceSysdigIPFilterRead,
		CreateContext: resourceSysdigIPFilterCreate,
		UpdateContext: resourceSysdigIPFilterUpdate,
		DeleteContext: resourceSysdigIPFilterDelete,
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

func resourceSysdigIPFilterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	ipFilter, err := client.GetIPFilterById(ctx, id)
	if err != nil {
		if err == v2.IPFilterNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = ipFilterToResourceData(ipFilter, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigIPFilterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFilter, err := ipFilterFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createdIPFilter, err := client.CreateIPFilter(ctx, ipFilter)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(createdIPFilter.ID))

	resourceSysdigIPFilterRead(ctx, d, m)

	return nil
}

func resourceSysdigIPFilterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ipFilter, err := ipFilterFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)

	}

	ipFilter.ID = id
	_, err = client.UpdateIPFilter(ctx, ipFilter, id)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceSysdigIPFilterRead(ctx, d, m)

	return nil
}

func resourceSysdigIPFilterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client, err := m.(SysdigClients).sysdigCommonClientV2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteIPFilter(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ipFilterToResourceData(ipFilter *v2.IPFilter, d *schema.ResourceData) error {
	err := d.Set("ip_range", ipFilter.IPRange)
	if err != nil {
		return err
	}
	err = d.Set("note", ipFilter.Note)
	if err != nil {
		return err
	}
	err = d.Set("enabled", ipFilter.Enabled)
	if err != nil {
		return err
	}

	return nil
}

func ipFilterFromResourceData(d *schema.ResourceData) (*v2.IPFilter, error) {
	return &v2.IPFilter{
		IPRange: d.Get("ip_range").(string),
		Note:    d.Get("note").(string),
		Enabled: d.Get("enabled").(bool),
	}, nil
}

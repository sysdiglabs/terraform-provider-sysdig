package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
	"time"
)

func resourceSysdigSecureList() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigListCreate,
		UpdateContext: resourceSysdigListUpdate,
		ReadContext:   resourceSysdigListRead,
		DeleteContext: resourceSysdigListDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"items": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"append": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSysdigListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list := listFromResourceData(d)
	list, err = client.CreateList(ctx, list)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(list.ID))
	d.Set("version", list.Version)

	return nil
}

func resourceSysdigListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	list := listFromResourceData(d)
	list.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	list.ID = id

	_, err = client.UpdateList(ctx, list)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSysdigListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	list, err := client.GetListById(ctx, id)

	if err != nil {
		d.SetId("")
	}

	d.Set("name", list.Name)
	d.Set("version", list.Version)
	d.Set("items", list.Items.Items)
	d.Set("append", list.Append)

	return nil
}

func resourceSysdigListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteList(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func listFromResourceData(d *schema.ResourceData) secure.List {
	list := secure.List{
		Name:   d.Get("name").(string),
		Append: d.Get("append").(bool),
		Items:  secure.Items{Items: []string{}},
	}

	items := d.Get("items").([]interface{})
	for _, item := range items {
		if item_str, ok := item.(string); ok {
			item_str = strings.TrimSpace(item_str)
			list.Items.Items = append(list.Items.Items, item_str)
		}
	}

	return list
}

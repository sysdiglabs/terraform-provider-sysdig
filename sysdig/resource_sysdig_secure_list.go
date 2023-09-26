package sysdig

import (
	"context"
	"strconv"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureList() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigListCreate,
		UpdateContext: resourceSysdigListUpdate,
		ReadContext:   resourceSysdigListRead,
		DeleteContext: resourceSysdigListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

func getSecureListClient(c SysdigClients) (v2.ListInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureListClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	list := listFromResourceData(d)
	list, err = client.CreateList(ctx, list)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	d.SetId(strconv.Itoa(list.ID))
	_ = d.Set("version", list.Version)

	return nil
}

func resourceSysdigListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureListClient(sysdigClients)
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
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureListClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	list, err := client.GetListByID(ctx, id)

	if err != nil {
		d.SetId("")
	}

	_ = d.Set("name", list.Name)
	_ = d.Set("version", list.Version)
	_ = d.Set("items", list.Items.Items)
	_ = d.Set("append", list.Append)

	return nil
}

func resourceSysdigListDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureListClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteList(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func listFromResourceData(d *schema.ResourceData) v2.List {
	list := v2.List{
		Name:   d.Get("name").(string),
		Append: d.Get("append").(bool),
		Items:  v2.Items{Items: []string{}},
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

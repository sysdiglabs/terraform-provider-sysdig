package sysdig

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecurePoliciesPriority() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		Create: resourceSysdigSecurePoliciesPriorityCreate,
		Read:   resourceSysdigSecurePoliciesPriorityRead,
		Update: resourceSysdigSecurePoliciesPriorityCreate, // Create and update have the same behaviour
		Delete: resourceSysdigSecurePoliciesPriorityDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// Removes duplicates from a int slice
func uniqueIntSlice(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func resourceSysdigSecurePoliciesPriorityCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	priority := policiesPriorityFromResourceData(d, client)

	priority, err := client.CreatePoliciesPriority(priority)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(priority.Version))
	d.Set("policies", priority.PolicyIds)
	d.Set("version", priority.Version)

	return nil
}

func policiesPriorityFromResourceData(d *schema.ResourceData, client secure.SysdigSecureClient) (priority secure.PoliciesPriority) {
	policies := d.Get("policies").([]interface{})

	var policiesId []int

	for _, id := range policies {
		if idStr, ok := id.(string); ok {
			idInt, _ := strconv.Atoi(idStr)
			policiesId = append(policiesId, idInt)
		}
	}

	priority = secure.PoliciesPriority{
		PolicyIds: policiesId,
		Version:   d.Get("version").(int),
	}

	if client != nil {
		policiesPriority, err := client.GetPoliciesPriority()
		if err != nil {
			return
		}
		priority.PolicyIds = append(priority.PolicyIds, policiesPriority.PolicyIds...)
		priority.PolicyIds = uniqueIntSlice(priority.PolicyIds)
		priority.Version = policiesPriority.Version
	}

	return
}

func resourceSysdigSecurePoliciesPriorityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(secure.SysdigSecureClient)

	priority, err := client.GetPoliciesPriority()

	if err != nil {
		d.SetId("")
	}

	d.Set("policies", priority.PolicyIds)
	d.Set("version", priority.Version)

	return nil
}

func resourceSysdigSecurePoliciesPriorityDelete(d *schema.ResourceData, meta interface{}) error {
	return nil // TODO checkout delete in priority order
}

package sysdig

import (
	"context"
	"net/http"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureRuleContainer() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleContainerCreate,
		UpdateContext: resourceSysdigRuleContainerUpdate,
		ReadContext:   resourceSysdigRuleContainerRead,
		DeleteContext: resourceSysdigRuleContainerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"matching": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"containers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func resourceSysdigRuleContainerCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleContainerFromResourceData(d)

	rule, err = client.CreateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	d.SetId(strconv.Itoa(rule.ID))
	_ = d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleContainerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, statusCode, err := client.GetRuleByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	if rule.Details.Containers == nil {
		return diag.Errorf("no container data for a container rule")
	}

	updateResourceDataForRule(d, rule)
	_ = d.Set("matching", rule.Details.Containers.MatchItems)
	_ = d.Set("containers", rule.Details.Containers.Items)

	return nil
}

func resourceSysdigRuleContainerUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleContainerFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleContainerDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleContainerFromResourceData(d *schema.ResourceData) v2.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = v2.RuleTypeContainer

	rule.Details.Containers = &v2.Containers{}
	rule.Details.Containers.MatchItems = d.Get("matching").(bool)
	rule.Details.Containers.Items = []string{}
	if containers, ok := d.Get("containers").([]any); ok {
		for _, rawContainer := range containers {
			if container, ok := rawContainer.(string); ok {
				rule.Details.Containers.Items = append(rule.Details.Containers.Items, container)
			}
		}
	}

	return rule
}

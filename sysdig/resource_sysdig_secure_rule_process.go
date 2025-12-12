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

func resourceSysdigSecureRuleProcess() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleProcessCreate,
		UpdateContext: resourceSysdigRuleProcessUpdate,
		ReadContext:   resourceSysdigRuleProcessRead,
		DeleteContext: resourceSysdigRuleProcessDelete,
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
			"processes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func resourceSysdigRuleProcessCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleProcessFromResourceData(d)

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
func resourceSysdigRuleProcessRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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
			return nil
		} else {
			return diag.FromErr(err)
		}
	}

	if rule.Details.Processes == nil {
		return diag.Errorf("no process data for a process rule")
	}

	updateResourceDataForRule(d, rule)
	_ = d.Set("matching", rule.Details.Processes.MatchItems)
	_ = d.Set("processes", rule.Details.Processes.Items)

	return nil
}

func resourceSysdigRuleProcessUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule := resourceSysdigRuleProcessFromResourceData(d)

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleProcessDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
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

func resourceSysdigRuleProcessFromResourceData(d *schema.ResourceData) v2.Rule {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = v2.RuleTypeProcess

	rule.Details.Processes = &v2.Processes{}
	rule.Details.Processes.MatchItems = d.Get("matching").(bool)
	rule.Details.Processes.Items = []string{}
	if processes, ok := d.Get("processes").([]any); ok {
		for _, rawProcess := range processes {
			if process, ok := rawProcess.(string); ok {
				rule.Details.Processes.Items = append(rule.Details.Processes.Items, process)
			}
		}
	}

	return rule
}

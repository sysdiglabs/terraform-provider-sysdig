package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureRuleFilesystem() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleFilesystemCreate,
		UpdateContext: resourceSysdigRuleFilesystemUpdate,
		ReadContext:   resourceSysdigRuleFilesystemRead,
		DeleteContext: resourceSysdigRuleFilesystemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"read_only": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"paths": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"read_write": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"paths": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		}),
	}
}

func resourceSysdigRuleFilesystemCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleFilesystemFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err = client.CreateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rule.ID))
	_ = d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleFilesystemRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := client.GetRuleByID(ctx, id)
	if err != nil {
		d.SetId("")
	}

	updateResourceDataForRule(d, rule)

	if rule.Details.ReadPaths == nil {
		return diag.Errorf("no readPaths for a filesystem rule")
	}

	if rule.Details.ReadWritePaths == nil {
		return diag.Errorf("no readWritePaths for a filesystem rule")
	}

	if len(rule.Details.ReadPaths.Items) > 0 {
		_ = d.Set("read_only", []map[string]interface{}{{
			"matching": rule.Details.ReadPaths.MatchItems,
			"paths":    rule.Details.ReadPaths.Items,
		}})

	}
	if len(rule.Details.ReadWritePaths.Items) > 0 {
		_ = d.Set("read_write", []map[string]interface{}{{
			"matching": rule.Details.ReadWritePaths.MatchItems,
			"paths":    rule.Details.ReadWritePaths.Items,
		}})

	}

	return nil
}

func resourceSysdigRuleFilesystemUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleFilesystemFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigRuleFilesystemDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
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

	return nil
}

func resourceSysdigRuleFilesystemFromResourceData(d *schema.ResourceData) (rule v2.Rule, err error) {
	rule = ruleFromResourceData(d)
	rule.Details.RuleType = v2.RuleTypeFilesystem

	rule.Details.ReadPaths = &v2.ReadPaths{
		MatchItems: true,
		Items:      []string{},
	}
	rule.Details.ReadWritePaths = &v2.ReadWritePaths{
		MatchItems: true,
		Items:      []string{},
	}

	if readOnlyRules, ok := d.Get("read_only").([]interface{}); ok && len(readOnlyRules) > 0 {
		rule.Details.ReadPaths.MatchItems = d.Get("read_only.0.matching").(bool)
		for _, path := range d.Get("read_only.0.paths").([]interface{}) {
			if pathStr, ok := path.(string); ok {
				rule.Details.ReadPaths.Items = append(rule.Details.ReadPaths.Items, pathStr)
			}
		}

	}

	if readWriteRules, ok := d.Get("read_write").([]interface{}); ok && len(readWriteRules) > 0 {
		rule.Details.ReadWritePaths.MatchItems = d.Get("read_write.0.matching").(bool)
		for _, path := range d.Get("read_write.0.paths").([]interface{}) {
			if pathStr, ok := path.(string); ok {
				rule.Details.ReadWritePaths.Items = append(rule.Details.ReadWritePaths.Items, pathStr)
			}
		}
	}
	return
}

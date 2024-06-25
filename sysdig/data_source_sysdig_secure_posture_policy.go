package sysdig

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecurePosturePolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecurePosturePolicyRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaNameKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaTypeKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaLinkKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaMinKubeVersionKey: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			SchemaMaxKubeVersionKey: {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			SchemaIsActiveKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			SchemaPlatformKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaGroupKey: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     createGroupSchema(1),
			},
		},
	}
}

func dataSourceSysdigSecurePosturePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Get("id").(string), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	policy, err := client.GetPosturePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set(SchemaIDKey, policy.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaNameKey, policy.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaDescriptionKey, policy.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaTypeKey, policy.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaLinkKey, policy.Link)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaMinKubeVersionKey, policy.MinKubeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaMaxKubeVersionKey, policy.MaxKubeVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaIsActiveKey, policy.IsActive)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaPlatformKey, policy.Platform)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set groups
	if _, err := setGroups(d, policy.RequirementsGroup); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(policy.ID)
	return nil
}

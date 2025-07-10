package sysdig

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecurePostureZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecurePostureZoneRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"author": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rules": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSysdigSecurePostureZoneRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getPostureZoneClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	postureZone, err := client.GetPostureZoneByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(postureZone.ID)
	err = d.Set("name", postureZone.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("description", postureZone.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("author", postureZone.Author)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("last_modified_by", postureZone.LastModifiedBy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("last_updated", postureZone.LastUpdated)
	if err != nil {
		return diag.FromErr(err)
	}

	pIDs := make([]int, len(postureZone.Policies))
	for i, p := range postureZone.Policies {
		id, err := strconv.Atoi(p.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		pIDs[i] = id
	}
	err = d.Set("policy_ids", pIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	scopes := make([]map[string]any, len(postureZone.Scopes))
	for i, s := range postureZone.Scopes {
		scopes[i] = map[string]any{
			"target_type": s.TargetType,
			"rules":       s.Rules,
		}
	}
	err = d.Set("scopes", scopes)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

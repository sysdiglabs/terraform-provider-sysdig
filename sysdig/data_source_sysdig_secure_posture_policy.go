package sysdig

import (
	"context"
	"fmt"
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
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{SchemaIDKey, SchemaNameKey},
			},
			SchemaNameKey: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{SchemaIDKey, SchemaNameKey},
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

func dataSourceSysdigSecurePosturePolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	var policyID int64

	if idRaw, hasID := d.GetOk(SchemaIDKey); hasID {
		policyID, err = strconv.ParseInt(idRaw.(string), 10, 64)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid policy id: %s", err))
		}
	} else if nameRaw, hasName := d.GetOk(SchemaNameKey); hasName {
		name := nameRaw.(string)
		policies, listErr := client.ListPosturePolicies(ctx)
		if listErr != nil {
			return diag.FromErr(fmt.Errorf("error listing posture policies: %s", listErr))
		}
		var matchedID string
		for _, p := range policies {
			if p.Name == name {
				matchedID = p.ID
				break
			}
		}
		if matchedID == "" {
			return diag.FromErr(fmt.Errorf("posture policy with name %q not found", name))
		}
		policyID, err = strconv.ParseInt(matchedID, 10, 64)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid policy id %q: %s", matchedID, err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("either id or name must be specified"))
	}

	policy, err := client.GetPosturePolicyByID(ctx, policyID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(policy.ID)

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

	groupsData, err := setGroups(d, policy.RequirementsGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set(SchemaGroupKey, groupsData)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

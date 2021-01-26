package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func resourceSysdigSecurePolicyAssignmentBundle() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigSecurePolicyAssignmentBundleCreate,
		ReadContext:   resourceSysdigSecurePolicyAssignmentBundleRead,
		UpdateContext: resourceSysdigSecurePolicyAssignmentBundleUpdate,
		DeleteContext: resourceSysdigSecurePolicyAssignmentBundleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"default"}, false),
			},
			"policy_assignment": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registry": {
							Type:     schema.TypeString,
							Required: true,
							//ValidateFunc: validation.StringInSlice([]string{"stop", "pause", "kill"}, false),
						},
						"repository": {
							Type:     schema.TypeString,
							Required: true,
							//ValidateFunc: validation.StringInSlice([]string{"stop", "pause", "kill"}, false),
						},
						"tag": {
							Type:     schema.TypeString,
							Required: true,
							//ValidateFunc: validation.StringInSlice([]string{"stop", "pause", "kill"}, false),
						},
						"policy_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
						},
						"whitelist_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceSysdigSecurePolicyAssignmentBundleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecurePolicyAssignmentBundleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecurePolicyAssignmentBundleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diagz diag.Diagnostics

	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)

	if d.Id() == "" {
		d.SetId(name)
	}

	providerBundle, err := client.GetPolicyAssignmentBundleByName(ctx, name)
	policyAssignments := make([]interface{}, len(providerBundle.Items), len(providerBundle.Items))

	for i, item := range providerBundle.Items {
		bundleItem := make(map[string]interface{})

		bundleItem["registry"] = item.Registry
		bundleItem["id"] = item.ID
		bundleItem["repository"] = item.Repository

		// policy ids
		policyIds := []string{}
		for _, policy := range item.Policies {
			policyIds = append(policyIds, policy)
		}
		bundleItem["policy_ids"] = policyIds

		// whitelist ids
		whitelistIds := []string{}
		for _, item := range item.Whitelist {
			whitelistIds = append(whitelistIds, item)
		}
		bundleItem["whitelist_ids"] = whitelistIds

		bundleItem["tag"] = item.Image.Value

		policyAssignments[i] = bundleItem
	}

	if err = d.Set("policy_assignment", policyAssignments); err != nil {
		return diag.FromErr(err)
	}

	return diagz
}

func resourceSysdigSecurePolicyAssignmentBundleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diagz diag.Diagnostics

	items := d.Get("policy_assignment").([]interface{})
	policyItems := []secure.PolicyAssignment{}

	for _, item := range items {
		i := item.(map[string]interface{})
		policyItem := secure.PolicyAssignment{
			ID: i["id"].(string),
		}
		policyItems = append(policyItems, policyItem)
	}

	name := d.Get("name").(string)

	policyBundle := secure.PolicyAssignmentBundle{
		Items: policyItems,
		Id:    name,
	}

	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.PutPolicyAssignmentBundle(ctx, policyBundle)
	if err != nil {
		return diag.FromErr(err)
	}

	// be sure to read the resource back
	resourceSysdigSecurePolicyAssignmentBundleRead(ctx, d, meta)

	return diagz
}

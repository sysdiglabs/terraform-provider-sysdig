package sysdig

import (
	"context"
	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cast"
	"time"
)

func resourceSysdigSecurePolicyAssignments() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigSecurePolicyAssignmentsCreate,
		ReadContext:   resourceSysdigSecurePolicyAssignmentsRead,
		UpdateContext: resourceSysdigSecurePolicyAssignmentsUpdate,
		DeleteContext: resourceSysdigSecurePolicyAssignmentsDelete,
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
			"default_policy_assignment": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"whitelist_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
					},
				},
			},
			"policy_assignment": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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

func resourceSysdigSecurePolicyAssignmentsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSysdigSecurePolicyAssignmentsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("policy_assignment") {
		return resourceSysdigSecurePolicyAssignmentsCreate(ctx, d, meta)
	}
	return nil
}

func resourceSysdigSecurePolicyAssignmentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diagz diag.Diagnostics

	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	name := "default"
	d.SetId(name) // only default is possible

	providerBundle, err := client.GetPolicyAssignments(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	policyAssignments := make([]interface{}, len(providerBundle.Items), len(providerBundle.Items))
	defaultAssignment := make([]interface{}, 1, 1)

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

		if i == len(providerBundle.Items)-1 {
			// we are on the list item which must be default
			defaultItem := make(map[string]interface{})
			defaultItem["whitelist_ids"] = whitelistIds
			defaultItem["policy_ids"] = policyIds

			defaultAssignment[0] = defaultItem
		} else {
			policyAssignments[i] = bundleItem
		}

	}

	if err = d.Set("policy_assignment", policyAssignments); err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("default_policy_assignment", defaultAssignment); err != nil {
		return diag.FromErr(err)
	}

	return diagz
}

func resourceSysdigSecurePolicyAssignmentsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diagz diag.Diagnostics

	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	items := d.Get("policy_assignment").([]interface{})
	policyItems := []secure.PolicyAssignment{}

	for _, item := range items {
		i := item.(map[string]interface{})

		policyItem := secure.PolicyAssignment{
			ID:         i["id"].(string),
			Name:       "",
			Whitelist:  cast.ToStringSlice(i["whitelist_ids"].(*schema.Set).List()),
			Policies:   cast.ToStringSlice(i["policy_ids"].(*schema.Set).List()),
			Registry:   i["registry"].(string),
			Repository: i["repository"].(string),
			Image: secure.PolicyImage{
				Type:  "tag",
				Value: i["tag"].(string),
			},
		}
		policyItems = append(policyItems, policyItem)
	}

	// handle default assignment
	// lookup existing to preserve `default` id
	assignments, err := client.GetPolicyAssignments(ctx, "default")
	if err != nil {
		return diag.FromErr(err)
	}

	var defaultWhitelistIds []string
	var defaultPolicyIds []string
	var defaultId string

	if defaultItem, ok := d.GetOk("default_policy_assignment"); ok {
		defaultItemL := defaultItem.([]interface{})[0].(map[string]interface{})
		defaultWhitelistIds = cast.ToStringSlice(defaultItemL["whitelist_ids"].(*schema.Set).List())
		defaultPolicyIds = cast.ToStringSlice(defaultItemL["policy_ids"].(*schema.Set).List())
	} else {
		defaultFromAPI := assignments.Items[len(assignments.Items)-1]

		defaultWhitelistIds = defaultFromAPI.Whitelist
		defaultPolicyIds = defaultFromAPI.Policies
		defaultId = defaultFromAPI.ID
	}

	defaultPolicy := secure.PolicyAssignment{
		ID:         defaultId,
		Registry:   "*",
		Repository: "*",
		Image: secure.PolicyImage{
			Type:  "tag",
			Value: "*",
		},
		Whitelist: defaultWhitelistIds,
		Policies:  defaultPolicyIds,
	}

	policyItems = append(policyItems, defaultPolicy)

	name := "default"

	policyBundle := secure.PolicyAssignmentBundle{
		Items: policyItems,
		Id:    name,
	}

	_, err = client.PutPolicyAssignments(ctx, policyBundle)
	if err != nil {
		return diag.FromErr(err)
	}

	// be sure to read the resource back
	resourceSysdigSecurePolicyAssignmentsRead(ctx, d, meta)

	return diagz
}

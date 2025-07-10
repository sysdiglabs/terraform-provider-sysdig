package sysdig

import (
	"context"
	"errors"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func deprecatedResourceSysdigSecureScanningPolicyAssignment() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		DeprecationMessage: "The legacy scanning engine has been deprecated. This resource will be removed in future releases.",
		CreateContext:      deprecatedResourceSysdigScanningPolicyAssignmentCreate,
		ReadContext:        deprecatedResourceSysdigScanningPolicyAssignmentRead,
		UpdateContext:      deprecatedResourceSysdigScanningPolicyAssignmentUpdate,
		DeleteContext:      deprecatedResourceSysdigScanningPolicyAssignmentDelete,
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
			"items": { // todo: validate that at least there is one "default" with */*:*
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"image": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:             schema.TypeString,
										Optional:         true,
										Default:          "tag",
										ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"tag"}, false)),
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"policy_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"registry": {
							Type:     schema.TypeString,
							Required: true,
						},
						"repository": {
							Type:     schema.TypeString,
							Required: true,
						},
						"whitelist_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"policy_bundle_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
		},
	}
}

func getDeprecatedSecureScanningPolicyAssignmentClient(c SysdigClients) (v2.DeprecatedScanningPolicyAssignmentInterface, error) {
	return c.sysdigSecureClientV2()
}

func deprecatedResourceSysdigScanningPolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyAssignmentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentList := deprecatedScanningPolicyAssignmentListFromResourceData(d)

	validation := deprecatedValidateScanningPolicyAssignment(scanningPolicyAssignmentList)
	if validation != nil {
		return validation
	}

	scanningPolicyAssignmentList, err = client.CreateDeprecatedScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList)
	if err != nil {
		return diag.FromErr(err)
	}

	deprecatedScanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

func deprecatedResourceSysdigScanningPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyAssignmentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentList, err := client.GetDeprecatedScanningPolicyAssignmentList(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	deprecatedScanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

func deprecatedResourceSysdigScanningPolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyAssignmentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	scanningPolicyAssignmentList := deprecatedScanningPolicyAssignmentListFromResourceData(d)

	validation := deprecatedValidateScanningPolicyAssignment(scanningPolicyAssignmentList)
	if validation != nil {
		return validation
	}

	scanningPolicyAssignmentList, err = client.CreateDeprecatedScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList) // As policy assignments is a list, update is the same than create
	if err != nil {
		return diag.FromErr(err)
	}

	deprecatedScanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

// As Policy Assignments cannot be empty (default assignment cannot be deleted), pushing the default one
func deprecatedResourceSysdigScanningPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getDeprecatedSecureScanningPolicyAssignmentClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	defaultImage := v2.DeprecatedScanningPolicyAssignmentImage{
		Type:  "tag",
		Value: "*",
	}
	defaultItem := v2.DeprecatedScanningPolicyAssignment{
		Name:         "default",
		Registry:     "*",
		Repository:   "*",
		Image:        defaultImage,
		PolicyIDs:    []string{"default"},
		WhitelistIDs: []string{},
	}

	scanningPolicyAssignmentList := v2.DeprecatedScanningPolicyAssignmentList{
		PolicyBundleID: "default", // this is forced because there is no other possible value
		Items:          []v2.DeprecatedScanningPolicyAssignment{defaultItem},
	}

	err = client.DeleteDeprecatedScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deprecatedScanningPolicyAssignmentListToResourceData(scanningPolicyAssignmentList *v2.DeprecatedScanningPolicyAssignmentList, d *schema.ResourceData) {
	d.SetId(scanningPolicyAssignmentList.PolicyBundleID)
	_ = d.Set("policy_bundle_id", scanningPolicyAssignmentList.PolicyBundleID)
	var items []map[string]any

	for _, item := range scanningPolicyAssignmentList.Items {
		itemInfo := deprecatedScanningPolicyAssignmentToResourceData(item)

		items = append(items, itemInfo)
	}

	_ = d.Set("items", items)
}

func deprecatedScanningPolicyAssignmentToResourceData(scanningPolicyAssignment v2.DeprecatedScanningPolicyAssignment) map[string]any {
	item := map[string]any{
		"id":            scanningPolicyAssignment.ID,
		"name":          scanningPolicyAssignment.Name,
		"registry":      scanningPolicyAssignment.Registry,
		"repository":    scanningPolicyAssignment.Repository,
		"policy_ids":    scanningPolicyAssignment.PolicyIDs,
		"whitelist_ids": scanningPolicyAssignment.WhitelistIDs,
	}

	image := []map[string]any{{
		"type":  scanningPolicyAssignment.Image.Type,
		"value": scanningPolicyAssignment.Image.Value,
	}}

	item["image"] = image

	return item
}

func deprecatedScanningPolicyAssignmentListFromResourceData(d *schema.ResourceData) v2.DeprecatedScanningPolicyAssignmentList {
	scanningPolicyAssignmentList := v2.DeprecatedScanningPolicyAssignmentList{
		PolicyBundleID: "default", // this is forced because there is no other possible value
	}

	scanningPolicyAssignmentList.Items = deprecatedScanningPolicyAssignmentFromResourceData(d)

	return scanningPolicyAssignmentList
}

func deprecatedScanningPolicyAssignmentFromResourceData(d *schema.ResourceData) (scanningPolicyAssignmentItems []v2.DeprecatedScanningPolicyAssignment) {
	for _, item := range d.Get("items").([]any) {
		assignmentInfo := item.(map[string]any)
		assignment := v2.DeprecatedScanningPolicyAssignment{
			Name:       assignmentInfo["name"].(string),
			Registry:   assignmentInfo["registry"].(string),
			Repository: assignmentInfo["repository"].(string),
		}

		assignment.PolicyIDs = []string{}
		policyIDsSet := assignmentInfo["policy_ids"].([]any)
		for _, policy := range policyIDsSet {
			assignment.PolicyIDs = append(assignment.PolicyIDs, policy.(string))
		}

		assignment.WhitelistIDs = []string{}
		whitelistIDsSet := assignmentInfo["whitelist_ids"].([]any)
		for _, policy := range whitelistIDsSet {
			assignment.WhitelistIDs = append(assignment.WhitelistIDs, policy.(string))
		}
		imageSet := assignmentInfo["image"].([]any)
		if len(imageSet) == 0 {
			return
		}
		for _, image := range imageSet {
			assignment.Image = v2.DeprecatedScanningPolicyAssignmentImage{
				Type:  image.(map[string]any)["type"].(string),
				Value: image.(map[string]any)["value"].(string),
			}
		}

		scanningPolicyAssignmentItems = append(scanningPolicyAssignmentItems, assignment)
	}
	return scanningPolicyAssignmentItems
}

// Validate during creation as ValidateFunc is not supported in TypeList/TypeSet https://github.com/hashicorp/terraform-plugin-sdk/issues/156
// This function validates the last Item from the assignment list applies to all (*/*:*) and the list of policies is not empty in any assignment
func deprecatedValidateScanningPolicyAssignment(scanningPolicyAssignmentList v2.DeprecatedScanningPolicyAssignmentList) diag.Diagnostics {
	for _, item := range scanningPolicyAssignmentList.Items {
		if len(item.PolicyIDs) == 0 {
			return diag.FromErr(errors.New("'policy_ids' list can not be empty"))
		}
	}

	// validate default assignment
	lastItem := scanningPolicyAssignmentList.Items[len(scanningPolicyAssignmentList.Items)-1]
	if lastItem.Image.Value != "*" || lastItem.Registry != "*" || lastItem.Repository != "*" {
		return diag.FromErr(errors.New("default policy assignment has to be registry='*', repository='*' and image.tag='*?"))
	}

	return nil
}

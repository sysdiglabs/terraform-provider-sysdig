package sysdig

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
)

func resourceSysdigSecureScanningPolicyAssignment() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigScanningPolicyAssignmentCreate,
		ReadContext:   resourceSysdigScanningPolicyAssignmentRead,
		UpdateContext: resourceSysdigScanningPolicyAssignmentUpdate,
		DeleteContext: resourceSysdigScanningPolicyAssignmentDelete,
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

func resourceSysdigScanningPolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentList := scanningPolicyAssignmentListFromResourceData(d)

	validation := validateScanningPolicyAssignment(scanningPolicyAssignmentList)
	if validation != nil {
		return validation
	}

	scanningPolicyAssignmentList, err = client.CreateScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList)
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

func resourceSysdigScanningPolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentList, err := client.GetScanningPolicyAssignmentList(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

func resourceSysdigScanningPolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}
	scanningPolicyAssignmentList := scanningPolicyAssignmentListFromResourceData(d)

	validation := validateScanningPolicyAssignment(scanningPolicyAssignmentList)
	if validation != nil {
		return validation
	}

	scanningPolicyAssignmentList, err = client.CreateScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList) // As policy assignments is a list, update is the same than create
	if err != nil {
		return diag.FromErr(err)
	}

	scanningPolicyAssignmentListToResourceData(&scanningPolicyAssignmentList, d)

	return nil
}

// As Policy Assignments cannot be empty (default assignment cannot be deleted), pushing the default one
func resourceSysdigScanningPolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	defaultImage := secure.ScanningPolicyAssignmentImage{
		Type:  "tag",
		Value: "*",
	}
	defaultItem := secure.ScanningPolicyAssignment{
		Name:         "default",
		Registry:     "*",
		Repository:   "*",
		Image:        defaultImage,
		PolicyIDs:    []string{"default"},
		WhitelistIDs: []string{},
	}

	scanningPolicyAssignmentList := secure.ScanningPolicyAssignmentList{
		PolicyBundleId: "default", // this is forced because there is no other possible value
		Items:          []secure.ScanningPolicyAssignment{defaultItem},
	}

	err = client.DeleteScanningPolicyAssignmentList(ctx, scanningPolicyAssignmentList)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func scanningPolicyAssignmentListToResourceData(scanningPolicyAssignmentList *secure.ScanningPolicyAssignmentList, d *schema.ResourceData) {
	d.SetId(scanningPolicyAssignmentList.PolicyBundleId)
	_ = d.Set("policy_bundle_id", scanningPolicyAssignmentList.PolicyBundleId)
	var items []map[string]interface{}

	for _, item := range scanningPolicyAssignmentList.Items {
		itemInfo := scanningPolicyAssignmentToResourceData(item)

		items = append(items, itemInfo)
	}

	_ = d.Set("items", items)
}

func scanningPolicyAssignmentToResourceData(scanningPolicyAssignment secure.ScanningPolicyAssignment) map[string]interface{} {
	item := map[string]interface{}{
		"id":            scanningPolicyAssignment.ID,
		"name":          scanningPolicyAssignment.Name,
		"registry":      scanningPolicyAssignment.Registry,
		"repository":    scanningPolicyAssignment.Repository,
		"policy_ids":    scanningPolicyAssignment.PolicyIDs,
		"whitelist_ids": scanningPolicyAssignment.WhitelistIDs,
	}

	image := []map[string]interface{}{{
		"type":  scanningPolicyAssignment.Image.Type,
		"value": scanningPolicyAssignment.Image.Value,
	}}

	item["image"] = image

	return item
}

func scanningPolicyAssignmentListFromResourceData(d *schema.ResourceData) secure.ScanningPolicyAssignmentList {
	scanningPolicyAssignmentList := secure.ScanningPolicyAssignmentList{
		PolicyBundleId: "default", // this is forced because there is no other possible value
	}

	scanningPolicyAssignmentList.Items = scanningPolicyAssignmentFromResourceData(d)

	return scanningPolicyAssignmentList

}

func scanningPolicyAssignmentFromResourceData(d *schema.ResourceData) (scanningPolicyAssignmentItems []secure.ScanningPolicyAssignment) {
	for _, item := range d.Get("items").([]interface{}) {
		assignmentInfo := item.(map[string]interface{})
		assignment := secure.ScanningPolicyAssignment{
			Name:       assignmentInfo["name"].(string),
			Registry:   assignmentInfo["registry"].(string),
			Repository: assignmentInfo["repository"].(string),
		}

		assignment.PolicyIDs = []string{}
		policyIDsSet := assignmentInfo["policy_ids"].([]interface{})
		for _, policy := range policyIDsSet {
			assignment.PolicyIDs = append(assignment.PolicyIDs, policy.(string))
		}

		assignment.WhitelistIDs = []string{}
		whitelistIDsSet := assignmentInfo["whitelist_ids"].([]interface{})
		for _, policy := range whitelistIDsSet {
			assignment.WhitelistIDs = append(assignment.WhitelistIDs, policy.(string))
		}
		imageSet := assignmentInfo["image"].([]interface{})
		if len(imageSet) == 0 {
			return
		}
		for _, image := range imageSet {
			assignment.Image = secure.ScanningPolicyAssignmentImage{
				Type:  image.(map[string]interface{})["type"].(string),
				Value: image.(map[string]interface{})["value"].(string),
			}
		}

		scanningPolicyAssignmentItems = append(scanningPolicyAssignmentItems, assignment)
	}
	return scanningPolicyAssignmentItems
}

// Validate during creation as ValidateFunc is not supported in TypeList/TypeSet https://github.com/hashicorp/terraform-plugin-sdk/issues/156
// This function validates the last Item from the assignment list applies to all (*/*:*) and the list of policies is not empty in any assignment
func validateScanningPolicyAssignment(scanningPolicyAssignmentList secure.ScanningPolicyAssignmentList) diag.Diagnostics {
	for _, item := range scanningPolicyAssignmentList.Items {
		if len(item.PolicyIDs) == 0 {
			return diag.FromErr(errors.New("'policy_ids' list can not be empty"))
		}
	}

	// validate default assignment
	lastItem := scanningPolicyAssignmentList.Items[len(scanningPolicyAssignmentList.Items)-1]
	if lastItem.Image.Value != "*" || lastItem.Registry != "*" || lastItem.Repository != "*" {
		return diag.FromErr(errors.New("Default policy assignment has to be registry='*', repository='*' and image.tag='*?"))
	}

	return nil
}

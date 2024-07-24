package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createGroupSchema(i int) *schema.Resource {
	if i == 5 {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"description": {
					Type:     schema.TypeString,
					Required: true,
				},
				"requirement": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"id": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"name": {
								Type:     schema.TypeString,
								Required: true,
							},
							"description": {
								Type:     schema.TypeString,
								Required: true,
							},
							"control": {
								Type:     schema.TypeList,
								Optional: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": {
											Type:     schema.TypeString,
											Required: true,
										},
										"enabled": {
											Type:     schema.TypeBool,
											Optional: true,
											Default:  true,
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     createGroupSchema(i + 1),
			},
			"requirement": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Required: true,
						},
						"control": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceSysdigSecurePosturePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecurePosturePolicyCreateOrUpdate,
		ReadContext:   resourceSysdigSecurePosturePolicyRead,
		DeleteContext: resourceSysdigSecurePosturePolicyDelete,
		UpdateContext: resourceSysdigSecurePosturePolicyCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaTypeKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Unknown",
			},
			SchemaLinkKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaMinKubeVersionKey: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			SchemaMaxKubeVersionKey: {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			SchemaIsActiveKey: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			SchemaPlatformKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaVersionConstraintKey: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaMinVersionKey: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						SchemaMaxVersionKey: {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						SchemaPlatformKey: {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			SchemaGroupKey: {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     createGroupSchema(1),
			},
		},
	}
}

func resourceSysdigSecurePosturePolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Extract 'group' field from Terraform configuration
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	groups := extractGroupsRecursive(d.Get(SchemaGroupKey))
	req := &v2.CreatePosturePolicy{
		ID:                 getStringValue(d, SchemaIDKey),
		Name:               getStringValue(d, SchemaNameKey),
		Type:               getStringValue(d, SchemaTypeKey),
		Description:        getStringValue(d, SchemaDescriptionKey),
		MinKubeVersion:     getFloatValue(d, SchemaMinKubeVersionKey),
		MaxKubeVersion:     getFloatValue(d, SchemaMaxKubeVersionKey),
		IsActive:           getBoolValue(d, SchemaIsActiveKey),
		Platform:           getStringValue(d, SchemaPlatformKey),
		VersionConstraints: getVersionConstraintsValue(d, SchemaVersionConstraintKey),
		Link:               getStringValue(d, SchemaLinkKey),
		RequirementGroups:  groups,
	}

	new, errStatus, err := client.CreateOrUpdatePosturePolicy(ctx, req)

	if err != nil {
		return diag.Errorf("Error creating new policy with groups. error status: %s err: %s", errStatus, err)
	}

	d.SetId(new.ID)
	resourceSysdigSecurePosturePolicyRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecurePosturePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
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

	err = setVersionConstraints(d, SchemaVersionConstraintKey, policy.VersionConstraints)

	if err != nil {
		return diag.FromErr(err)
	}
	// Set groups
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

func resourceSysdigSecurePosturePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeletePosturePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setGroups(d *schema.ResourceData, groups []v2.RequirementsGroup) ([]interface{}, error) {
	var groupsData []interface{}
	for _, group := range groups {
		groupData := map[string]interface{}{
			"id":          group.ID,
			"name":        group.Name,
			"description": group.Description,
		}

		// Recursively set nested groups and requirements
		if len(group.Requirements) > 0 {
			requirementsData := setRequirements(group.Requirements)
			groupData["requirement"] = requirementsData
		}
		if len(group.Folders) > 0 {
			nestedGroupsData, err := setGroups(d, group.Folders)
			if err != nil {
				return nil, err
			}
			groupData["group"] = nestedGroupsData
		}
		groupsData = append(groupsData, groupData)
	}
	return groupsData, nil
}

func setRequirements(requirements []v2.Requirement) []interface{} {
	var requirementsData []interface{}
	for _, req := range requirements {
		reqData := map[string]interface{}{
			"id":          req.ID,
			"name":        req.Name,
			"description": req.Description,
		}

		// Set controls for each requirement
		if len(req.Controls) > 0 {
			controlsData := setControls(req.Controls)
			reqData["control"] = controlsData
		}

		requirementsData = append(requirementsData, reqData)
	}
	return requirementsData
}

func setControls(controls []v2.Control) []interface{} {
	var controlsData []interface{}
	for _, ctrl := range controls {
		ctrlData := map[string]interface{}{
			"name":    ctrl.Name,
			"enabled": ctrl.Status,
		}
		controlsData = append(controlsData, ctrlData)
	}
	return controlsData
}

// Helper function to retrieve string value from ResourceData and handle nil case
func getStringValue(d *schema.ResourceData, key string) string {
	if value, ok := d.GetOk(key); ok {
		return value.(string)
	}
	return ""
}

// Helper function to retrieve version constraints value from ResourceData and handle nil case
func getVersionConstraintsValue(d *schema.ResourceData, key string) []v2.VersionConstraint {
	pvc := []v2.VersionConstraint{}
	versionContraintsMap, ok := d.Get(key).([]interface{})
	if !ok {
		return nil
	}
	for _, vc := range versionContraintsMap {
		vcMap := vc.(map[string]interface{})
		minVersion := 0.0
		maxVersion := 0.0
		if vcMap["min_version"] != nil {
			minVersion = vcMap["min_version"].(float64)
		}
		if vcMap["max_version"] != nil {
			maxVersion = vcMap["max_version"].(float64)
		}
		versionConstraint := v2.VersionConstraint{
			MinVersion: minVersion,
			MaxVersion: maxVersion,
			Platform:   vcMap["platform"].(string),
		}
		pvc = append(pvc, versionConstraint)
	}
	return pvc
}

// Helper function to retrieve float64 value from ResourceData and handle nil case
func getFloatValue(d *schema.ResourceData, key string) float64 {
	if value, ok := d.GetOk(key); ok {
		return value.(float64)
	}
	return 0
}

// Helper function to retrieve bool value from ResourceData and handle nil case
func getBoolValue(d *schema.ResourceData, key string) bool {
	if value, ok := d.GetOk(key); ok {
		return value.(bool)
	}
	return false
}

func extractGroupsRecursive(data interface{}) []v2.CreateRequirementsGroup {
	var groups []v2.CreateRequirementsGroup

	switch d := data.(type) {
	case []interface{}:
		for _, item := range d {
			groups = append(groups, extractGroupsRecursive(item)...)
		}
	case map[string]interface{}:
		group := v2.CreateRequirementsGroup{
			ID:          d["id"].(string),
			Name:        d["name"].(string),
			Description: d["description"].(string),
		}

		if reqs, ok := d["requirement"].([]interface{}); ok {
			for _, reqData := range reqs {
				reqMap := reqData.(map[string]interface{})
				requirement := v2.CreateRequirement{
					ID:          reqMap["id"].(string),
					Name:        reqMap["name"].(string),
					Description: reqMap["description"].(string),
				}

				if controlsData, ok := reqMap["control"].([]interface{}); ok {
					for _, controlData := range controlsData {
						controlMap := controlData.(map[string]interface{})
						control := v2.CreateRequirementControl{
							Name:    controlMap["name"].(string),
							Enabled: controlMap["enabled"].(bool),
						}
						requirement.Controls = append(requirement.Controls, control)
					}
				}

				group.Requirements = append(group.Requirements, requirement)
			}
		}

		if subGroups, ok := d["group"]; ok {
			group.Folders = extractGroupsRecursive(subGroups)
		}

		groups = append(groups, group)
	}

	return groups
}

// Helper function to set version constraints in the Terraform schema
func setVersionConstraints(d *schema.ResourceData, key string, constraints []v2.VersionConstraint) error {
	var constraintsData []interface{}
	for _, vc := range constraints {
		constraint := map[string]interface{}{
			"min_version": vc.MinVersion,
			"max_version": vc.MaxVersion,
			"platform":    vc.Platform,
		}
		constraintsData = append(constraintsData, constraint)
	}
	if err := d.Set(key, constraintsData); err != nil {
		return err
	}
	return nil
}

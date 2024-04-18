package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecurePosturePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecurePosturePolicyCreate,
		ReadContext:   resourceSysdigSecurePosturePolicyRead,
		DeleteContext: resourceSysdigSecurePosturePolicyDelete,
		UpdateContext: resourceSysdigSecurePosturePolicyUpdate,
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
			SchemaTypeKey: {
				Type:     schema.TypeInt,
				Optional: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaLinkKey: {
				Type:     schema.TypeString,
				Optional: true,
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
			},
			SchemaPlatformKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			SchemaGroupsKey: { // level1
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
						"requirements": {
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
									"controls": {
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
												},
											},
										},
									},
								},
							},
						},
						"groups": { // level2
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Required: true,
									},
									"requirements": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"description": {
													Type:     schema.TypeString,
													Required: true,
												},
												"controls": {
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
															},
														},
													},
												},
											},
										},
									},
									"groups": { //level3
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"description": {
													Type:     schema.TypeString,
													Required: true,
												},
												"requirements": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
															"description": {
																Type:     schema.TypeString,
																Required: true,
															},
															"controls": {
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
																		},
																	},
																},
															},
														},
													},
												},
												"groups": { // level4
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:     schema.TypeString,
																Required: true,
															},
															"description": {
																Type:     schema.TypeString,
																Required: true,
															},
															"requirements": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"description": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"controls": {
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
																					},
																				},
																			},
																		},
																	},
																},
															},
															"groups": {
																Type:     schema.TypeList,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"name": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"description": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"requirements": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": {
																						Type:     schema.TypeString,
																						Required: true,
																					},
																					"description": {
																						Type:     schema.TypeString,
																						Required: true,
																					},
																					"controls": {
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
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
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

func resourceSysdigSecurePosturePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Extract 'groups' field from Terraform configuration
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	groups := extractGroupsRecursive(d.Get(SchemaGroupsKey))

	req := &v2.CreatePosturePolicy{
		Name:               getStringValue(d, SchemaNameKey),
		Description:        getStringValue(d, SchemaDescriptionKey),
		MinKubeVersion:     getFloatValue(d, SchemaMinKubeVersionKey),
		MaxKubeVersion:     getFloatValue(d, SchemaMaxKubeVersionKey),
		IsActive:           getBoolValue(d, SchemaIsActiveKey),
		Platform:           getStringValue(d, SchemaPlatformKey),
		Link:               getStringValue(d, SchemaLinkKey),
		RequirementFolders: groups,
	}
	fmt.Println("requestttttttttttttttttttttttttttttttttttttttttt: ", req)
	new, errStatus, err := client.CreatePosturePolicy(ctx, req)
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

	fmt.Println("get policyyyyyyyyyyyyyyyyyyyyyyyy: ", policy)

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

	return nil

}

func resourceSysdigSecurePosturePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil

}

func resourceSysdigSecurePosturePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

// Helper function to retrieve string value from ResourceData and handle nil case
func getStringValue(d *schema.ResourceData, key string) string {
	if value, ok := d.GetOk(key); ok {
		return value.(string)
	}
	return ""
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
			Name:        d["name"].(string),
			Description: d["description"].(string),
		}

		if reqs, ok := d["requirements"].([]interface{}); ok {
			for _, reqData := range reqs {
				reqMap := reqData.(map[string]interface{})
				requirement := v2.CreateRequirement{
					Name:        reqMap["name"].(string),
					Description: reqMap["description"].(string),
				}

				if controlsData, ok := reqMap["controls"].([]interface{}); ok {
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

		if subGroups, ok := d["groups"]; ok {
			group.Folders = extractGroupsRecursive(subGroups)
		}

		groups = append(groups, group)
	}

	return groups
}

package sysdig

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	// debug
	// "fmt"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureCompositePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigCompositePolicyCreate,
		ReadContext:   resourceSysdigCompositePolicyRead,
		UpdateContext: resourceSysdigCompositePolicyUpdate,
		DeleteContext: resourceSysdigCompositePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSysdigSecureCompositePolicyImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},
		// IMPORTANT: composite.Policy.Rules and composite.Policy.RuleNames are read-only attributes.
		// They're not used as a source for rule updates, so it's okay to drop those attributes in TF.
		// To update the rules, composite.Rules values are used instead.
		// https://github.com/draios/secure-backend/blob/main/policies/api/handler_policies.go#L1120
		Schema: map[string]*schema.Schema{
			// IMPORTANT: Type is implicit: It's automatically added upon conversion to JSON
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "malware",
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"malware"}, false)),
			},
			"name":                  NameSchema(),
			"description":           DescriptionSchema(),
			"enabled":               EnabledSchema(),
			"severity":              SeveritySchema(),
			"scope":                 ScopeSchema(),
			"version":               VersionSchema(),
			"notification_channels": NotificationChannelsSchema(),
			"runbook":               RunbookSchema(),
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   ReadOnlyIntSchema(),
						"name": ReadOnlyStringSchema(),
						// Do not allow switching off individual rules
						// "enabled":     EnabledSchema(),
						"description": DescriptionSchema(),
						"tags":        TagsSchema(),
						"details": {
							Type:     schema.TypeList,
							MaxItems: 1, // There can only ever be a single details block per rule
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_managed_hashes": BoolSchema(),
									"additional_hashes":  HashesSchema(),
									"ignore_hashes":      HashesSchema(),
								},
							},
						},
					},
				},
			},
			"actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prevent_malware": PreventMalwareActionSchema(),
						"container":       ContainerActionSchema(),
						"capture":         CaptureActionSchema(),
					},
				},
			},
		}, // Schema end
	}
}

func getSecureCompositePolicyClient(c SysdigClients) (v2.CompositePolicyInterface, error) {
	return c.sysdigSecureClientV2()
}

// TODO: Move
func schemaSetToList(values interface{}) []string {
	v := values.(*schema.Set).List()

	x := make([]string, len(v))
	for i := range v {
		x[i] = v[i].(string)
	}
	return x
}

func compositePolicyFromResourceData(d *schema.ResourceData) (v2.PolicyRulesComposite, error) {
	policy := &v2.PolicyRulesComposite{
		Policy: &v2.Policy{},
		Rules:  []*v2.RuntimePolicyRule{},
	}
	err := malwarePolicyReducer(policy, d)
	if err != nil {
		return *policy, err
	}

	return *policy, nil
}

// TODO: Rename
func compositePolicyToResourceData(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	return malwareTFResourceReducer(d, *policy)
}

func resourceSysdigCompositePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := compositePolicyFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err = client.CreateCompositePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	compositePolicyToResourceData(&policy, d)

	return nil
}

func resourceSysdigCompositePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := compositePolicyFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCompositePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigCompositePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	policy, statusCode, err := client.GetCompositePolicyByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	err = compositePolicyToResourceData(&policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigCompositePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := compositePolicyFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if policy.Policy.ID == 0 {
		return diag.FromErr(errors.New("Policy ID is missing"))
	}

	err = client.DeleteCompositePolicy(ctx, policy.Policy.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigSecureCompositePolicyImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return nil, err
	}

	policy, err := compositePolicyFromResourceData(d)
	if err != nil {
		return nil, err
	}

	if policy.Policy.ID == 0 {
		return nil, errors.New("Policy ID is missing")
	}

	policy, _, err = client.GetCompositePolicyByID(ctx, policy.Policy.ID)
	if err != nil {
		return nil, err
	}

	if policy.Policy.IsDefault || policy.Policy.TemplateId != 0 {
		return nil, errors.New("unable to import policy that is not a custom policy")
	}

	err = compositePolicyToResourceData(&policy, d)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

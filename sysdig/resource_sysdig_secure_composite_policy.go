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
						"id":          ReadOnlyIntSchema(),
						"name":        ReadOnlyStringSchema(),
						"enabled":     EnabledSchema(),
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

func schemaSetToList(values interface{}) []string {
	v := values.(*schema.Set).List()

	x := make([]string, len(v))
	for i := range v {
		x[i] = v[i].(string)
	}
	return x
}

func compositePolicyFromResourceData(d *schema.ResourceData) v2.PolicyRulesComposite {
	policy := &v2.PolicyRulesComposite{
		Policy: &v2.Policy{},
		Rules:  []*v2.RuntimePolicyRule{},
	}
	commonPolicyFromResourceData(policy.Policy, d)

	policy.Policy.Description = d.Get("description").(string)
	policy.Policy.Severity = d.Get("severity").(int)
	// policy.Policy.Version = ???

	policy.Policy.Rules = []*v2.PolicyRule{}
	policy.Rules = []*v2.RuntimePolicyRule{}
	if _, ok := d.GetOk("rules"); ok {
		// TODO: Iterate over a list of rules instead of hard-coding the index values
		// TODO: Should we assume that only a single Malware rule can be attached to a policy?

		policy.Policy.Type = "malware"
		// TODO: What origin should we use?
		// https://github.com/draios/secure-backend/blob/main/policies/model/model.go#L1576
		// policy.Policy.Origin = "Terraform"

		additionalHashes := map[string][]string{}
		if items, ok := d.GetOk("rules.0.details.0.additional_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]interface{}) {
				item := item.(map[string]interface{})
				k := item["hash"].(string)
				v := schemaSetToList(item["hash_aliases"])
				additionalHashes[k] = v
			}
		}

		// TODO: Extract into a function
		ignoreHashes := map[string][]string{}
		if items, ok := d.GetOk("rules.0.details.0.ignore_hashes"); ok { // TODO: Do not hardcode the indexes
			for _, item := range items.([]interface{}) {
				item := item.(map[string]interface{})
				k := item["hash"].(string)
				v := schemaSetToList(item["hash_aliases"])
				ignoreHashes[k] = v
			}
		}

		tags := schemaSetToList(d.Get("rules.0.tags"))
		rule := &v2.RuntimePolicyRule{
			// TODO: Do not hardcode the indexes
			Name:        d.Get("rules.0.name").(string),
			Description: d.Get("rules.0.description").(string),
			Tags:        tags,
			Details: v2.MalwareRuleDetails{
				RuleType:         v2.ElementType("MALWARE"), // TODO: Use const
				UseManagedHashes: d.Get("rules.0.details.0.use_managed_hashes").(bool),
				AdditionalHashes: additionalHashes,
				IgnoreHashes:     ignoreHashes,
			},
		}

		id := v2.FlexInt(d.Get("rules.0.id").(int))
		if int(id) != 0 {
			rule.Id = &id
		} else {
			// TODO: Panic?
			// panic(fmt.Sprintf("id is nil: %s, %s", d.Get("rules.0.name"), d.Get("rules.0.id")))
		}
		policy.Rules = append(policy.Rules, rule)

		policy.Policy.Rules = append(policy.Policy.Rules, &v2.PolicyRule{
			Name:    d.Get("rules.0.name").(string),
			Enabled: d.Get("rules.0.enabled").(bool),
		})

		return *policy
	}

	// TODO: Add other types: ML, AWS_ML, DRIFT, etc.

	return *policy
}

func compositePolicyToResourceData(policy *v2.PolicyRulesComposite, d *schema.ResourceData) {
	commonPolicyToResourceData(policy.Policy, d)

	_ = d.Set("description", policy.Policy.Description)
	_ = d.Set("severity", policy.Policy.Severity)
	if policy.Policy.Type != "" {
		_ = d.Set("type", policy.Policy.Type)
	} else {
		// _ = d.Set("type", "malware") // TODO
	}

	actions := compositePolicyDataSourceActionsToResourceData(policy.Policy.Actions)
	d.Set("actions", actions)

	enabledByRuleName := map[string]bool{}
	for _, rule := range policy.Policy.Rules {
		enabledByRuleName[rule.Name] = rule.Enabled
	}

	if len(policy.Rules) > 0 {
		// TODO: Extract into a function
		rules := []map[string]interface{}{}

		for _, r := range policy.Rules {
			// TODO: Single element only
			additionalHashes := []map[string]interface{}{}
			for k, v := range r.Details.(*v2.MalwareRuleDetails).AdditionalHashes {
				additionalHashes = append(additionalHashes, map[string]interface{}{
					"hash":         k,
					"hash_aliases": v,
				})
			}

			ignoreHashes := []map[string]interface{}{}
			for k, v := range r.Details.(*v2.MalwareRuleDetails).IgnoreHashes {
				ignoreHashes = append(ignoreHashes, map[string]interface{}{
					"hash":         k,
					"hash_aliases": v,
				})
			}

			details := []map[string]interface{}{{}}
			details[0] = map[string]interface{}{
				"use_managed_hashes": r.Details.(*v2.MalwareRuleDetails).UseManagedHashes,
				"additional_hashes":  additionalHashes,
				"ignore_hashes":      ignoreHashes,
			}

			rules = append(rules, map[string]interface{}{
				"id":          *policy.Rules[0].Id,
				"name":        policy.Rules[0].Name,
				"enabled":     enabledByRuleName[policy.Rules[0].Name],
				"description": policy.Rules[0].Description,
				"tags":        policy.Rules[0].Tags,
				"details":     details,
			})

		}

		_ = d.Set("rules", rules)
	}
	// TODO: Add other policy rule type: ML, AWS ML, DRIFT, etc.
}

func validateCompositePolicy(policy *v2.PolicyRulesComposite) error {
	// TODO: Add other validation rules
	if len(policy.Rules) == 0 {
		return errors.New("Policy is missing rules")
	}
	return nil
}

func resourceSysdigCompositePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := compositePolicyFromResourceData(d)
	if err = validateCompositePolicy(&policy); err != nil {
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

	policy := compositePolicyFromResourceData(d)
	if err = validateCompositePolicy(&policy); err != nil {
		return diag.FromErr(err)
	}
	policy.Policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.Policy.ID = id

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

	compositePolicyToResourceData(&policy, d)

	return nil
}

func resourceSysdigCompositePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeleteCompositePolicy(ctx, id)
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

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, err
	}

	policy, _, err := client.GetCompositePolicyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if policy.Policy.IsDefault || policy.Policy.TemplateId != 0 {
		return nil, errors.New("unable to import policy that is not a custom policy")
	}

	return []*schema.ResourceData{d}, nil
}

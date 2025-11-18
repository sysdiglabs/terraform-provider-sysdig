package sysdig

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureOktaMLPolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigOktaMLPolicyCreate,
		ReadContext:   resourceSysdigOktaMLPolicyRead,
		UpdateContext: resourceSysdigOktaMLPolicyUpdate,
		DeleteContext: resourceSysdigOktaMLPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSysdigSecureOktaMLPolicyImportState,
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
				Default:          policyTypeOktaML,
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{policyTypeOktaML}, false)),
			},
			"name":                  NameSchema(),
			"description":           DescriptionSchema(),
			"enabled":               EnabledSchema(),
			"severity":              SeveritySchema(),
			"scope":                 ScopeSchema(),
			"version":               VersionSchema(),
			"notification_channels": NotificationChannelsSchema(),
			"runbook":               RunbookSchema(),
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   ReadOnlyIntSchema(),
						"name": ReadOnlyStringSchema(),
						// Do not allow switching off individual rules
						// "enabled":     EnabledSchema(),
						"description":             DescriptionSchema(),
						"tags":                    TagsSchema(),
						"version":                 VersionSchema(),
						"anomalous_console_login": MLRuleThresholdAndSeveritySchema(),
					},
				},
			},
		}, // Schema end
	}
}

func oktaMLPolicyFromResourceData(d *schema.ResourceData) (v2.PolicyRulesComposite, error) {
	policy := v2.PolicyRulesComposite{
		Policy: &v2.Policy{},
		Rules:  []*v2.RuntimePolicyRule{},
	}
	err := oktaMLPolicyReducer(&policy, d)
	if err != nil {
		return policy, err
	}

	return policy, nil
}

func oktaMLPolicyToResourceData(policy *v2.PolicyRulesComposite, d *schema.ResourceData) error {
	return oktaMLTFResourceReducer(d, *policy)
}

func resourceSysdigOktaMLPolicyCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := oktaMLPolicyFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err = client.CreateCompositePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	err = oktaMLPolicyToResourceData(&policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigOktaMLPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := oktaMLPolicyFromResourceData(d)
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

func resourceSysdigOktaMLPolicyRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policy, statusCode, err := client.GetCompositePolicyByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	err = oktaMLPolicyToResourceData(&policy, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigOktaMLPolicyDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureCompositePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCompositePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigSecureOktaMLPolicyImportState(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	client, err := getSecureCompositePolicyClient(meta.(SysdigClients))
	if err != nil {
		return nil, err
	}

	policy, err := oktaMLPolicyFromResourceData(d)
	if err != nil {
		return nil, err
	}

	if policy.Policy.ID == 0 {
		return nil, errors.New("policy ID is missing")
	}

	policy, _, err = client.GetCompositePolicyByID(ctx, policy.Policy.ID)
	if err != nil {
		return nil, err
	}

	if policy.Policy.IsDefault || policy.Policy.TemplateID != 0 {
		return nil, errors.New("unable to import policy that is not a custom policy")
	}

	err = oktaMLPolicyToResourceData(&policy, d)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

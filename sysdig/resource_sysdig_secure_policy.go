package sysdig

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var validatePolicyType = validation.StringInSlice([]string{
	"falco",
	"list_matching",
	"k8s_audit",
	"aws_cloudtrail",
	"gcp_auditlog",
	"azure_platformlogs",
  "okta",
  "github",
	"malware",
	"drift",
	"aws_machine_learning",
	"machine_learning",
}, false)

func resourceSysdigSecurePolicy() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigPolicyCreate,
		ReadContext:   resourceSysdigPolicyRead,
		UpdateContext: resourceSysdigPolicyUpdate,
		DeleteContext: resourceSysdigPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
		},

		DeprecationMessage: "The sysdig_secure_policy resource is being replaced by sysdig_secure_custom_policy, " +
			"sysdig_secure_managed_policy, and sysdig_secure_managed_ruleset depending on the type of policy.",

		Schema: createPolicySchema(map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "falco",
				ValidateDiagFunc: validateDiagFunc(validatePolicyType),
			},
			"severity": {
				Type:             schema.TypeInt,
				Default:          4,
				Optional:         true,
				ValidateDiagFunc: validateDiagFunc(validation.IntBetween(0, 7)),
			},
			"rule_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		}),
	}
}

func getSecurePolicyClient(c SysdigClients) (v2.PolicyInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := policyFromResourceData(d)
	policy, err = client.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	policyToResourceData(&policy, d)

	return nil
}

// Saves the resource data information for the common fields of the policy
func commonPolicyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	if policy.ID != 0 {
		d.SetId(strconv.Itoa(policy.ID))
	}

	_ = d.Set("name", policy.Name)
	_ = d.Set("scope", policy.Scope)
	_ = d.Set("enabled", policy.Enabled)
	_ = d.Set("version", policy.Version)
	_ = d.Set("runbook", policy.Runbook)

	actions := []map[string]interface{}{{}}
	for _, action := range policy.Actions {
		if action.Type != "POLICY_ACTION_CAPTURE" {
			action := strings.Replace(action.Type, "POLICY_ACTION_", "", 1)
			actions[0]["container"] = strings.ToLower(action)
			// d.Set("actions.0.container", strings.ToLower(action))
		} else {
			actions[0]["capture"] = []map[string]interface{}{{
				"seconds_after_event":  action.AfterEventNs / 1000000000,
				"seconds_before_event": action.BeforeEventNs / 1000000000,
				"name":                 action.Name,
				"filter":               action.Filter,
				"bucket_name":          action.BucketName,
				"folder":               action.Folder,
			}}
		}
	}

	currentContainerAction := d.Get("actions.0.container").(string)
	currentCaptureAction := d.Get("actions.0.capture").([]interface{})
	// If the policy retrieved from service has no actions and the current state is default values,
	// then do not set the "actions" key as it may cause terraform to think there has been a state change
	if len(policy.Actions) > 0 || currentContainerAction != "" || len(currentCaptureAction) > 0 {
		_ = d.Set("actions", actions)
	}

	_ = d.Set("notification_channels", policy.NotificationChannelIds)
}

func policyToResourceData(policy *v2.Policy, d *schema.ResourceData) {
	commonPolicyToResourceData(policy, d)

	_ = d.Set("description", policy.Description)
	_ = d.Set("severity", policy.Severity)
	if policy.Type != "" {
		_ = d.Set("type", policy.Type)
	} else {
		_ = d.Set("type", "falco")
	}

	_ = d.Set("rule_names", policy.RuleNames)
}

func commonPolicyFromResourceData(policy *v2.Policy, d *schema.ResourceData) {
	policy.Name = d.Get("name").(string)
	policy.Enabled = d.Get("enabled").(bool)
	policy.Runbook = d.Get("runbook").(string)
	policy.Scope = d.Get("scope").(string)

	addActionsToPolicy(d, policy)

	policy.NotificationChannelIds = []int{}
	notificationChannelIdSet := d.Get("notification_channels").(*schema.Set)
	for _, id := range notificationChannelIdSet.List() {
		policy.NotificationChannelIds = append(policy.NotificationChannelIds, id.(int))
	}
}

func policyFromResourceData(d *schema.ResourceData) v2.Policy {
	policy := &v2.Policy{}
	commonPolicyFromResourceData(policy, d)

	policy.Description = d.Get("description").(string)
	policy.Severity = d.Get("severity").(int)
	policy.Type = d.Get("type").(string)

	policy.RuleNames = []string{}
	rule_names := d.Get("rule_names").(*schema.Set)
	for _, name := range rule_names.List() {
		if rule_name, ok := name.(string); ok {
			rule_name = strings.TrimSpace(rule_name)
			policy.RuleNames = append(policy.RuleNames, rule_name)
		}
	}

	return *policy
}

func addActionsToPolicy(d *schema.ResourceData, policy *v2.Policy) {
	policy.Actions = []v2.Action{}
	actions := d.Get("actions").([]interface{})
	if len(actions) == 0 {
		return
	}

	preventMalwareAction, ok := d.GetOkExists("actions.0.prevent_malware")
	if ok && preventMalwareAction.(bool) {
		policy.Actions = append(policy.Actions, v2.Action{Type: "POLICY_ACTION_PREVENT_MALWARE"})
	}

	containerAction := d.Get("actions.0.container").(string)
	if containerAction != "" {
		containerAction = strings.ToUpper("POLICY_ACTION_" + containerAction)

		policy.Actions = append(policy.Actions, v2.Action{Type: containerAction})
	}

	if captureAction := d.Get("actions.0.capture").([]interface{}); len(captureAction) > 0 {
		afterEventNs := d.Get("actions.0.capture.0.seconds_after_event").(int) * 1000000000
		beforeEventNs := d.Get("actions.0.capture.0.seconds_before_event").(int) * 1000000000
		name := d.Get("actions.0.capture.0.name").(string)
		filter := d.Get("actions.0.capture.0.filter").(string)
		bucketName := d.Get("actions.0.capture.0.bucket_name").(string)
		folder := d.Get("actions.0.capture.0.folder").(string)
		policy.Actions = append(policy.Actions, v2.Action{
			Type:                 "POLICY_ACTION_CAPTURE",
			IsLimitedToContainer: false,
			AfterEventNs:         afterEventNs,
			BeforeEventNs:        beforeEventNs,
			Name:                 name,
			Filter:               filter,
			StorageType:          "S3",
			BucketName:           bucketName,
			Folder:               folder,
		})
	}
}

func resourceSysdigPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecurePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())
	policy, statusCode, err := client.GetPolicyByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	policyToResourceData(&policy, d)

	return nil
}

func resourceSysdigPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := strconv.Atoi(d.Id())

	err = client.DeletePolicy(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecurePolicyClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	policy := policyFromResourceData(d)
	policy.Version = d.Get("version").(int)

	id, _ := strconv.Atoi(d.Id())
	policy.ID = id

	_, err = client.UpdatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

var sendPoliciesToAgentsOnce sync.Once

func sendPoliciesToAgents(ctx context.Context, clients SysdigClients) error {
	var err error
	sendPoliciesToAgentsOnce.Do(func() {
		tflog.Info(ctx, "Sending policies to agents")
		var client v2.PolicyInterface
		client, err = getSecurePolicyClient(clients)
		if err != nil {
			return
		}

		// When running as a cleanup hook, the terraform context is in a cancelled state.
		// Using a background context with a deadline will allow us to complete this request.
		backgroundCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
		defer cancel()
		err = client.SendPoliciesToAgents(backgroundCtx)
	})
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error in sendPoliciesToAgents: %s", err.Error()))
	}
	return err
}

package sysdig

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/spf13/cast"
)

var validateFalcoRuleSource = validation.StringInSlice([]string{"syscall", "k8s_audit", "aws_cloudtrail", "gcp_auditlog", "azure_platformlogs", "awscloudtrail", "okta", "github"}, false)

func resourceSysdigSecureRuleFalco() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigRuleFalcoCreate,
		UpdateContext: resourceSysdigRuleFalcoUpdate,
		ReadContext:   resourceSysdigRuleFalcoRead,
		DeleteContext: resourceSysdigRuleFalcoDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"priority": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "warning",
				ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{"emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"}, false)),
			},
			"source": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "",
				ValidateDiagFunc: validateDiagFunc(validateFalcoRuleSource),
			},
			"minimum_engine_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"append": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"exceptions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"comps": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		}),
	}
}

func resourceSysdigRuleFalcoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleFalcoFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err = client.CreateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	d.SetId(strconv.Itoa(rule.ID))
	_ = d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleFalcoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getSecureRuleClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, statusCode, err := client.GetRuleByID(ctx, id)
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}

	if rule.Details.Append != nil && !(*(rule.Details.Append)) {
		if rule.Details.Condition == nil {
			return diag.Errorf("no condition data for a falco rule")
		}
	}

	updateResourceDataForRule(d, rule)
	if rule.Details.Condition != nil {
		_ = d.Set("condition", rule.Details.Condition.Condition)
	}
	_ = d.Set("output", rule.Details.Output)
	_ = d.Set("priority", strings.ToLower(rule.Details.Priority))
	_ = d.Set("source", rule.Details.Source)
	if rule.Details.MinimumEngineVersion != nil {
		_ = d.Set("minimum_engine_version", *rule.Details.MinimumEngineVersion)
	}
	if rule.Details.Append != nil {
		_ = d.Set("append", *rule.Details.Append)
	}
	if err := updateResourceDataExceptions(d, rule.Details.Exceptions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateResourceDataExceptions(d *schema.ResourceData, ruleExceptions []*v2.Exception) error {
	exceptions := make([]any, 0, len(ruleExceptions))
	for _, exception := range ruleExceptions {
		if exception == nil {
			return fmt.Errorf("exception is nil")
		}
		valuesData, err := json.Marshal(exception.Values)
		if err != nil {
			return fmt.Errorf("error marshalling exception values '%+v': %s", exception.Values, err)
		}
		fields, err := fieldOrCompsToStringSlice(exception.Fields)
		if err != nil {
			return fmt.Errorf("error converting exception fields '%+v': %s", exception.Fields, err)
		}
		comps, err := fieldOrCompsToStringSlice(exception.Comps)
		if err != nil {
			return fmt.Errorf("error converting exception comps '%+v': %s", exception.Comps, err)
		}

		exceptions = append(exceptions, map[string]any{
			"name":   exception.Name,
			"comps":  comps,
			"values": string(valuesData),
			"fields": fields,
		})
	}
	_ = d.Set("exceptions", exceptions)
	return nil
}

func fieldOrCompsToStringSlice(fields any) ([]string, error) {
	elements := []string{}
	switch t := fields.(type) {
	case []interface{}:
		for _, field := range t {
			elements = append(elements, field.(string))
		}
	case string:
		elements = append(elements, t)
	case nil:
		// do nothing
	default:
		return nil, fmt.Errorf("unexpected type: %T", t)
	}
	return elements, nil
}

func resourceSysdigRuleFalcoUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleFalcoFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleFalcoDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sysdigClients := meta.(SysdigClients)
	client, err := getSecureRuleClient(sysdigClients)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	sysdigClients.AddCleanupHook(sendPoliciesToAgents)

	return nil
}

func resourceSysdigRuleFalcoFromResourceData(d *schema.ResourceData) (v2.Rule, error) {
	rule := ruleFromResourceData(d)
	rule.Details.RuleType = v2.RuleTypeFalco

	appendMode, appendModeIsSet := d.GetOk("append")
	if appendModeIsSet {
		ptr := appendMode.(bool)
		rule.Details.Append = &ptr
	}

	if source, ok := d.GetOk("source"); ok && source.(string) != "" {
		rule.Details.Source = source.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return v2.Rule{}, errors.New("source must be set when append = false")
	}

	if output, ok := d.GetOk("output"); ok && output.(string) != "" {
		rule.Details.Output = output.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return v2.Rule{}, errors.New("output must be set when append = false")
	}

	if priority, ok := d.GetOk("priority"); ok && priority.(string) != "" {
		rule.Details.Priority = priority.(string)
	} else if !appendModeIsSet || !(appendMode.(bool)) {
		return v2.Rule{}, errors.New("priority must be set when append = false")
	}

	minimumEngineVersionInterface, ok := d.GetOk("minimum_engine_version")
	if ok {
		minimumEngineVersion := minimumEngineVersionInterface.(int)
		rule.Details.MinimumEngineVersion = &minimumEngineVersion
	}

	rule.Details.Condition = &v2.Condition{
		Condition:  d.Get("condition").(string),
		Components: []interface{}{},
	}

	if exceptionsField, ok := d.GetOk("exceptions"); ok {
		falcoExceptions := []*v2.Exception{}
		for _, exception := range exceptionsField.([]interface{}) {
			exceptionMap := exception.(map[string]interface{})
			newFalcoException := &v2.Exception{
				Name: exceptionMap["name"].(string),
			}

			fields := cast.ToStringSlice(exceptionMap["fields"])
			if len(fields) >= 1 {
				newFalcoException.Fields = fields
			}

			comps := cast.ToStringSlice(exceptionMap["comps"])
			if len(comps) >= 1 {
				newFalcoException.Comps = comps
			}

			values := cast.ToString(exceptionMap["values"])
			err := json.Unmarshal([]byte(values), &newFalcoException.Values)
			if err != nil {
				return v2.Rule{}, err
			}

			falcoExceptions = append(falcoExceptions, newFalcoException)
		}
		rule.Details.Exceptions = falcoExceptions
	}

	return rule, nil
}

package sysdig

import (
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"reflect"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Creates a rule with the default schema that a Secure Rule should have,
// which fields will be overwritten by the schema in the parameter.
func createRuleSchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	ruleSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}

	for k, v := range original {
		ruleSchema[k] = v
	}

	return ruleSchema
}

// Retrieves the common rule fields for a rule from a resource data.
func ruleFromResourceData(d *schema.ResourceData) v2.Rule {
	rule := v2.Rule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Version:     d.Get("version").(int),
	}

	rule.Tags = getTagsFromResourceData(d)

	return rule
}

// Saves in the resource data the information from the common fields of the rule.
func updateResourceDataForRule(d *schema.ResourceData, rule v2.Rule) {
	currentTags := getTagsFromResourceData(d)
	newTags := append([]string{}, rule.Tags...)
	sort.Strings(currentTags)
	sort.Strings(newTags)
	areTagsSame := reflect.DeepEqual(currentTags, newTags)

	_ = d.Set("name", rule.Name)
	_ = d.Set("description", rule.Description)
	if !areTagsSame {
		_ = d.Set("tags", rule.Tags)
	}
	_ = d.Set("version", rule.Version)

}

func getTagsFromResourceData(d *schema.ResourceData) []string {
	result := []string{}
	if tags, ok := d.Get("tags").([]interface{}); ok {
		for _, rawTag := range tags {
			if tag, ok := rawTag.(string); ok {
				result = append(result, tag)
			}
		}
	}

	return result
}

func getSecureRuleClient(c SysdigClients) (v2.RuleInterface, error) {
	return c.sysdigSecureClientV2()
}

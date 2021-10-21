package sysdig

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/internal/client/secure"
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
func ruleFromResourceData(d *schema.ResourceData) secure.Rule {
	rule := secure.Rule{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Version:     d.Get("version").(int),
	}

	rule.Tags = []string{}
	if tags, ok := d.Get("tags").([]interface{}); ok {
		for _, rawTag := range tags {
			if tag, ok := rawTag.(string); ok {
				rule.Tags = append(rule.Tags, tag)
			}
		}
	}
	return rule
}

// Saves in the resource data the information from the common fields of the rule.
func updateResourceDataForRule(d *schema.ResourceData, rule secure.Rule) {
	err := d.Set("name", rule.Name)
	if err != nil {
		log.Println("error assigning 'name'")
	}

	err = d.Set("description", rule.Description)
	if err != nil {
		log.Println("error assigning 'description'")
	}

	err = d.Set("tags", rule.Tags)
	if err != nil {
		log.Println("error assigning 'tags'")
	}

	err = d.Set("version", rule.Version)
	if err != nil {
		log.Println("error assigning 'version'")
	}

}

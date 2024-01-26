package sysdig

import (
	// "github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ReadOnlyIntSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Required: false, // Read-only value
		Optional: false, // Read-only value
		Computed: true,
	}
}

func ReadOnlyStringSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: false, // Read-only value
		Optional: false, // Read-only value
		Computed: true,
	}
}

func BoolSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func BoolComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
}

// Can be omitted for composite policies
func RuleNamesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// Can be omitted for Composite policies
func RulesSchema() *schema.Schema {
	return &schema.Schema{
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
	}
}

func ScopeSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	}
}

func ScopeComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func RunbookSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
}

func RunbookComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func VersionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
}

func NameSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func DescriptionSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func DescriptionComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func NotificationChannelsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

func NotificationChannelsComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

func EnabledSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
}

func EnabledComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
}

func SeveritySchema() *schema.Schema {
	return &schema.Schema{
		Type:             schema.TypeInt,
		Default:          4,
		Optional:         true,
		ValidateDiagFunc: validateDiagFunc(validation.IntBetween(0, 7)),
	}
}

func SeverityComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
}

func PreventMalwareActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
}

func PreventMalwareActionComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
}

func ContainerActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"stop", "pause", "kill"}, false),
	}
}

func ContainerActionComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
}

func CaptureActionSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"seconds_after_event": {
					Type:             schema.TypeInt,
					Required:         true,
					ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(0)),
				},
				"seconds_before_event": {
					Type:             schema.TypeInt,
					Required:         true,
					ValidateDiagFunc: validateDiagFunc(validation.IntAtLeast(0)),
				},
				"name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"filter": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"folder": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "/",
				},
			},
		},
	}
}

func CaptureActionComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"seconds_after_event": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"seconds_before_event": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"filter": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"bucket_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"folder": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func HashesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hash": {
					Type:     schema.TypeString,
					Required: true,
				},
				"hash_aliases": {
					Type:     schema.TypeSet,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func HashesComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hash": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"hash_aliases": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func TagsComputedSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Computed: true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

// Creates the common policy schema that is shared between policy resources
func createPolicySchema(original map[string]*schema.Schema) map[string]*schema.Schema {
	policySchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"notification_channels": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeInt,
			},
		},
		"runbook": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"actions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"container": ContainerActionSchema(),
					"capture":   CaptureActionSchema(),
				},
			},
		},
	}

	for k, v := range original {
		policySchema[k] = v
	}

	return policySchema
}

package sysdig

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	TeamSchemaThemeKey          = "theme"
	TeamSchemaNameKey           = "name"
	TeamSchemaDescriptionKey    = "description"
	TeamSchemaScopeByKey        = "scope_by"
	TeamSchemaFilterKey         = "filter"
	TeamSchemaUserRolesKey      = "user_roles"
	TeamSchemaUserRolesEmailKey = "email"
	TeamSchemaUserRolesRoleKey  = "role"
	TeamSchemaDefaultTeamKey    = "default_team"
	TeamSchemaVersionKey        = "version"
)

func createBaseTeamSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		TeamSchemaThemeKey: {
			Type: schema.TypeString,
		},
		TeamSchemaNameKey: {
			Type: schema.TypeString,
		},
		TeamSchemaDescriptionKey: {
			Type: schema.TypeString,
		},
		TeamSchemaScopeByKey: {
			Type: schema.TypeString,
		},
		TeamSchemaFilterKey: {
			Type: schema.TypeString,
		},
		TeamSchemaUserRolesKey: {
			Type: schema.TypeSet,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					TeamSchemaUserRolesEmailKey: {
						Type: schema.TypeString,
					},
					TeamSchemaUserRolesRoleKey: {
						Type: schema.TypeString,
					},
				},
			},
		},
		TeamSchemaDefaultTeamKey: {
			Type: schema.TypeBool,
		},
		TeamSchemaVersionKey: {
			Type: schema.TypeInt,
		},
	}
}

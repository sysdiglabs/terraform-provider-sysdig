package sysdig

import (
	"testing"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestSSOGroupMappingToResourceData_SortsTeamIDs(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSysdigSSOGroupMapping().Schema, map[string]any{})

	gm := &v2.SSOGroupMapping{
		GroupName:        "test-group",
		StandardTeamRole: "ROLE_TEAM_STANDARD",
		TeamMap: &v2.SSOGroupMappingTeamMap{
			IsForAllTeams: false,
			TeamIDs:       []int{3, 1, 2},
		},
	}

	err := ssoGroupMappingToResourceData(gm, d)
	require.NoError(t, err)

	teamMaps := d.Get("team_map").([]any)
	require.Len(t, teamMaps, 1)

	teamMap := teamMaps[0].(map[string]any)
	teamIDs := teamMap["team_ids"].([]any)
	require.Equal(t, []any{1, 2, 3}, teamIDs)
}

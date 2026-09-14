package sysdig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// The standalone component resource and the account resource's nested component block share
// accountComponent.Schema, so ForceNew may only be set on the standalone resource's own copy:
// on the nested block it would replace the entire cloud account when a component is renamed.
func TestCloudauthAccountComponentIdentityIsForceNew(t *testing.T) {
	standalone := resourceSysdigSecureCloudauthAccountComponent().Schema

	for _, key := range []string{SchemaAccountID, SchemaType, SchemaInstance} {
		if !standalone[key].ForceNew {
			t.Errorf("%s must be ForceNew: it identifies the component in the API path and in the resource id", key)
		}
	}

	if standalone[SchemaVersion].ForceNew {
		t.Errorf("%s must stay updatable in place", SchemaVersion)
	}

	nested := resourceSysdigSecureCloudauthAccount().Schema[SchemaComponent].Elem.(*schema.Resource).Schema
	for key, attribute := range nested {
		if attribute.ForceNew {
			t.Errorf("nested %s block: %s must not be ForceNew, it would replace the cloud account", SchemaComponent, key)
		}
	}
}

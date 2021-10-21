package sysdig

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Temporary wrapper for validate functions.
//
// Deprecated: use your own functions, this wrapper will be removed as
// soon as the new validate functions are supported by the SDK
func validateDiagFunc(validateFunc func(interface{}, string) ([]string, []error)) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		warnings, errs := validateFunc(i, fmt.Sprintf("%+v", path))
		var diags diag.Diagnostics
		for _, warning := range warnings {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  warning,
			})
		}
		for _, err := range errs {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error(),
			})
		}
		return diags
	}
}

// parseAzureCreds splits an Azure Trusted Identity into a tenantID and a clientID
func parseAzureCreds(azureTrustedIdentity string) (tenantID string, clientID string, err error) {
	tokens := strings.Split(azureTrustedIdentity, ":")
	if len(tokens) != 2 {
		return "", "", errors.New("Not a valid Azure Trusted Identity")
	}
	return tokens[0], tokens[1], nil
}

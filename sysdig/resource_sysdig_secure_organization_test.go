//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccSecureOrganization(t *testing.T) {
	// needs an actual existing gcp projectID
	accID := "mycapitalprojet"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: secureOrgWithAccountID(accID),
			},
			{
				ResourceName:      "sysdig_secure_organization.sample-org",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureOrgWithAccountID(accountID string) string {
	// this is a base64 encoded service account key, that should exist apriori in the project (mycapitalprojet)
	test_service_account_key_encoded := "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAibXljYXBpdGFscHJvamV0IiwKICAicHJpdmF0ZV9rZXlfaWQiOiAiNzRhMDZkMmE3YWNkN2Q3MmMyMDNhMzllYTZkYmE1Yzc0YTI3MWU5MiIsCiAgInByaXZhdGVfa2V5IjogIi0tLS0tQkVHSU4gUFJJVkFURSBLRVktLS0tLVxuTUlJRXZRSUJBREFOQmdrcWhraUc5dzBCQVFFRkFBU0NCS2N3Z2dTakFnRUFBb0lCQVFDK202ZTZlaVp3VWE5M1xuM1p4K3RuNEpqNmFCM1FTakpjQ3NSdkhFeThMTE1Za2FTbjdYSERBTEJXR3F0UXhCdGFIUnRWdXRzQWlMY29lNlxublJKZWNDTEl0K1hnWjFmcnFMeDJGTVZENzFUbW1JUUU3T3ZHNlBkR2pUVlFldTZjbFR2UjU0SHlZcDE0RlVxL1xuVjFlTUdHS29IK2JzZG9YL2I3T1Q5Z2RCTk80bVNRU29mZXQvQVlYVEVvSkJrVUlKd1YrV21CQkNUeE92UlBXbFxuOUMyNHN4alJocThVYWFkVlVQMWl6V3lURDJYb0oxejF5M0VUbnpJZGFVUExFQ2RoT0JHT0VWeEwxTXBmbWNadVxuV0tmdVM1cFJPR3R2VHFXVjZvd25IY3dVZjIwK0dYQnpSTFgzWG5qT2NPcmVjaTlhclVaeXNpV01ONW1tdDAvNFxuTkJWRmI1akZBZ01CQUFFQ2dnRUFSRGd4cit4SUZhc213aC92QXVzTjhSNFkyaUhncHdPOEVlelNXejBTV3VjL1xueXNrZTNKNmFFMWU1dlA5UGc3VWVWWkF6WDliQk9DZWxySVRMTGtHME1XS2dROUM1QnY4OWRJVzZsTFgwRFJSSlxuSy9vZjZQRTRqMmU2elRNeWM1aDE4SXFMVjlVenh1YlgrZU9vMGR0b2RBUDNBbXJwU3FNUUFzVHJrOHI5OFhWU1xuVFV2cERLekgyMjR0V01DL0VwRTNtaUtYNUpmRkRCaDBZTTlxV0lGbm80Q0FRSGhYRmFiTU5iMW9SN0M4cmtDWlxuRm45WXVjQjJYRkdjNm5pQ0pZRXFMYUJKVW9uZDFlS0dCTi9OV2JGdDJVendYS2dLZldDZnJYMitjalBNYTkvQ1xuWUcxVnRBdzJRUGJweFcxeHBZOVRGWTZxYU93L2lpUUNVa2xtUExBZmdRS0JnUUR0WHJ4OWtMcy9objREM2w5Y1xuclNoZlVlNzByMXFmT0FMMk9EWldrdUxVd2VucXNqb1ZLeStvbUo3UVJIZFllREViMitZdlM5cllrWWV2ck5UelxuQ0pEcW1qNUpDaEJHbzFIOG9HdTlLK3UyQmo3WHpWc0pqZFlUQ0phOGJnNjhERmtpbHlvbm0rZ3g0OStDb1Y1dlxuQ0hYZWRNckVtNlo0QWpzV2ZZTndobXBNbHdLQmdRRE5rVjhMMzE0bUZCeEJpdHNGb1pKQkd0QmdURHhMajd1SVxuRHcyeVk5b0xGR204VnhjM05Qcmt0K3FIUytCTXFaNHB6SEhXb3RJazJFZWlJWDFzbU9EMDJTclY3c2s1c245ZFxuRzJVMWtkZWFJb1RTWkpFTENLZzg1TzMxb2JOcDgzT3B0ZDB5SStrUlZsZXA3LzBsdnRLTnFGR05DWGN4MEJwTVxuY2RTeXBsUkZBd0tCZ0YvNFdwc2w1aDhFQUhVTjlsNWhBQjZ4NEx2N0hkZWI3TTZoNFk1Vkt2SzhTQmdFNFNqblxuNGdmM1ZOWjlxQWNUNlQ4TFJHREErWVZ3S1h6a2t1Q0VDUnRoSzJlYWN3UXNTaHlxdTRTcmVreUk3K1dPZUkwL1xuVkZzenNNWVVkVTZnYTNWcHlyaGk5NWtjT2FUMkcxa25BWWprallxNko2OERyK0lpOHY2T3lmR1hBb0dCQU1IeFxuZ0RITVdLQW1ZdzQzT2pLUzRGQ0tRc1JIeUs4bGVUR0J1bE51djMycWthTitxMG1MczVYc0t0bmc3VXFHME5Ed1xuc1FwbWJVc1R2bW1wblJMREhhSUQ3ZFVPeDB5bkttQ21neE5LZUpaVU1PbnF0YWtxVHNlODJRRGd3VXVadzZyL1xuQ1NUUUdva2Y0KzlSbTQxci9teGx2Q01MSmlpYUJPWFFrM0xGV0VZUEFvR0FVMldIVjBrNVBNdnNiKzhwZm4zaVxudEdFWjNJNlZvWGthdHc4VzJJL0YvdzJzOENjeEpoa0VvYzFNaysrVTVFRnk5YlBlb3JQQjRrM1ViNVZIbms0blxuN3JjZGl1a3ZLUjVNajdldjl3RnlNdGcwWEdVeDZxbC8waE1oLzZ1SmlBcnIvcVE0ekpEY0E0SW5pUW9QdmI0ZFxuUGRxY2p6Q3Bhc3JHTnJlQ2pnczBnWW89XG4tLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tXG4iLAogICJjbGllbnRfZW1haWwiOiAic3lzZGlnLXNlY3VyZS1wcm92aWRlci10ZXN0QG15Y2FwaXRhbHByb2pldC5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgImNsaWVudF9pZCI6ICIxMDYxODY3MTY5ODc0MDI4MDA2NjgiLAogICJhdXRoX3VyaSI6ICJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9vYXV0aDIvYXV0aCIsCiAgInRva2VuX3VyaSI6ICJodHRwczovL29hdXRoMi5nb29nbGVhcGlzLmNvbS90b2tlbiIsCiAgImF1dGhfcHJvdmlkZXJfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLAogICJjbGllbnRfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9yb2JvdC92MS9tZXRhZGF0YS94NTA5L3N5c2RpZy1zZWN1cmUtcHJvdmlkZXItdGVzdCU0MG15Y2FwaXRhbHByb2pldC5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgInVuaXZlcnNlX2RvbWFpbiI6ICJnb29nbGVhcGlzLmNvbSIKfQo="

	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "%s"
  provider_type = "PROVIDER_GCP"
  enabled       = "true"
  feature {
	secure_config_posture {
	  enabled    = "true"
	  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
	}
	secure_identity_entitlement {
	  enabled    = true
	  components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
	}
  }
  component {
	type                       = "COMPONENT_SERVICE_PRINCIPAL"
	instance                   = "secure-posture"
	service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
  component {
    type     				   = "COMPONENT_SERVICE_PRINCIPAL"
    instance 				   = "secure-onboarding"
    service_principal_metadata = jsonencode({
      gcp = {
        key = "%s"
      }
    })
  }
  lifecycle {
	ignore_changes = [component]
  }
}

resource "sysdig_secure_organization" "sample-org" {
  management_account_id		= sysdig_secure_cloud_auth_account.sample.id
}
`, accountID, test_service_account_key_encoded, test_service_account_key_encoded)
}

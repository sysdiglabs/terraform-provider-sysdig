//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

func TestAccGCPSecureCloudAuthAccount(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: secureGCPCloudAuthAccountMinimumConfiguration(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureGCPCloudAuthAccountMinimumConfiguration(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
  provider_id   = "sample-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = true
}`, accountID)
}

func TestAccGCPSecureCloudAuthAccountFC(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: secureGCPCloudAuthAccountWithFC(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample-1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureGCPCloudAuthAccountWithFC(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample-1" {
  provider_id   = "sample-1-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = true
  feature {
	secure_config_posture {
	  enabled    = true
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
}
`, accountID, getEncodedServiceAccountKey("sample-1", accountID))
}

func getEncodedServiceAccountKey(resourceName string, accountID string) string {
	type sample_service_account_key struct {
		Type                    string `json:"type"`
		ProjectId               string `json:"project_id"`
		PrivateKeyId            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientId                string `json:"client_id"`
		AuthUri                 string `json:"auth_uri"`
		TokenUri                string `json:"token_uri"`
		AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
		ClientX509CertUrl       string `json:"client_x509_cert_url"`
		UniverseDomain          string `json:"universe_domain"`
	}
	test_service_account_key := &sample_service_account_key{
		Type:                    "service_account",
		ProjectId:               fmt.Sprintf("%s-%s", resourceName, accountID),
		PrivateKeyId:            "xxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		PrivateKey:              "-----BEGIN PRIVATE KEY-----\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n-----END PRIVATE KEY-----\n",
		ClientEmail:             fmt.Sprintf("some-sa-name@%s-%s.iam.gserviceaccount.com", resourceName, accountID),
		ClientId:                "some-client-id",
		AuthUri:                 "https://some-auth-uri",
		TokenUri:                "https://some-token-uri",
		AuthProviderX509CertUrl: "https://some-authprovider-cert-url",
		ClientX509CertUrl:       "https://some-client-cert-url",
		UniverseDomain:          "googleapis.com",
	}

	test_service_account_key_bytes, err := json.Marshal(test_service_account_key)
	if err != nil {
		fmt.Printf("Failed to marshal test_service_account_key: %v", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, test_service_account_key_bytes, "", "  ")
	if err != nil {
		fmt.Printf("Failed to indent test_service_account_key: %v", err)
	}
	out.WriteByte('\n')

	test_service_account_key_encoded := b64.StdEncoding.EncodeToString(out.Bytes())
	return test_service_account_key_encoded
}

func TestAccAzureSecureCloudAccount(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: secureCloudAuthAccountMinimumConfigurationAzure(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureCloudAuthAccountMinimumConfigurationAzure(accountId string) string {
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
	  provider_id   = "sample-%s"
	  provider_type = "PROVIDER_AZURE"
	  enabled       = true
	  provider_tenant_id = "%s"
	  provider_alias = "some-alias"
	}`, accountId, randomTenantId)
}

func TestAccAzureSecureCloudAccountFC(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: secureAzureCloudAuthAccountWithFC(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample-1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureCloudAuthAccountWithFC(accountID string) string {
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "sample-1" {
			provider_id   = "sample-1-%s"
			provider_type = "PROVIDER_AZURE"
			enabled       = true
			provider_tenant_id = "%s"
			provider_alias = "some-alias"
			feature {
				secure_config_posture {
					enabled    = true
					components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
				}
			}
			component {
				type                       = "COMPONENT_SERVICE_PRINCIPAL"
				instance                   = "secure-posture"
				service_principal_metadata = jsonencode({
					azure = {
						active_directory_service_principal = {
							id                        = "some-id"
							account_enabled           = true
							display_name              = "some-display-name"
							app_display_name          = "some-app-display-name"
							app_id                    = "some-app-id"
							app_owner_organization_id = "some-app-owner-organization-id"
						}
					}
				})
			}
		}`, accountID, randomTenantId)
}

func TestAccAzureSecureCloudAccountFCThreatDetection(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: secureAzureCloudAuthAccountWithFCThreatDetection(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample-1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureCloudAuthAccountWithFCThreatDetection(accountID string) string {
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "sample-1" {
			provider_id   = "sample-1-%s"
			provider_type = "PROVIDER_AZURE"
			enabled       = true
			provider_tenant_id = "%s"
			feature {
				secure_threat_detection {
					enabled    = true
					components = ["COMPONENT_EVENT_BRIDGE/secure-runtime"]
				  }
			}
			component {
				type                       = "COMPONENT_EVENT_BRIDGE"
				instance                   = "secure-runtime"
				event_bridge_metadata = jsonencode({
					azure = {
						event_hub_metadata= {
							event_hub_name      = "event-hub-name"
							event_hub_namespace = "event-hub-namespace"
							consumer_group      = "consumer-group"
						}
					}
				})
			}
		}`, accountID, randomTenantId)
}

func TestGCPAgentlesScanningOnboarding(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accID := rText()
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
				Config: getResourceForGCPAgentlessScan(accID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "provider_type", "PROVIDER_GCP"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "feature.0.secure_agentless_scanning.0.enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "feature.0.secure_agentless_scanning.0.components.0", "COMPONENT_SERVICE_PRINCIPAL/secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "component.0.type", "COMPONENT_SERVICE_PRINCIPAL"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "component.0.instance", "secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp-agentless-scanning", "component.0.service_principal_metadata", "{\"gcp\":{\"email\":\"email_value\",\"workload_identity_federation\":{\"pool_provider_id\":\"pool_provider_id_value\"}}}"),
				),
			},
			{
				ResourceName:            "sysdig_secure_cloud_auth_account.gcp-agentless-scanning",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"component"},
			},
		},
	})
}

func getResourceForGCPAgentlessScan(projectID string) string {
	return fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "gcp-agentless-scanning" {
			provider_id   = "gcp-agentless-test-%s"
			provider_type = "PROVIDER_GCP"
			enabled       = true

		    feature {
			  secure_agentless_scanning {
			    enabled    = true
			    components = ["COMPONENT_SERVICE_PRINCIPAL/secure-scanning"]
			  }
		    }

			component {
				type                       = "COMPONENT_SERVICE_PRINCIPAL"
				instance                   = "secure-scanning"
				service_principal_metadata = jsonencode({
					gcp = {
						workload_identity_federation = {
							pool_provider_id = "pool_provider_id_value"
						}
						email = "email_value"
					}
				})
			}
		}`, projectID)
}

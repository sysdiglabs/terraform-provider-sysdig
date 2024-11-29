//go:build tf_acc_sysdig_secure || tf_acc_sysdig_common

package sysdig_test

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

/************
* GCP tests
************/
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

func TestAccGCPSecureCloudAuthAccountConfigPosture(t *testing.T) {
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
				Config: secureGCPCloudAuthAccountWithConfigPosture(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.gcp_config_posture",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureGCPCloudAuthAccountWithConfigPosture(accountID string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "gcp_config_posture" {
  provider_id   = "gcp-cspm-test-%s"
  provider_type = "PROVIDER_GCP"
  enabled       = true
  feature {
    secure_config_posture {
      enabled    = true
      components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
    }
    secure_identity_entitlement {
      enabled    = true
      components = ["COMPONENT_WEBHOOK_DATASOURCE/secure-runtime"]
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
		type                        = "COMPONENT_WEBHOOK_DATASOURCE"
		instance                    = "secure-runtime"
		webhook_datasource_metadata = jsonencode({
			gcp = {
				webhook_datasource = {
					pubsub_topic_name      = "pubsub_topic_name_value"
					sink_name              = "sink_name_value"
					push_subscription_name = "push_subscription_name_value"
					push_endpoint          = "push_endpoint_value"
				}
			  service_principal = {
					workload_identity_federation = {
						pool_id          = "pool_id_value"
						pool_provider_id = "pool_provider_id_value"
						project_number   = "123456789011"
					}
					email = "email_value"
				}
			}
		})
	}
}
`, accountID, getEncodedServiceAccountKey("gcp-cspm-test", accountID))
}

func TestAccGCPSecureCloudAuthAccountAgentlesScanning(t *testing.T) {
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
				Config: fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "gcp_agentless_scanning" {
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
		}`, accID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "provider_type", "PROVIDER_GCP"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "feature.0.secure_agentless_scanning.0.enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "feature.0.secure_agentless_scanning.0.components.0", "COMPONENT_SERVICE_PRINCIPAL/secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.type", "COMPONENT_SERVICE_PRINCIPAL"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.instance", "secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.service_principal_metadata", "{\"gcp\":{\"email\":\"email_value\",\"workload_identity_federation\":{\"pool_provider_id\":\"pool_provider_id_value\"}}}"),
				),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.gcp_agentless_scanning",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGCPSecureCloudAuthAccountAgentlesScanningWithInventory(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accountID := rText()
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
				Config: fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "gcp_agentless_scanning" {
			provider_id   = "gcp-agentless-test-%s"
			provider_type = "PROVIDER_GCP"
			enabled       = true

		    feature {
              secure_config_posture {
              enabled    = true
                components = ["COMPONENT_SERVICE_PRINCIPAL/secure-posture"]
              }

			  secure_agentless_scanning {
			    enabled    = true
			    components = ["COMPONENT_SERVICE_PRINCIPAL/secure-scanning"]
			  }
		    }

            component {
				type     = "COMPONENT_SERVICE_PRINCIPAL"
				instance = "secure-posture"
				service_principal_metadata = jsonencode({
				  gcp = {
					key = "%s"
				  }
				})
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
		}`, accountID, getEncodedServiceAccountKey("gcp-agentless-test", accountID)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "provider_type", "PROVIDER_GCP"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "feature.0.secure_agentless_scanning.0.enabled", "true"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "feature.0.secure_agentless_scanning.0.components.0", "COMPONENT_SERVICE_PRINCIPAL/secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.type", "COMPONENT_SERVICE_PRINCIPAL"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.instance", "secure-posture"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.0.service_principal_metadata", "{\"gcp\":{\"key\":\""+getEncodedServiceAccountKey("gcp-agentless-test", accountID)+"\"}}"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.1.instance", "secure-scanning"),
					resource.TestCheckResourceAttr("sysdig_secure_cloud_auth_account.gcp_agentless_scanning", "component.1.service_principal_metadata", "{\"gcp\":{\"email\":\"email_value\",\"workload_identity_federation\":{\"pool_provider_id\":\"pool_provider_id_value\"}}}"),
				),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.gcp_agentless_scanning",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGCPSecureCloudAuthAccountThreatDetection(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }
	accountID := rText()
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
				Config: fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "gcp_threat_detection" {
			provider_id   = "gcp-cdr-test-%s"
			provider_type = "PROVIDER_GCP"
			enabled       = true
		    feature {
			  secure_threat_detection {
			    enabled    = true
			    components = ["COMPONENT_WEBHOOK_DATASOURCE/secure-runtime", "COMPONENT_SERVICE_PRINCIPAL/secure-runtime"]
			  }
		    }
			component {
				type                        = "COMPONENT_WEBHOOK_DATASOURCE"
				instance                    = "secure-runtime"
				webhook_datasource_metadata = jsonencode({
					gcp = {
						webhook_datasource = {
							pubsub_topic_name      = "pubsub_topic_name_value"
							sink_name              = "sink_name_value"
							push_subscription_name = "push_subscription_name_value"
							push_endpoint          = "push_endpoint_value"
						}
					}
				})
			}
			component {
				type                       = "COMPONENT_SERVICE_PRINCIPAL"
				instance                   = "secure-runtime"
				service_principal_metadata = jsonencode({
					gcp = {
						workload_identity_federation = {
							pool_id          = "pool_id_value"
							pool_provider_id = "pool_provider_id_value"
							project_number   = "123456789011"
						}
						email = "email_value"
					}
				})
			}
		}`, accountID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.gcp_threat_detection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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

/*************
* Azure tests
*************/
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
				Config: secureAzureCloudAuthAccountMinimumConfiguration(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureCloudAuthAccountMinimumConfiguration(accountId string) string {
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

func TestAccAzureSecureCloudAccountConfigPosture(t *testing.T) {
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
				Config: secureAzureCloudAuthAccountWithConfigPosture(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.azure_config_posture",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureCloudAuthAccountWithConfigPosture(accountID string) string {
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "azure_config_posture" {
			provider_id   = "azure-cspm-test-%s"
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

func TestAccAzureSecureCloudAccountThreatDetection(t *testing.T) {
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
				Config: secureAzureCloudAuthAccountWithThreatDetection(accID),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.azure_threat_detection",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureAzureCloudAuthAccountWithThreatDetection(accountID string) string {
	rID := func() string { return acctest.RandStringFromCharSet(36, acctest.CharSetAlphaNum) }
	randomTenantId := rID()

	return fmt.Sprintf(`
		resource "sysdig_secure_cloud_auth_account" "azure_threat_detection" {
			provider_id   = "azure-cdr-test-%s"
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

/************
* AWS tests
************/
func TestAccAWSSecureCloudAccountThreatDetection(t *testing.T) {
	accountID := fmt.Sprintf("%012d", rand.Intn(99999999999))
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
				Config: fmt.Sprintf(`
				resource "sysdig_secure_cloud_auth_account" "aws_account_%s" {
					enabled       = true
					provider_id   = "%s"
					provider_type = "PROVIDER_AWS"

					feature {

						secure_threat_detection {
							enabled    = true
							components = ["COMPONENT_EVENT_BRIDGE/secure-runtime"]
						}
					}
					component {
						type     = "COMPONENT_EVENT_BRIDGE"
						instance = "secure-runtime"
						event_bridge_metadata = jsonencode({
							aws = {
								role_name = "sysdig-secure-events-ezsz"
								rule_name = "sysdig-secure-events-ezsz"
							}
						})
					}
				}`, accountID, accountID),
			},
			{
				ResourceName:      fmt.Sprintf("sysdig_secure_cloud_auth_account.aws_account_%s", accountID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAWSSecureCloudAccountConfigPostureAndAgentlessScanning(t *testing.T) {
	accountID := fmt.Sprintf("%012d", rand.Intn(99999999999))
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
				Config: fmt.Sprintf(`
				resource "sysdig_secure_cloud_auth_account" "aws_account_%s" {
					enabled       = true
					provider_id   = "%s"
					provider_type = "PROVIDER_AWS"

					feature {

						secure_config_posture {
							enabled    = true
							components = ["COMPONENT_TRUSTED_ROLE/secure-posture"]
						}

						secure_agentless_scanning {
							enabled    = true
							components = ["COMPONENT_TRUSTED_ROLE/secure-scanning", "COMPONENT_CRYPTO_KEY/secure-scanning"]
						}
					}
					component {
						type     = "COMPONENT_TRUSTED_ROLE"
						instance = "secure-scanning"
						trusted_role_metadata = jsonencode({
							aws = {
								role_name = "sysdig-secure-scanning-ob1o"
							}
						})
					}
					component {
						type     = "COMPONENT_CRYPTO_KEY"
						instance = "secure-scanning"
						crypto_key_metadata = jsonencode({
							aws = {
								kms = {
									alias    = "alias/sysdig-secure-scanning-ob1o"
									regions  = [
										"us-east-1",
										"us-west-2",
									]
								}
							}
						})
					}
					component {
						type     = "COMPONENT_TRUSTED_ROLE"
						instance = "secure-posture"
						trusted_role_metadata = jsonencode({
							aws = {
								role_name = "sysdig-secure-bu1k"
							}
						})
					}
				}`, accountID, accountID),
			},
			{
				ResourceName:      fmt.Sprintf("sysdig_secure_cloud_auth_account.aws_account_%s", accountID),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

/*************
* Oracle tests
*************/
func TestAccOracleSecureCloudAccountRoot(t *testing.T) {
	rID := func() string { return acctest.RandStringFromCharSet(60, acctest.CharSetAlphaNum) }
	tenantID := rID()
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
				Config: secureOracleCloudAuthAccountMinimumConfiguration(tenantID, false),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOracleSecureCloudAccountCompartment(t *testing.T) {
	rID := func() string { return acctest.RandStringFromCharSet(60, acctest.CharSetAlphaNum) }
	tenantID := rID()
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
				Config: secureOracleCloudAuthAccountMinimumConfiguration(tenantID, true),
			},
			{
				ResourceName:      "sysdig_secure_cloud_auth_account.sample",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func secureOracleCloudAuthAccountMinimumConfiguration(tenantID string, compartmentOnboard bool) string {
	rID := func() string { return acctest.RandStringFromCharSet(60, acctest.CharSetAlphaNum) }
	compartmentID := rID()
	if !compartmentOnboard {
		compartmentID = tenantID
	}
	return fmt.Sprintf(`
resource "sysdig_secure_cloud_auth_account" "sample" {
	  provider_id   = "%s"
	  provider_type = "PROVIDER_ORACLECLOUD"
	  enabled       = true
	  provider_tenant_id = "%s"

# TODO: add component back when https://github.com/draios/secure-backend/pull/38958 is merged
#	  component {
#		type     = "COMPONENT_SERVICE_PRINCIPAL"
#		instance = "secure-onboarding"
#		version  = "v0.1.0"
#		service_principal_metadata = jsonencode({
#		  oci = {
#			api_key = {
#			  user_id = "user-id"
#			}
#		  }
#		})
#	  }
	}`, fmt.Sprintf("ocid1.tenancy.oc1..%s", compartmentID), fmt.Sprintf("ocid1.tenancy.oc1..%s", tenantID))

}

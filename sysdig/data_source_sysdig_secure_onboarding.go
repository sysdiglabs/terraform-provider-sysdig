package sysdig

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
)

func getSecureOnboardingClient(c SysdigClients) (v2.OnboardingSecureInterface, error) {
	return c.sysdigSecureClientV2()
}

func dataSourceSysdigSecureTrustedCloudIdentity() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTrustedCloudIdentityRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"aws", "gcp", "azure"}, false),
			},
			"identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"azure_service_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gov_identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_gov_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"aws_gov_role_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTrustedCloudIdentityRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	// get trusted identity for commercial backend
	identity, err := client.GetTrustedCloudIdentitySecure(ctx, d.Get("cloud_provider").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// get trusted identity for regulatory backend, such as govcloud
	// XXX: only supported for aws currently. update when supported for other providers
	var trustedRegulation map[string]string
	if d.Get("cloud_provider").(string) == "aws" {
		trustedRegulation, err = client.GetTrustedCloudRegulationAssetsSecure(ctx, d.Get("cloud_provider").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(identity)

	provider := d.Get("cloud_provider")
	switch provider {
	case "aws":
		// set the commercial identity
		_ = d.Set("identity", identity)
		// if identity is an ARN, attempt to extract certain fields
		parsedArn, err := arn.Parse(identity)
		if err == nil {
			_ = d.Set("aws_account_id", parsedArn.AccountID)
			if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
				_ = d.Set("aws_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
			}
		}

		// set the gov regulation based identity (only supported for aws currently)
		err = d.Set("gov_identity", trustedRegulation["trustedIdentityGov"])
		if err != nil {
			return diag.FromErr(err)
		}
		// if identity is an ARN, attempt to extract certain fields
		parsedArn, err = arn.Parse(trustedRegulation["trustedIdentityGov"])
		if err == nil {
			_ = d.Set("aws_gov_account_id", parsedArn.AccountID)
			if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
				_ = d.Set("aws_gov_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
			}
		}
	case "gcp":
		// set the commercial identity
		_ = d.Set("identity", identity)
		// if identity is an ARN, attempt to extract certain fields
		parsedArn, err := arn.Parse(identity)
		if err == nil {
			_ = d.Set("aws_account_id", parsedArn.AccountID)
			if parsedArn.Service == "iam" && strings.HasPrefix(parsedArn.Resource, "role/") {
				_ = d.Set("aws_role_name", strings.TrimPrefix(parsedArn.Resource, "role/"))
			}
		}
	case "azure":
		// set the commercial identity
		_ = d.Set("identity", identity)
		// if identity is an Azure tenantID/clientID, separate into each part
		tenantID, spID, err := parseAzureCreds(identity)
		if err == nil {
			_ = d.Set("azure_tenant_id", tenantID)
			_ = d.Set("azure_service_principal_id", spID)

		}
	}
	return nil
}

func dataSourceSysdigSecureTrustedAzureApp() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTrustedAzureAppRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"config_posture", "onboarding", "threat_detection", "vm_workload_scanning"}, false),
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTrustedAzureAppRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	app := d.Get("name").(string)
	registration, err := client.GetTrustedAzureAppSecure(ctx, app)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(app)
	for k, v := range registration {
		fmt.Printf("%s, %s\n", k, snakeCase(k))
		err = d.Set(snakeCase(k), v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func dataSourceSysdigSecureTenantExternalID() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTenantExternalIDRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"external_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureTenantExternalIDRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	externalID, err := client.GetTenantExternalIDSecure(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(externalID)
	err = d.Set("external_id", externalID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func dataSourceSysdigSecureAgentlessScanningAssets() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureAgentlessScanningAssetsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"aws": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"azure": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"backend": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gcp": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureAgentlessScanningAssetsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	assets, err := client.GetAgentlessScanningAssetsSecure(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	assetsAws, _ := assets["aws"].(map[string]any)
	assetsAzure, _ := assets["azure"].(map[string]any)
	assetsBackend, _ := assets["backend"].(map[string]any)
	assetsGcp, _ := assets["gcp"].(map[string]any)

	d.SetId("agentlessScanningAssets")
	err = d.Set("aws", map[string]any{
		"account_id": assetsAws["accountId"],
	})
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("azure", map[string]any{
		"service_principal_id": assetsAzure["servicePrincipalId"],
		"tenant_id":            assetsAzure["tenantId"],
	})
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("backend", map[string]any{
		"cloud_id": assetsBackend["cloudId"],
		"type":     assetsBackend["type"],
	})
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("gcp", map[string]any{
		"worker_identity": assetsGcp["workerIdentity"],
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func dataSourceSysdigSecureCloudIngestionAssets() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureCloudIngestionAssetsRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"aws", "gcp", "azure"}, false),
			},
			"cloud_provider_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"COMPONENT_WEBHOOK_DATASOURCE"}, false),
			},
			"aws": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"gcp_routing_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gcp_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource form the file and loads it in Terraform
func dataSourceSysdigSecureCloudIngestionAssetsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	cloudProvider := ""
	if v, ok := d.GetOk("cloud_provider"); ok {
		cloudProvider = v.(string)
	}

	cloudProviderID := ""
	if v, ok := d.GetOk("cloud_provider_id"); ok {
		cloudProviderID = v.(string)
	}

	componentType := ""
	if v, ok := d.GetOk("component_type"); ok {
		componentType = v.(string)
	}

	assets, err := client.GetCloudIngestionAssetsSecure(ctx, cloudProvider, cloudProviderID, componentType)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("cloudIngestionAssets")

	// Set GCP data if available
	if gcpAssets, ok := assets["gcp"].(map[string]any); ok {
		if routingKey, exists := gcpAssets["routingKey"]; exists {
			if err := d.Set("gcp_routing_key", routingKey); err != nil {
				return diag.FromErr(err)
			}
		}

		if metadata, exists := gcpAssets["metadata"]; exists {
			if err := d.Set("gcp_metadata", metadata); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Set AWS data if available
	if awsAssets, ok := assets["aws"].(map[string]any); ok {
		awsData := map[string]any{
			"eventBusARN":    awsAssets["eventBusARN"],
			"eventBusARNGov": awsAssets["eventBusARNGov"],
		}

		// Add SNS specific fields if available
		if awsAssets["snsRoutingKey"] != nil {
			awsData["sns_routing_key"] = awsAssets["snsRoutingKey"]
		}

		if snsMetadata, ok := awsAssets["snsMetadata"].(map[string]any); ok && snsMetadata != nil {
			if ingestionURL, exists := snsMetadata["ingestionURL"]; exists {
				awsData["sns_routing_url"] = ingestionURL
			}
		}

		// Add EventBridge specific fields if available
		if awsAssets["ebRoutingKey"] != nil {
			awsData["eb_routing_key"] = awsAssets["ebRoutingKey"]
		}

		if ebMetadata, ok := awsAssets["ebMetadata"].(map[string]any); ok && ebMetadata != nil {
			if ingestionURL, exists := ebMetadata["ingestionURL"]; exists {
				awsData["eb_routing_url"] = ingestionURL
			}
			if apiKey, exists := ebMetadata["apiKey"]; exists {
				awsData["eb_api_key"] = apiKey
			}
		}

		if err := d.Set("aws", awsData); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func dataSourceSysdigSecureTrustedOracleApp() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigSecureTrustedOracleAppRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"config_posture", "onboarding"}, false),
			},
			"tenancy_ocid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_ocid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_ocid": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// Retrieves the information of a resource from the file and loads it in Terraform
func dataSourceSysdigSecureTrustedOracleAppRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureOnboardingClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	app := d.Get("name").(string)
	trustedIdentityGroup, err := client.GetTrustedOracleAppSecure(ctx, app)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(app)
	for k, v := range trustedIdentityGroup {
		fmt.Printf("%s, %s\n", k, snakeCase(k))
		err = d.Set(snakeCase(k), v)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func snakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

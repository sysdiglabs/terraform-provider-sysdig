package sysdig

import (
	"context"
	"strconv"
	"strings"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigMonitorCloudAccount() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigMonitorCloudAccountCreate,
		DeleteContext: resourceSysdigMonitorCloudAccountDelete,
		ReadContext:   resourceSysdigMonitorCloudAccountRead,
		UpdateContext: resourceSysdigMonitorCloudAccountUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			"cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"integration_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"role_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"access_key_id": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"additional_options": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func getMonitorCloudAccountClient(c SysdigClients) (v2.CloudAccountMonitorInterface, error) {
	return c.sysdigMonitorClientV2()
}

func resourceSysdigMonitorCloudAccountCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("integration_type").(string) == "Cost" {
		cloudAccount := monitorCloudAccountForCostFromResourceData(data)
		cloudAccountCreated, err := client.CreateCloudAccountMonitorForCost(ctx, &cloudAccount)

		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(cloudAccountCreated.Id)
	} else {
		cloudAccount := monitorCloudAccountFromResourceData(data)
		cloudAccountCreated, err := client.CreateCloudAccountMonitor(ctx, &cloudAccount)

		if err != nil {
			return diag.FromErr(err)
		}

		data.SetId(strconv.Itoa(cloudAccountCreated.Id))
	}

	return nil
}

func resourceSysdigMonitorCloudAccountDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteCloudAccountMonitor(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigMonitorCloudAccountRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("integration_type").(string) == "Cost" {
		cloudAccount, err := client.GetCloudAccountMonitorForCost(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		err = monitorCloudAccountForCostToResourceData(data, cloudAccount)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		cloudAccount, err := client.GetCloudAccountMonitor(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		err = monitorCloudAccountToResourceData(data, cloudAccount)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceSysdigMonitorCloudAccountUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client, err := getMonitorCloudAccountClient(i.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if data.Get("integration_type").(string) == "Cost" {
		putObjectAccount := monitorCloudAccountForCostFromResourceDataPutMethod(data)

		cloudAccount, err := client.GetCloudAccountMonitorForCost(ctx, id)
		if err != nil {
			return diag.FromErr(err)
		}

		completeEmptyFieldsForUpdateRequest(&putObjectAccount, cloudAccount)

		_, err = client.UpdateCloudAccountMonitorForCost(ctx, &putObjectAccount)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		cloudAccount := monitorCloudAccountFromResourceData(data)

		id, err := strconv.Atoi(data.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateCloudAccountMonitor(ctx, id, &cloudAccount)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func monitorCloudAccountFromResourceData(data *schema.ResourceData) v2.CloudAccountMonitor {
	return v2.CloudAccountMonitor{
		Platform:          data.Get("cloud_provider").(string),
		IntegrationType:   data.Get("integration_type").(string),
		AdditionalOptions: data.Get("additional_options").(string),
		Credentials: v2.CloudAccountCredentialsMonitor{
			AccountId:   data.Get("account_id").(string),
			RoleName:    data.Get("role_name").(string),
			SecretKey:   data.Get("secret_key").(string),
			AccessKeyId: data.Get("access_key_id").(string),
		},
	}
}

func monitorCloudAccountForCostFromResourceData(data *schema.ResourceData) v2.CloudAccountMonitorForCost {

	configuration := v2.CloudCostConfiguration{}

	if config, ok := data.Get("config").(map[string]interface{}); ok {
		if val, exists := config["athena_bucket_name"]; exists {
			configuration.AthenaBucketName = val.(string)
		}
		if val, exists := config["athena_database_name"]; exists {
			configuration.AthenaDatabaseName = val.(string)
		}
		if val, exists := config["athena_region"]; exists {
			configuration.AthenaRegion = val.(string)
		}
		if val, exists := config["athena_workgroup"]; exists {
			configuration.AthenaWorkgroup = val.(string)
		}
		if val, exists := config["athena_table_name"]; exists {
			configuration.AthenaTableName = val.(string)
		}
		if val, exists := config["spot_prices_bucket_name"]; exists {
			configuration.SpotPricesBucketName = val.(string)
		}
	}

	return v2.CloudAccountMonitorForCost{
		Feature:       "Cloud cost",
		Platform:      data.Get("cloud_provider").(string),
		Configuration: configuration,
		Credentials: v2.CloudAccountCredentialsMonitor{
			AccountId: data.Get("account_id").(string),
			RoleName:  data.Get("role_name").(string),
		},
	}
}

func monitorCloudAccountForCostFromResourceDataPutMethod(data *schema.ResourceData) v2.CloudAccountCostProvider {
	return v2.CloudAccountCostProvider{
		Provider: data.Get("cloud_provider").(string),
		RoleArn:  "arn:aws:iam::" + data.Get("account_id").(string) + ":role/" + data.Get("role_name").(string),
		Config: v2.CloudConfigForCost{
			AthenaBucketName:     data.Get("config").(map[string]interface{})["athena_bucket_name"].(string),
			AthenaDatabaseName:   data.Get("config").(map[string]interface{})["athena_database_name"].(string),
			AthenaRegion:         data.Get("config").(map[string]interface{})["athena_region"].(string),
			AthenaWorkgroup:      data.Get("config").(map[string]interface{})["athena_workgroup"].(string),
			AthenaTableName:      data.Get("config").(map[string]interface{})["athena_table_name"].(string),
			SpotPricesBucketName: data.Get("config").(map[string]interface{})["spot_prices_bucket_name"].(string),
			IntegrationType:      data.Get("integration_type").(string),
		},
	}
}

func monitorCloudAccountToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudAccountMonitor) error {
	err := data.Set("cloud_provider", cloudAccount.Platform)
	if err != nil {
		return err
	}

	err = data.Set("integration_type", cloudAccount.IntegrationType)
	if err != nil {
		return err
	}

	err = data.Set("additional_options", cloudAccount.AdditionalOptions)
	if err != nil {
		return err
	}

	err = data.Set("account_id", cloudAccount.Credentials.AccountId)
	if err != nil {
		return err
	}

	err = data.Set("role_name", strings.Split(cloudAccount.Credentials.RoleName, ":role/")[1])
	if err != nil {
		return err
	}

	err = data.Set("secret_key", cloudAccount.Credentials.SecretKey)
	if err != nil {
		return err
	}

	err = data.Set("access_key_id", cloudAccount.Credentials.AccessKeyId)
	if err != nil {
		return err
	}

	return nil
}

func monitorCloudAccountForCostToResourceData(data *schema.ResourceData, cloudAccount *v2.CloudAccountCostProvider) error {
	err := data.Set("cloud_provider", cloudAccount.Provider)
	if err != nil {
		return err
	}

	err = data.Set("config", map[string]interface{}{
		"athena_bucket_name":      cloudAccount.Config.AthenaBucketName,
		"athena_database_name":    cloudAccount.Config.AthenaDatabaseName,
		"athena_region":           cloudAccount.Config.AthenaRegion,
		"athena_workgroup":        cloudAccount.Config.AthenaWorkgroup,
		"athena_table_name":       cloudAccount.Config.AthenaTableName,
		"spot_prices_bucket_name": cloudAccount.Config.SpotPricesBucketName,
	})
	if err != nil {
		return err
	}

	err = data.Set("integration_type", cloudAccount.Config.IntegrationType)
	if err != nil {
		return err
	}

	err = data.Set("account_id", cloudAccount.ProviderId)
	if err != nil {
		return err
	}

	err = data.Set("role_name", strings.Split(cloudAccount.RoleArn, ":role/")[1])
	if err != nil {
		return err
	}

	return nil
}

func completeEmptyFieldsForUpdateRequest(putObject *v2.CloudAccountCostProvider, currentStatusObject *v2.CloudAccountCostProvider) {
	putObject.ExternalId = currentStatusObject.ExternalId
	putObject.CredentialsType = currentStatusObject.CredentialsType
	putObject.CustomerId = currentStatusObject.CustomerId
	putObject.ProviderId = currentStatusObject.ProviderId
	putObject.CredentialsId = currentStatusObject.CredentialsId
	putObject.Feature = currentStatusObject.Feature
	putObject.Enabled = currentStatusObject.Enabled
}

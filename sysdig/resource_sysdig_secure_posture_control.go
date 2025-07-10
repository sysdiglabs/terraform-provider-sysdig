package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecurePostureControl() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecurePostureControlCreateOrUpdate,
		ReadContext:   resourceSysdigSecurePostureContorlRead,
		DeleteContext: resourceSysdigSecurePostureControlDelete,
		UpdateContext: resourceSysdigSecurePostureControlCreateOrUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaResourceKindKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaResourceRegoKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaResourceRemediationDetailsKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaResourceSeverityKey: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"High", "Medium", "Low"}, false),
			},
		},
	}
}

func resourceSysdigSecurePostureControlCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	// Extract 'group' field from Terraform configuration
	client, err := getPostureControlClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &v2.SaveControlRequest{
		ID:                 getStringValue(d, SchemaIDKey),
		Name:               getStringValue(d, SchemaNameKey),
		Description:        getStringValue(d, SchemaDescriptionKey),
		ResourceKind:       getStringValue(d, SchemaResourceKindKey),
		Rego:               getStringValue(d, SchemaResourceRegoKey),
		Severity:           getStringValue(d, SchemaResourceSeverityKey),
		RemediationDetails: getStringValue(d, SchemaResourceRemediationDetailsKey),
	}

	control, errStatus, err := client.CreateOrUpdatePostureControl(ctx, req)
	if err != nil {
		return diag.Errorf("Error saving control. error status: %s err: %s", errStatus, err)
	}

	d.SetId(control.ID)
	resourceSysdigSecurePostureContorlRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecurePostureContorlRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getPostureControlClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	control, err := client.GetPostureControlByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set(SchemaIDKey, control.ID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaNameKey, control.Name)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaDescriptionKey, control.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaResourceRemediationDetailsKey, control.RemediationDetails)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaResourceKindKey, control.ResourceKind)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaResourceRegoKey, control.Rego)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaResourceSeverityKey, control.Severity)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecurePostureControlDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getPostureControlClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeletePostureControlByID(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getPostureControlClient(c SysdigClients) (v2.PostureControlInterface, error) {
	var client v2.PostureControlInterface
	var err error
	switch c.GetClientType() {
	case IBMSecure:
		client, err = c.ibmSecureClient()
		if err != nil {
			return nil, err
		}
	default:
		client, err = c.sysdigSecureClientV2()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

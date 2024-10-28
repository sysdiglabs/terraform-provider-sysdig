package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureAcceptPostureRisk() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureAcceptPostureControlCreate,
		ReadContext:   resourceSysdigSecureAcceptPostureControlRead,
		DeleteContext: resourceSysdigSecureAcceptPostureControlDelete,
		UpdateContext: resourceSysdigSecureAcceptPostureControlUpdate,
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
			SchemaControlNameKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaZoneNameKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaDescriptionKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			SchemaFilterKey: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			SchemaReasonKey: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Risk Owned", "Risk Transferred", "Risk Avoided", "Risk Mitigated", "Risk Not Relevant", "Custom"}, false),
			},
			SchemaExpiresInKey: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"7 Days", "30 Days", "60 Days", "90 Days", "Custom", "Never"}, false),
			},
			SchemaExpiresAtKey: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			SchemaIsExpiredKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			SchemaAcceptanceDateKey: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			SchemaUsernameKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaTypeKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaIsSystemKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			SchemaAcceptPeriodKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSysdigSecureAcceptPostureControlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Extract 'group' field from Terraform configuration
	client, err := getPostureAcceptRiskClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &v2.AccepetPostureRiskRequest{
		AcceptanceID: d.Id(),
		ControlName:  d.Get(SchemaControlNameKey).(string),
		ZoneName:     d.Get(SchemaZoneNameKey).(string),
		Description:  d.Get(SchemaDescriptionKey).(string),
		Filter:       d.Get(SchemaFilterKey).(string),
		Reason:       d.Get(SchemaReasonKey).(string),
	}

	expiresIn := d.Get(SchemaExpiresInKey).(string)
	if expiresIn == "7 Days" {
		req.ExpiresAt = time.Now().AddDate(0, 0, 7).UTC().UnixMilli()
	} else if expiresIn == "30 Days" {
		req.ExpiresAt = time.Now().AddDate(0, 0, 30).UTC().UnixMilli()
	} else if expiresIn == "60 Days" {
		req.ExpiresAt = time.Now().AddDate(0, 0, 60).UTC().UnixMilli()
	} else if expiresIn == "90 Days" {
		req.ExpiresAt = time.Now().AddDate(0, 0, 90).UTC().UnixMilli()
	} else if expiresIn == "Never" {
		req.ExpiresAt = 0
	} else {
		req.ExpiresAt = d.Get(SchemaExpiresAtKey).(int64)
	}

	if req.ZoneName == "" && req.Filter == "" || req.ZoneName != "" && req.Filter != "" {
		return diag.Errorf("Error creating accept risk. Either a zone name must be provided to accept all resources for control in a specific zone, or a filter must be provided to accept a specific resource.")
	}

	acceptance, errStatus, err := client.SaveAcceptPostureRisk(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating accept risk. error status: %s err: %s", errStatus, err)
	}
	d.SetId(acceptance.Data.AcceptanceID)
	resourceSysdigSecureAcceptPostureControlRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecureAcceptPostureControlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Extract 'group' field from Terraform configuration
	client, err := getPostureAcceptRiskClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &v2.UpdateAccepetPostureRiskRequest{
		AcceptanceID: d.Id(),
		Acceptance:   v2.UpdateAcceptPostureRiskFields{},
	}
	expiresIn := d.Get(SchemaExpiresInKey).(string)
	var millis int64
	if expiresIn == "7 Days" {
		req.Acceptance.AcceptPeriod = "7"
		millis = time.Now().AddDate(0, 0, 7).UTC().UnixMilli()
	} else if expiresIn == "30 Days" {
		req.Acceptance.AcceptPeriod = "30"
		millis = time.Now().AddDate(0, 0, 30).UTC().UnixMilli()
	} else if expiresIn == "60 Days" {
		req.Acceptance.AcceptPeriod = "60"
		millis = time.Now().AddDate(0, 0, 60).UTC().UnixMilli()
	} else if expiresIn == "90 Days" {
		req.Acceptance.AcceptPeriod = "90"
		millis = time.Now().AddDate(0, 0, 90).UTC().UnixMilli()
	} else if expiresIn == "Never" {
		req.Acceptance.AcceptPeriod = "Never"
		millis = 0
	} else {
		req.Acceptance.AcceptPeriod = "Custom"
		req.Acceptance.ExpiresAt = d.Get(SchemaExpiresAtKey).(string)
	}
	req.Acceptance.ExpiresAt = fmt.Sprintf("%d", millis)
	req.Acceptance.Description = d.Get(SchemaDescriptionKey).(string)
	req.Acceptance.Reason = d.Get(SchemaReasonKey).(string)

	acceptance, errStatus, err := client.UpdateAcceptancePostureRisk(ctx, req)
	if err != nil {
		return diag.Errorf("Error updating accept risk. ID: %s, error status: %s err: %s", req.AcceptanceID, errStatus, err)
	}
	d.SetId(acceptance.AcceptanceID)
	resourceSysdigSecureAcceptPostureControlRead(ctx, d, meta)
	return nil
}

func resourceSysdigSecureAcceptPostureControlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPostureAcceptRiskClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	acceptance, errStatus, err := client.GetAcceptancePostureRisk(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if errStatus != "" {
		return diag.Errorf("Error reading accept risk. error status: %s", errStatus)
	}

	err = d.Set(SchemaControlNameKey, acceptance.Data.ControlName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaZoneNameKey, acceptance.Data.ZoneName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaDescriptionKey, acceptance.Data.Description)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaFilterKey, acceptance.Data.Filter)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaReasonKey, acceptance.Data.Reason)
	if err != nil {
		return diag.FromErr(err)
	}

	acceptanceDate, err := strconv.ParseInt(acceptance.Data.AcceeptanceDate, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set(SchemaAcceptanceDateKey, acceptanceDate)
	if err != nil {
		return diag.FromErr(err)
	}
	expiresAt, err := strconv.ParseInt(acceptance.Data.ExpiresAt, 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set(SchemaExpiresAtKey, expiresAt)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaUsernameKey, acceptance.Data.UserName)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaTypeKey, acceptance.Data.Type)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaIsSystemKey, acceptance.Data.IsSystem)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaAcceptPeriodKey, acceptance.Data.AcceptPeriod)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set(SchemaIsExpiredKey, acceptance.Data.IsExpired)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSysdigSecureAcceptPostureControlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPostureAcceptRiskClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}
	id := d.Id()
	err = client.DeleteAcceptancePostureRisk(ctx, &v2.DeleteAcceptPostureRisk{AcceptanceID: id})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getPostureAcceptRiskClient(c SysdigClients) (v2.PostureAcceptRiskInterface, error) {
	var client v2.PostureAcceptRiskInterface
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

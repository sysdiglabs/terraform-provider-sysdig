package sysdig

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSysdigSecureZonePosturePolicyAssignment() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureZonePosturePolicyAssignmentCreate,
		ReadContext:   resourceSysdigSecureZonePosturePolicyAssignmentRead,
		UpdateContext: resourceSysdigSecureZonePosturePolicyAssignmentUpdate,
		DeleteContext: resourceSysdigSecureZonePosturePolicyAssignmentDelete,
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
			SchemaZoneIDKey: {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			SchemaPolicyIDsKey: {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func getZonePolicyAssignmentClient(clients SysdigClients) (v2.ZonePolicyAssignmentInterface, error) {
	var client v2.ZonePolicyAssignmentInterface
	var err error
	switch clients.GetClientType() {
	case IBMSecure:
		client, err = clients.ibmSecureClient()
	default:
		client, err = clients.sysdigSecureClientV2()
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

func resourceSysdigSecureZonePosturePolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZonePolicyAssignmentClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID := d.Get(SchemaZoneIDKey).(int)
	req := &v2.ZonePolicyAssignmentRequest{
		PolicyIDs: expandIntSet(d.Get(SchemaPolicyIDsKey).(*schema.Set)),
	}

	_, err = client.CreateZonePolicyAssignment(ctx, zoneID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating zone policy assignment: %w", err))
	}

	d.SetId(strconv.Itoa(zoneID))
	return resourceSysdigSecureZonePosturePolicyAssignmentRead(ctx, d, m)
}

func resourceSysdigSecureZonePosturePolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZonePolicyAssignmentClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid zone id %q: %w", d.Id(), err))
	}

	assignment, err := client.GetZonePolicyAssignment(ctx, zoneID)
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading zone policy assignment for zone %d: %w", zoneID, err))
	}

	_ = d.Set(SchemaZoneIDKey, zoneID)
	_ = d.Set(SchemaPolicyIDsKey, assignment.PolicyIDs)
	return nil
}

func resourceSysdigSecureZonePosturePolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZonePolicyAssignmentClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid zone id %q: %w", d.Id(), err))
	}

	req := &v2.ZonePolicyAssignmentRequest{
		PolicyIDs: expandIntSet(d.Get(SchemaPolicyIDsKey).(*schema.Set)),
	}

	_, err = client.UpdateZonePolicyAssignment(ctx, zoneID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating zone policy assignment for zone %d: %w", zoneID, err))
	}

	return resourceSysdigSecureZonePosturePolicyAssignmentRead(ctx, d, m)
}

func resourceSysdigSecureZonePosturePolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client, err := getZonePolicyAssignmentClient(m.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	zoneID, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("invalid zone id %q: %w", d.Id(), err))
	}

	if err := client.DeleteZonePolicyAssignment(ctx, zoneID); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting zone policy assignment for zone %d: %w", zoneID, err))
	}

	d.SetId("")
	return nil
}

func expandIntSet(set *schema.Set) []int {
	result := make([]int, 0, set.Len())
	for _, v := range set.List() {
		result = append(result, v.(int))
	}
	return result
}

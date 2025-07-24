package sysdig

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSysdigSecureAutomation() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		CreateContext: resourceSysdigSecureAutomationCreate,
		ReadContext:   resourceSysdigSecureAutomationRead,
		UpdateContext: resourceSysdigSecureAutomationUpdate,
		DeleteContext: resourceSysdigSecureAutomationDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name for the automation (will override name in JSON)",
			},
			"automation_json": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "JSON configuration from Sysdig UI (name field will be replaced with 'name' field)",
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress diff if only the name changed (since we control the name)
					return suppressAutomationJSONDiff(old, new, d.Get("name").(string))
				},
			},
			"automation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The automation ID returned by Sysdig API",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the automation is enabled",
			},
			"customer_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Customer ID",
			},
			"team_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Team ID",
			},
			"author": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Author of the automation",
			},
			"group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Automation group",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the automation",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the automation was created",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the automation was last updated",
			},
			"last_executed_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When the automation was last executed",
			},
		},
	}
}

func getSecureAutomationClient(c SysdigClients) (v2.AutomationInterface, error) {
	return c.sysdigSecureClientV2()
}

func resourceSysdigSecureAutomationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureAutomationClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	// Prepare the JSON with Terraform name override
	automationJSON, err := prepareAutomationJSON(d, "")
	if err != nil {
		return diag.FromErr(err)
	}

	// Create the automation
	automation, err := client.CreateAutomation(ctx, automationJSON)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the resource ID to the automation ID
	d.SetId(automation.Automation.ID)

	// Update all computed fields
	err = automationToResourceData(&automation.Automation, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureAutomationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureAutomationClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	automationID := d.Id()
	automation, err := client.GetAutomationByID(ctx, automationID)
	if err != nil {
		// If not found, remove from state
		d.SetId("")
		return diag.FromErr(err)
	}

	// Update all fields from API response
	err = automationToResourceData(&automation.Automation, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureAutomationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureAutomationClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	automationID := d.Id()

	// Prepare the JSON with existing ID and Terraform name override
	automationJSON, err := prepareAutomationJSON(d, automationID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the automation
	automation, err := client.UpdateAutomation(ctx, automationID, automationJSON)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update all computed fields
	err = automationToResourceData(&automation.Automation, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigSecureAutomationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client, err := getSecureAutomationClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	automationID := d.Id()
	err = client.DeleteAutomation(ctx, automationID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// Helper function to prepare automation JSON with name override and ID injection to ensure terraform is the single source of truth
func prepareAutomationJSON(d *schema.ResourceData, existingID string) ([]byte, error) {
	jsonStr := d.Get("automation_json").(string)
	terraformName := d.Get("name").(string)

	// Parse the JSON
	var automation map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &automation); err != nil {
		return nil, fmt.Errorf("invalid automation_json: %v", err)
	}

	// Get the automation object
	automationObj, ok := automation["automation"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("automation_json must contain 'automation' object")
	}

	// OVERRIDE the name with Terraform name (this is the key sync point) and ensures terraforms name is what we keep
	automationObj["name"] = terraformName

	// For updates, ensure ID is present
	if existingID != "" {
		automationObj["id"] = existingID
	}

	return json.Marshal(automation)
}

// Simple helper function to copy automation details to resource data
func automationToResourceData(automation *v2.AutomationDetails, d *schema.ResourceData) error {
	_ = d.Set("automation_id", automation.ID)
	_ = d.Set("enabled", automation.Enabled)
	_ = d.Set("customer_id", automation.CustomerID)
	_ = d.Set("team_id", automation.TeamID)
	_ = d.Set("author", automation.Author)
	_ = d.Set("group", automation.Group)
	_ = d.Set("version", automation.Version)
	_ = d.Set("created_at", automation.CreatedAt)

	if automation.UpdatedAt != nil {
		_ = d.Set("updated_at", *automation.UpdatedAt)
	}
	if automation.LastExecutedAt != nil {
		_ = d.Set("last_executed_at", *automation.LastExecutedAt)
	}

	return nil
}

// Helper function to suppress diff when only the name field has changed in JSON
func suppressAutomationJSONDiff(old, new, terraformName string) bool {
	if old == "" || new == "" {
		return false
	}

	// Parse both JSON strings
	var oldAutomation, newAutomation map[string]interface{}
	if err := json.Unmarshal([]byte(old), &oldAutomation); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newAutomation); err != nil {
		return false
	}

	// Get automation objects
	oldObj, ok1 := oldAutomation["automation"].(map[string]interface{})
	newObj, ok2 := newAutomation["automation"].(map[string]interface{})
	if !ok1 || !ok2 {
		return false
	}

	// Set both names to the terraform name for comparison
	oldObj["name"] = terraformName
	newObj["name"] = terraformName

	// Remove computed fields that might differ
	delete(oldObj, "id")
	delete(oldObj, "updatedAt")
	delete(oldObj, "lastExecutedAt")
	delete(oldObj, "customerId")
	delete(oldObj, "teamId")
	delete(oldObj, "author")
	delete(oldObj, "createdAt")
	delete(oldObj, "validationError")
	delete(oldObj, "rateLimitDiscardedEvents")
	delete(oldObj, "latestRateLimitHit")
	delete(oldObj, "rateLimitDiscardedNodes")
	delete(oldObj, "latestNodeRateLimitHit")

	delete(newObj, "id")
	delete(newObj, "updatedAt")
	delete(newObj, "lastExecutedAt")
	delete(newObj, "customerId")
	delete(newObj, "teamId")
	delete(newObj, "author")
	delete(newObj, "createdAt")
	delete(newObj, "validationError")
	delete(newObj, "rateLimitDiscardedEvents")
	delete(newObj, "latestRateLimitHit")
	delete(newObj, "rateLimitDiscardedNodes")
	delete(newObj, "latestNodeRateLimitHit")

	// Compare the normalised objects
	return reflect.DeepEqual(oldObj, newObj)
}

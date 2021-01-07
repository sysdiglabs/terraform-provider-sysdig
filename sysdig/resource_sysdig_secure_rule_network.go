package sysdig

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func resourceSysdigSecureRuleNetwork() *schema.Resource {
	timeout := 30 * time.Second

	return &schema.Resource{
		CreateContext: resourceSysdigRuleNetworkCreate,
		UpdateContext: resourceSysdigRuleNetworkUpdate,
		ReadContext:   resourceSysdigRuleNetworkRead,
		DeleteContext: resourceSysdigRuleNetworkDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(timeout),
			Read:   schema.DefaultTimeout(timeout),
			Update: schema.DefaultTimeout(timeout),
			Delete: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleSchema(map[string]*schema.Schema{
			"block_inbound": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"block_outbound": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"tcp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"udp": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		}),
	}
}

func resourceSysdigRuleNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleNetworkFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err = client.CreateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(rule.ID))
	d.Set("version", rule.Version)

	return nil
}

// Retrieves the information of a resource form the file and loads it in Terraform
func resourceSysdigRuleNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := client.GetRuleByID(ctx, id)

	if err != nil {
		d.SetId("")
	}
	updateResourceDataForRule(d, rule)

	d.Set("block_inbound", rule.Details.AllInbound)
	d.Set("block_outbound", rule.Details.AllOutbound)

	if rule.Details.TCPListenPorts == nil {
		return diag.Errorf("no tcpListenPorts for a filesystem rule")
	}

	if rule.Details.UDPListenPorts == nil {
		return diag.Errorf("no udpListenPorts for a filesystem rule")
	}

	if len(rule.Details.TCPListenPorts.Items) > 0 {
		tcpPorts := []int{}
		for _, port := range rule.Details.TCPListenPorts.Items {
			intPort, err := strconv.Atoi(port)
			if err != nil {
				return diag.FromErr(err)
			}
			tcpPorts = append(tcpPorts, intPort)
		}
		d.Set("tcp", []map[string]interface{}{{
			"matching": rule.Details.TCPListenPorts.MatchItems,
			"ports":    tcpPorts,
		}})
	}
	if len(rule.Details.UDPListenPorts.Items) > 0 {
		udpPorts := []int{}
		for _, port := range rule.Details.UDPListenPorts.Items {
			intPort, err := strconv.Atoi(port)
			if err != nil {
				return diag.FromErr(err)
			}
			udpPorts = append(udpPorts, intPort)
		}
		d.Set("udp", []map[string]interface{}{{
			"matching": rule.Details.UDPListenPorts.MatchItems,
			"ports":    udpPorts,
		}})
	}

	return nil
}

func resourceSysdigRuleNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := resourceSysdigRuleNetworkFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	rule.Version = d.Get("version").(int)
	rule.ID, _ = strconv.Atoi(d.Id())

	_, err = client.UpdateRule(ctx, rule)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigRuleNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(SysdigClients).sysdigSecureClient()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteRule(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSysdigRuleNetworkFromResourceData(d *schema.ResourceData) (rule secure.Rule, err error) {
	rule = ruleFromResourceData(d)
	rule.Details.RuleType = "NETWORK"

	rule.Details.TCPListenPorts = &secure.TCPListenPorts{}
	rule.Details.UDPListenPorts = &secure.UDPListenPorts{}

	rule.Details.AllInbound = d.Get("block_inbound").(bool)
	rule.Details.AllOutbound = d.Get("block_outbound").(bool)

	rule.Details.TCPListenPorts.Items = []string{}
	if tcpRules := d.Get("tcp").([]interface{}); len(tcpRules) > 0 {
		rule.Details.TCPListenPorts.MatchItems = d.Get("tcp.0.matching").(bool)
		for _, port := range d.Get("tcp.0.ports").(*schema.Set).List() {
			portStr := port.(int)
			rule.Details.TCPListenPorts.Items = append(rule.Details.TCPListenPorts.Items, strconv.Itoa(portStr))
		}
	}

	rule.Details.UDPListenPorts.Items = []string{}
	if udpRules, ok := d.Get("udp").([]interface{}); ok && len(udpRules) > 0 {
		rule.Details.UDPListenPorts.MatchItems = d.Get("udp.0.matching").(bool)
		for _, port := range d.Get("udp.0.ports").(*schema.Set).List() {
			portStr := port.(int)
			rule.Details.UDPListenPorts.Items = append(rule.Details.UDPListenPorts.Items, strconv.Itoa(portStr))
		}
	}

	return
}

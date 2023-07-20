package sysdig

import (
	"context"
	"strconv"
	"time"

	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSysdigSecureRuleNetwork() *schema.Resource {
	timeout := 5 * time.Minute

	return &schema.Resource{
		ReadContext: dataSourceSysdigRuleNetworkRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(timeout),
		},

		Schema: createRuleDataSourceSchema(map[string]*schema.Schema{
			"block_inbound": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"block_outbound": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"tcp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"matching": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Computed: true,
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
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeSet,
							Computed: true,
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

func dataSourceSysdigRuleNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return commonDataSourceSysdigRuleRead(ctx, d, meta, v2.RuleTypeNetwork, networkRuleDataSourceToResourceData)
}

func networkRuleDataSourceToResourceData(rule v2.Rule, d *schema.ResourceData) diag.Diagnostics {
	_ = d.Set("block_inbound", !rule.Details.AllInbound)
	_ = d.Set("block_outbound", !rule.Details.AllOutbound)

	if rule.Details.TCPListenPorts == nil {
		return diag.Errorf("no tcpListenPorts for a network rule")
	}

	if rule.Details.UDPListenPorts == nil {
		return diag.Errorf("no udpListenPorts for a network rule")
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
		_ = d.Set("tcp", []map[string]interface{}{{
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
		_ = d.Set("udp", []map[string]interface{}{{
			"matching": rule.Details.UDPListenPorts.MatchItems,
			"ports":    udpPorts,
		}})
	}

	return nil
}

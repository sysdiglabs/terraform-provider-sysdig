package sysdig

import (
	"context"
	v2 "github.com/draios/terraform-provider-sysdig/sysdig/internal/client/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func dataSourceSysdigSecurePosturePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSysdigSecurePosturePoliciesRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			SchemaIDKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			SchemaPoliciesKey: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						SchemaIDKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaNameKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaTypeKey: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						SchemaKindKey: {
							Type:     schema.TypeInt,
							Computed: true,
						},
						SchemaDescriptionKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaVersionKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaApiVersionKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaLinkKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaAuthorsKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaPublishedDateKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaMinKubeVersionKey: {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						SchemaMaxKubeVersionKey: {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						SchemaIsCustomKey: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						SchemaIsActiveKey: {
							Type:     schema.TypeBool,
							Computed: true,
						},
						SchemaPlatformKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						SchemaZonesKey: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									SchemaIDKey: {
										Type:     schema.TypeString,
										Computed: true,
									},
									SchemaNameKey: {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func getPosturePolicyClient(c SysdigClients) (v2.PosturePolicyInterface, error) {
	var client v2.PosturePolicyInterface
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

func dataSourceSysdigSecurePosturePoliciesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := getPosturePolicyClient(meta.(SysdigClients))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListPosturePolicies(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	policies := make([]map[string]interface{}, len(resp))
	for i, p := range resp {
		zones := make([]map[string]interface{}, len(p.Zones))
		for j, z := range p.Zones {
			zones[j] = map[string]interface{}{
				SchemaIDKey:   z.ID,
				SchemaNameKey: z.Name,
			}
		}
		policies[i] = map[string]interface{}{
			SchemaIDKey:             p.ID,
			SchemaNameKey:           p.Name,
			SchemaTypeKey:           p.Type,
			SchemaKindKey:           p.Kind,
			SchemaDescriptionKey:    p.Description,
			SchemaVersionKey:        p.Version,
			SchemaApiVersionKey:     p.ApiVersion,
			SchemaLinkKey:           p.Link,
			SchemaAuthorsKey:        p.Authors,
			SchemaPublishedDateKey:  p.PublishedData,
			SchemaMinKubeVersionKey: p.MinKubeVersion,
			SchemaMaxKubeVersionKey: p.MaxKubeVersion,
			SchemaIsCustomKey:       p.IsCustom,
			SchemaIsActiveKey:       p.IsActive,
			SchemaPlatformKey:       p.Platform,
			SchemaZonesKey:          zones,
		}
	}

	err = d.Set(SchemaPoliciesKey, policies)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("0")
	return nil
}

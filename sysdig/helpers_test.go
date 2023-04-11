//go:build unit

package sysdig_test

import (
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyOnSchema(t *testing.T) {
	s := map[string]*schema.Schema{
		"key1": {},
		"key2": {
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key1": {},
					"key2": {
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key1": {},
							},
						},
					},
				},
			},
		},
	}

	sysdig.ApplyOnSchema(s, func(s *schema.Schema) {
		s.Required = true
	})

	assert.True(t, s["key1"].Required)
	assert.True(t, s["key2"].Required)
	nestedSchema := s["key2"].Elem.(*schema.Resource).Schema
	assert.True(t, nestedSchema["key1"].Required)
	assert.True(t, nestedSchema["key2"].Elem.(*schema.Resource).Schema["key1"].Required)
}

func TestMergeMap(t *testing.T) {
	m1 := map[string]int{
		"k1": 1,
		"k2": 2,
	}
	m2 := map[string]int{
		"k2": 3,
		"k3": 4,
	}
	sysdig.MergeMap(m1, m2)
	assert.Equal(t, 1, m1["k1"])
	assert.Equal(t, 3, m1["k2"])
	assert.Equal(t, 4, m1["k3"])
}

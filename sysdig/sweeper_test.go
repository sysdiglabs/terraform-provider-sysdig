package sysdig_test

import (
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var provider *schema.Provider

func TestMain(m *testing.M) {
	provider = sysdig.Provider()
	resource.TestMain(m)
}

//go:build tf_acc_sysdig || tf_acc_sysdig_secure

package sysdig_test

func TestAccManagedPolicyDataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_SECURE_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_SECURE_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: managedPolicyDataSource(),
			},
		},
	})
}

func managedPolicyDataSource() string {
	return `
data "sysdig_secure_managed_policy" "example" {
	name = "Sysdig Runtime Threat Detection"
	type = "falco"
}
`
}

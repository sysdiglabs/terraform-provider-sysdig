//go:build tf_acc_sysdig_monitor

package sysdig_test

import (
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

func TestCustomerProviderKeys(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"sysdig": func() (*schema.Provider, error) {
				return sysdig.Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:             monitorCustomerProviderKey(),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func monitorCustomerProviderKey() string {
	return `
resource "sysdig_monitor_cloud_account" "provider" {
 cloud_provider = "GCP"
 integration_type = "API"
 account_id = "joe-test-project-372418"
 additional_options = "ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiam9lLXRlc3QtcHJvamVjdC0zNzI0MTgiLAogICJwcml2YXRlX2tleV9pZCI6ICI2MzVlMTk0ZGNkODI4MWU5ZWE1YWZlMmJjNjdlMGIwYjY0OGI2YTM4IiwKICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXG5NSUlFdkFJQkFEQU5CZ2txaGtpRzl3MEJBUUVGQUFTQ0JLWXdnZ1NpQWdFQUFvSUJBUURkZlFrR0haN3d5Z3RKXG5kVnlKVHM5M3g5SldMeXQ4bTh2cHZNYzVzUDltMzRBQ1k0bTNNMGNxcE8yZWVkRTlWSURIeVR0Z25RZ3Fja2ZNXG5KTDlFZE9HQURHOEJlcFh5WmdiWjR2c2NvdDB4emFTTDRkQ3pvVDJqUEkvV3JEdFFaSGR4K1hOc0VpZFFzYjlQXG5qTndnN29DSE5XSkY3MXRNZCt3dUphejFMTlJOME5sRnU3b01BV3N0KzluQjIwS2FLZDZpTWY3Tlo1bko4N1lwXG5YV3hsM2RkSnRmV2RmTmgvSU1zZGhvSTBGMEhIMUIwdnBoOGl5L01kbDBsSExpeUNUdXVDS01PZ0NhNFRubENXXG5PcjJkaGJSSG5aV1ZtaXd4dVR4cllGRVRVeUZERDJiNVhjV3loQlh5ajJZNVpnVldqQWhmMlVyblJRTjJvMVVHXG44WVVBUzEwVEFnTUJBQUVDZ2dFQVF6MGJtYk14VnFrSGp5ZmxUVHZUS09wTkZPUGlBRVN0dlVvVmN4S2tIbDlZXG5WYUZSSkFBWnFUMERjL3BJUnFXYUtNeVN6WXd1ZC9CVWtvbFBWV0ZrT2NMTWlqYmtRWCt1c2NQQjl0b01hM3VoXG42ZU5HUDlvQnc4WDFacmJIbE9yREJpTXo0b21LVE9tQkNnM1puOWUzeGhRelBzYmd3UkNnN3d0NSs3NDl2MWRJXG5vVUh4MitXcDZ2d3ZQV0NVRDFHb0RDSElXWUthNG9kN2FvWmhnSDlhQm1ubkZ6dWYydXJuRTVuYUMzR1o3dUVsXG4zUCt0eVIvTlFLR2RZWm5yY3NTM1B6MWtoTFBMczFTL25kRk5MU1pOemVxMWpOajVRdWYxd05wdDBvYjJsZ1RGXG5jYmQrV2lOYWF5M0gvOW5DQXFJMzUzOUNQVlZpK0pNOEtPYlVYVjJ0WlFLQmdRRHdLeXM4bElJdFZSRE9yTHljXG4vaG1DdXZPdHJlUFR0TFB6MlFYZFZzNVhKanNEd3Z3MWJoeXc2VGN6UERNZXFSYUhQNUtCZW5kOEwwc3d4MjYxXG5HMGUyRXVGcnBSejk1VWNCNWFxSHkvcHp5MEtLQ0krbitIUURlUlNlSks2aytNdjJldzVpODhobUxtVTFVTW9QXG5kbDIraFBBWDRzL0dXTXQ2c3Vhb2kraXFKd0tCZ1FEc0ZxT1NHZmIvMGFXVzZwVSt0OWRxc21nRkV2R3NsQk9lXG5ZM0NLU3dlNDFhcW1ZcTRZaEp5bVg1dmtjWEl1SFlsVUtKdTNIOU9vbzBoeW12OHcxR2RwRHY5eDE1TUJQakd5XG5kOEZyeU1tK0FMVXpKaXgzSFJHcTl5VHd5cjdtaGE0Qm5kcmI4K0tTRk1hRFEzNWNkZ3JpUWgzdEdIOCsycWtrXG5TemMwWXdPbE5RS0JnR3UyNE1obHp0Q29FMGF1WUZXRS9Vb05zUmFYSTlRaWVvY0dNY1FvbDVpc2s5RkhGVGlkXG5idzdGT2pXbmJVSDJFaDJNbkplbnBva3k2T1V5dk90TEZlbUtKRUhVSnVHVWdEbFFtU0FZa3ZaMkZoeTBaRUd3XG5nOCsrOFVsUUtHZmpFZzgwOTZuWHJteHRxSVMxL0RuZEc0UkVPUzV0VWtTaU5IaU9YamIvc05VSEFvR0FFbzVJXG45dS9CZ1NQYUx2MXJFNDNoaVlwU01Ldm5nTmYybnNsVURCcVBsZEI5WkN4M1lJZnp4QVBadmQvSXlLVWJxUml6XG4vSFdzN2lFL1RYcXZPZ2hIeEhNZ1VyTk40NWdlMGRjbHhiSDNZVTZ1NzBFOTEzTGFjNlNQSzduVHZVeWVlNVFMXG5vcVFObDh1NE9wTHdlSlh5and3QlRDUlR3LzN0czJPU0NEVU1FVTBDZ1lCdW9ocGdMMUxGT2dXc3R6M1B5UUVtXG5POHdvRVFtelVINFU3SUx4U1VCVkJ5dnpqa1NZL2Z3NzVQc1NQMUl2eS9DYi9yQWRoRFI2MU9taVk5RWFJQXBmXG42YmN5cDAySElrNGhKb0IybTE4RmNYejQ5cGFPUncxb0I2bkVDTXg5Tk5jbGx4cFV4ak9SSk1idGxZV2Y3SUZ4XG5jTlQySDljQUdZVDZ2c2RsekJIOFd3PT1cbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIsCiAgImNsaWVudF9lbWFpbCI6ICJ0ZXN0LTg3MUBqb2UtdGVzdC1wcm9qZWN0LTM3MjQxOC5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgImNsaWVudF9pZCI6ICIxMTcwMDIxNjYzNjIzODM2MzMyMzMiLAogICJhdXRoX3VyaSI6ICJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9vYXV0aDIvYXV0aCIsCiAgInRva2VuX3VyaSI6ICJodHRwczovL29hdXRoMi5nb29nbGVhcGlzLmNvbS90b2tlbiIsCiAgImF1dGhfcHJvdmlkZXJfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLAogICJjbGllbnRfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9yb2JvdC92MS9tZXRhZGF0YS94NTA5L3Rlc3QtODcxJTQwam9lLXRlc3QtcHJvamVjdC0zNzI0MTguaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0="
}
`
}

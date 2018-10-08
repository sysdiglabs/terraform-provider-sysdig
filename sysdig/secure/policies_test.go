package secure_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func TestCreatePolicy(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	policy, err := sysdigSecureClient.CreatePolicy(aPolicy())

	assert.NotEqual(t, policy.ID, 0)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(policy.ID)
}

func TestCreatePolicyFailsWhenPolicyDoesNotHaveAllRequiredFields(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	_, err := sysdigSecureClient.CreatePolicy(aPolicyWithoutNameAndDescription())

	assert.NotNil(t, err)
}

func TestUpdatePolicy(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	created, err := sysdigSecureClient.CreatePolicy(aPolicy())
	created.Severity = 5

	updated, err := sysdigSecureClient.UpdatePolicy(created)

	assert.Equal(t, created.Severity, updated.Severity)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(created.ID)
}

func TestUpdatePolicyFailsWhenPolicyDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentID := 9838
	policy := aPolicy()
	policy.ID = nonExistentID

	_, err := sysdigSecureClient.UpdatePolicy(policy)

	assert.NotNil(t, err)
}

func TestGetPolicyById(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	created, err := sysdigSecureClient.CreatePolicy(aPolicy())

	retrieved, err := sysdigSecureClient.GetPolicyById(created.ID)

	assert.Equal(t, created, retrieved)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(created.ID)
}

func TestGetPolicyByIdFailsWhenPolicyDoesNotExist(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	nonExistentID := 9838

	_, err := sysdigSecureClient.GetPolicyById(nonExistentID)

	assert.NotNil(t, err)
}

func aPolicy() secure.Policy {
	return secure.Policy{
		Name:           "Write dpkg stuff",
		Description:    "an attempt to write to the dpkg database by any non-dpkg related program",
		Severity:       4,
		ContainerScope: true,
		HostScope:      true,
		Enabled:        true,
		Scope:          "host.ip.private = \"10.0.0.1\"",
		FalcoConfiguration: secure.FalcoConfiguration{
			RuleNameRegEx: "Mkdir binary dirs",
		},
		Actions: []secure.Action{
			secure.Action{
				Type: "POLICY_ACTION_PAUSE",
			},
			secure.Action{
				Type:                 "POLICY_ACTION_CAPTURE",
				AfterEventNs:         10000000000,
				BeforeEventNs:        10000000000,
				IsLimitedToContainer: false,
			},
		},
	}
}

func aPolicyWithoutNameAndDescription() secure.Policy {
	return secure.Policy{
		Severity:       4,
		ContainerScope: true,
		HostScope:      true,
		Enabled:        true,
	}
}

package sysdig_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig"
)

// == Policies ================================================================

func TestCreatePolicy(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	policy, err := sysdigSecureClient.CreatePolicy(aPolicy())

	assert.NotEqual(t, policy.ID, 0)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(policy.ID)
}

func TestCreatePolicyFailsWhenPolicyDoesNotHaveAllRequiredFields(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	_, err := sysdigSecureClient.CreatePolicy(aPolicyWithoutNameAndDescription())

	assert.NotNil(t, err)
}

func TestUpdatePolicy(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	created, err := sysdigSecureClient.CreatePolicy(aPolicy())
	created.Severity = 5

	updated, err := sysdigSecureClient.UpdatePolicy(created)

	assert.Equal(t, created.Severity, updated.Severity)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(created.ID)
}

func TestUpdatePolicyFailsWhenPolicyDoesNotExist(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentID := 9838
	policy := aPolicy()
	policy.ID = nonExistentID

	_, err := sysdigSecureClient.UpdatePolicy(policy)

	assert.NotNil(t, err)
}

func TestGetPolicyById(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	created, err := sysdigSecureClient.CreatePolicy(aPolicy())

	retrieved, err := sysdigSecureClient.GetPolicyById(created.ID)

	assert.Equal(t, created, retrieved)
	assert.Nil(t, err)

	// Cleanup Sysdig Secure
	defer sysdigSecureClient.DeletePolicy(created.ID)
}

func TestGetPolicyByIdFailsWhenPolicyDoesNotExist(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	nonExistentID := 9838

	_, err := sysdigSecureClient.GetPolicyById(nonExistentID)

	assert.NotNil(t, err)
}

func aPolicy() sysdig.Policy {
	return sysdig.Policy{
		Name:           "Write dpkg stuff",
		Description:    "an attempt to write to the dpkg database by any non-dpkg related program",
		Severity:       4,
		ContainerScope: true,
		HostScope:      true,
		Enabled:        true,
		Scope:          "host.ip.private = \"10.0.0.1\"",
		FalcoConfiguration: sysdig.FalcoConfiguration{
			RuleNameRegEx: "Mkdir binary dirs",
		},
		Actions: []sysdig.Action{
			sysdig.Action{
				Type: "POLICY_ACTION_PAUSE",
			},
			sysdig.Action{
				Type:                 "POLICY_ACTION_CAPTURE",
				AfterEventNs:         10000000000,
				BeforeEventNs:        10000000000,
				IsLimitedToContainer: false,
			},
		},
	}
}

func aPolicyWithoutNameAndDescription() sysdig.Policy {
	return sysdig.Policy{
		Severity:       4,
		ContainerScope: true,
		HostScope:      true,
		Enabled:        true,
	}
}

// == User Rules Files ========================================================

func TestCreateUserRulesFile(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	userRulesFile, err := sysdigSecureClient.GetUserRulesFile()
	updatedUserRulesFile, err := sysdigSecureClient.UpdateUserRulesFile(userRulesFile)

	assert.Equal(t, updatedUserRulesFile.Content, userRulesFile.Content)
	assert.Equal(t, updatedUserRulesFile.Version, userRulesFile.Version+1)
	assert.Nil(t, err)
}

// == Notification Channels ===================================================

func aNotificationChannel() sysdig.NotificationChannel {
	return sysdig.NotificationChannel{
		Name:    "Example Channel",
		Enabled: true,
		Type:    "EMAIL",
		Options: sysdig.NotificationChannelOptions{
			EmailRecipients: []string{"root@localhost.com"},
		},
	}
}

func aNotificationChannelWithoutRecipients() sysdig.NotificationChannel {
	return sysdig.NotificationChannel{
		Name:    "Example channel",
		Enabled: true,
		Type:    "EMAIL",
	}
}

func TestCreateNotificationChannel(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)
	assert.NotEqual(t, 0, channel.ID)
}

func TestCreateNotificationChannelWithoutRecipientsFails(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	_, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannelWithoutRecipients())
	assert.NotNil(t, err)
}

func TestUpdateNotificationChannel(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)
	assert.Equal(t, "Example Channel", channel.Name)

	channel.Name = "Changed Name"
	newChannel, err := sysdigSecureClient.UpdateNotificationChannel(channel)

	assert.Nil(t, err)
	assert.Equal(t, "Changed Name", newChannel.Name)

}

func TestUpdateNotificationChannelFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)

	nonExistentId := 1
	channel.ID = nonExistentId
	channel.Name = "Changed Name"
	_, err = sysdigSecureClient.UpdateNotificationChannel(channel)

	assert.NotNil(t, err)

}

func TestGetNotificationChannelById(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)

	newChannel, err := sysdigSecureClient.GetNotificationChannelById(channel.ID)

	assert.Nil(t, err)
	assert.Equal(t, channel.ID, newChannel.ID)
	assert.Equal(t, channel.Version, newChannel.Version)
	assert.Equal(t, channel.Name, newChannel.Name)
	assert.Equal(t, channel.Type, newChannel.Type)

	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)
}

func TestGetNotificationChannelByIdFailsWhenDoesNotExist(t *testing.T) {
	sysdigSecureClient := sysdig.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentId := 1
	_, err := sysdigSecureClient.GetNotificationChannelById(nonExistentId)

	assert.NotNil(t, err)
}

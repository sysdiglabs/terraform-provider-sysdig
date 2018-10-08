package secure_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

func TestCreateNotificationChannel(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)
	assert.NotEqual(t, 0, channel.ID)
}

func TestCreateNotificationChannelWithoutRecipientsFails(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	_, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannelWithoutRecipients())
	assert.NotNil(t, err)
}

func TestUpdateNotificationChannel(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
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
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	channel, err := sysdigSecureClient.CreateNotificationChannel(aNotificationChannel())
	assert.Nil(t, err)
	defer sysdigSecureClient.DeleteNotificationChannel(channel.ID)

	nonExistentID := 1
	channel.ID = nonExistentID
	channel.Name = "Changed Name"
	_, err = sysdigSecureClient.UpdateNotificationChannel(channel)

	assert.NotNil(t, err)

}

func TestGetNotificationChannelById(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")
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
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	nonExistentID := 1
	_, err := sysdigSecureClient.GetNotificationChannelById(nonExistentID)

	assert.NotNil(t, err)
}

func aNotificationChannel() secure.NotificationChannel {
	return secure.NotificationChannel{
		Name:    "Example Channel",
		Enabled: true,
		Type:    "EMAIL",
		Options: secure.NotificationChannelOptions{
			EmailRecipients: []string{"root@localhost.com"},
		},
	}
}

func aNotificationChannelWithoutRecipients() secure.NotificationChannel {
	return secure.NotificationChannel{
		Name:    "Example channel",
		Enabled: true,
		Type:    "EMAIL",
	}
}

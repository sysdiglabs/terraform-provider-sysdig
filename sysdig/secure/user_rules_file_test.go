package secure_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/draios/terraform-provider-sysdig/sysdig/secure"
)

// == User Rules Files ========================================================

func TestCreateUserRulesFile(t *testing.T) {
	sysdigSecureClient := secure.NewSysdigSecureClient(os.Getenv("SYSDIG_SECURE_API_TOKEN"), "https://secure.sysdig.com")

	userRulesFile, err := sysdigSecureClient.GetUserRulesFile()
	updatedUserRulesFile, err := sysdigSecureClient.UpdateUserRulesFile(userRulesFile)

	assert.Equal(t, updatedUserRulesFile.Content, userRulesFile.Content)
	assert.Equal(t, updatedUserRulesFile.Version, userRulesFile.Version+1)
	assert.Nil(t, err)
}

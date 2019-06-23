package gomniauth_traq

import (
	"testing"

	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/oauth2"
	"github.com/stretchr/gomniauth/test"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserInterface(t *testing.T) {

	var user common.User = new(User)

	assert.NotNil(t, user)

}

func TestNewUser(t *testing.T) {

	testProvider := new(test.TestProvider)
	testProvider.On("Name").Return("providerName")

	data := objx.MSI(
		traqKeyID, "123435467890",
		traqKeyName, "Mathew",
		traqKeyDisplayName, "Mat",
		traqKeyIconFileId, "aaaa")
	creds := &common.Credentials{Map: objx.MSI(oauth2.OAuth2KeyAccessToken, "ABC123")}

	user := NewUser(data, creds, testProvider)

	if assert.NotNil(t, user) {

		assert.Equal(t, data, user.Data())

		assert.Equal(t, "Mathew", user.Name())
		assert.Equal(t, "", user.Email())
		assert.Equal(t, "Mat", user.Nickname())
		assert.Equal(t, traqUserIconPath+"aaaa", user.AvatarURL())

		// check provider credentials
		creds := user.ProviderCredentials()[testProvider.Name()]
		if assert.NotNil(t, creds) {
			assert.Equal(t, "ABC123", creds.Get(oauth2.OAuth2KeyAccessToken).Str())
			assert.Equal(t, "123435467890", creds.Get(common.CredentialsKeyID).Str())
		}

	}

	mock.AssertExpectationsForObjects(t, testProvider.Mock)

}

func TestIDForProvider(t *testing.T) {

	user := new(User)
	user.data = objx.MSI(
		common.UserKeyProviderCredentials,
		map[string]*common.Credentials{
			"traq": &common.Credentials{Map: objx.MSI(common.CredentialsKeyID, "traqid")}})

	assert.Equal(t, "traqid", user.IDForProvider("traq"))

}

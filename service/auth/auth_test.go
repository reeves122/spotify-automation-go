package auth

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/reeves122/spotify-automation-go/mocks/mock_adapter"
	"github.com/reeves122/spotify-automation-go/service/storage"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

var testToken = &oauth2.Token{
	AccessToken:  "test1",
	RefreshToken: "test2",
	TokenType:    "Bearer",
	Expiry:       time.Unix(1644696995, 0),
}

func cleanUp(cacheDir string) {
	_ = os.RemoveAll(cacheDir)
	os.Clearenv()
}

func Test_NewAuth(t *testing.T) {
	assert.NotNil(t, NewAuth(nil, nil))
}

// Test_CheckForResponseUrl_Missing tests with RESPONSE_CODE missing
func Test_CheckForResponseUrl_Missing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	a := auth{spotify: mockWrapper, storage: s}

	mockWrapper.EXPECT().GetAuthURL().Return("https://dummyurl")

	assert.Error(t, a.checkForResponseUrl())
}

// Test_CheckForResponseUrl_Present tests with RESPONSE_CODE present
func Test_CheckForResponseUrl_Present(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	a := auth{spotify: mockWrapper, storage: s}

	_ = os.Setenv("RESPONSE_CODE", "abc123")

	assert.NoError(t, a.checkForResponseUrl())
}

// Test_CreateAndSaveToken_Success tests saving a token without error
func Test_CreateAndSaveToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	defer cleanUp("test")
	a := auth{spotify: mockWrapper, storage: s}

	_ = os.Setenv("RESPONSE_CODE", "abc123")
	mockWrapper.EXPECT().GetTokenFromResponseCode("abc123").Return(testToken, nil)

	result, err := a.createAndSaveToken("token.json")
	assert.NoError(t, err)
	assert.Equal(t, testToken, result)
}

// Test_CreateAndSaveToken_MissingEnv tests with RESPONSE_CODE missing
func Test_CreateAndSaveToken_MissingEnv(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	defer cleanUp("test")
	a := auth{spotify: mockWrapper, storage: s}

	mockWrapper.EXPECT().GetTokenFromResponseCode("").Return(nil, fmt.Errorf("test error"))

	result, err := a.createAndSaveToken("token.json")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// Test_Login_WithoutToken tests without loading a previous token
func Test_Login_WithoutToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	defer cleanUp("test")
	a := auth{spotify: mockWrapper, storage: s}

	_ = os.Setenv("RESPONSE_CODE", "abc123")
	mockWrapper.EXPECT().CreateAuthenticator("http://test")
	mockWrapper.EXPECT().GetTokenFromResponseCode("abc123").Return(testToken, nil)
	mockWrapper.EXPECT().LoginAndCreateClient(testToken)
	mockWrapper.EXPECT().GetToken().Return(testToken, nil)

	err := a.Login("http://test", "token.json")
	assert.NoError(t, err)
}

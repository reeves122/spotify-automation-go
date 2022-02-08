package auth

import (
	"os"

	"github.com/reeves122/spotify-automation-go/adapter"
	"github.com/reeves122/spotify-automation-go/service"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type auth struct {
	spotify adapter.SpotifyWrapperInterface
	storage service.StorageInterface
}

func NewAuth(spotify adapter.SpotifyWrapperInterface, storage service.StorageInterface) *auth {
	return &auth{
		spotify: spotify,
		storage: storage,
	}
}

func (a *auth) Login(redirectURL string, tokenFile string) error {
	a.spotify.CreateAuthenticator(redirectURL)

	token, err := a.storage.LoadToken(tokenFile)
	if err != nil {
		a.checkForResponseUrl()
		token = a.createAndSaveToken(tokenFile)
	}

	log.Info("Logging in using token")

	a.spotify.LoginAndCreateClient(token)
	return nil
}

func (a *auth) checkForResponseUrl() {
	if responseUrl := os.Getenv("RESPONSE_CODE"); responseUrl == "" {

		log.Info("Access token not found. Please use the below URL to authorize this " +
			"application and then set the RESPONSE_CODE env variable to the code " +
			"spotify responds with and run this application again")

		log.Info(a.spotify.GetAuthURL())
		os.Exit(0)
	}
}

func (a *auth) createAndSaveToken(tokenFile string) *oauth2.Token {
	log.Info("Attempting to get token using RESPONSE_CODE")
	token, err := a.spotify.GetTokenFromResponseCode(os.Getenv("RESPONSE_CODE"))
	if err != nil {
		log.Error("Unable to get token")
		panic(err)
	}

	log.Infof("Saving token to file: %s\n", tokenFile)
	err = a.storage.SaveToken(token, tokenFile)
	if err != nil {
		log.Error("Unable to save token to file: ", tokenFile)
		panic(err)
	}
	return token
}

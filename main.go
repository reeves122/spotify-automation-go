package main

import (
	"os"

	"github.com/reeves122/spotify-automation-go/adapter/spotifywrapper"
	"github.com/reeves122/spotify-automation-go/service/auth"
	"github.com/reeves122/spotify-automation-go/service/storage"
	"github.com/reeves122/spotify-automation-go/service/util"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func main() {
	log.SetLevel(log.DebugLevel)
	username := checkAndGetEnv("USER_NAME")
	redirectURL := checkAndGetEnv("REDIRECT_URL")
	tokenFile := checkAndGetEnv("TOKEN_FILE")
	cacheDir := checkAndGetEnv("CACHE_DIR")
	dislikedPrefix := checkAndGetEnv("DISLIKED_PREFIX")
	queueSuffix := checkAndGetEnv("QUEUE_SUFFIX")
	_ = checkAndGetEnv("SPOTIFY_ID")
	_ = checkAndGetEnv("SPOTIFY_SECRET")

	wrapper := spotifywrapper.NewWrapper(spotify.Client{}, spotifyauth.Authenticator{})

	storageService := storage.NewStorage(cacheDir, false)
	authService := auth.NewAuth(wrapper, storageService)
	err := authService.Login(redirectURL, tokenFile)
	if err != nil {
		panic(err)
	}

	utilService := util.NewUtil(wrapper, storageService, dislikedPrefix, queueSuffix)
	playlists, err := utilService.GetAllPlaylistsForUser(username)
	if err != nil {
		panic(err)
	}

	err = utilService.UpdateLocalCache(playlists)
	if err != nil {
		panic(err)
	}

	disliked, err := utilService.LoadAllDislikedTracks(playlists)
	if err != nil {
		panic(err)
	}

	err = utilService.ScanPlaylistsForDislikedTracks(playlists, disliked, username)
	if err != nil {
		panic(err)
	}

	err = utilService.ProcessQueuePlaylists(playlists, username)
	if err != nil {
		panic(err)
	}

	err = utilService.FindPossibleDuplicateTracks(playlists)
	if err != nil {
		panic(err)
	}

	log.Info("Done processing!")
}

func checkAndGetEnv(envVar string) string {
	value := os.Getenv(envVar)
	if value == "" {
		log.Errorf("%s env variable must be set", envVar)
		os.Exit(1)
	}
	log.Debugf("%s is set to '%s'", envVar, value)
	return value
}

package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type storage struct {
	cacheDir string
}

func NewStorage(cacheDir string, relativePath bool) *storage {
	if relativePath {
		cwd, _ := os.Getwd()
		cacheDir = filepath.Join(cwd, cacheDir)
	}

	createCacheDir(cacheDir)

	return &storage{
		cacheDir: cacheDir,
	}
}

func createCacheDir(cacheDir string) {
	_ = os.Mkdir(cacheDir, 0770)
}

// LoadToken loads the auth token from JSON file and parses it
func (s *storage) LoadToken(fileName string) (*oauth2.Token, error) {
	fileName = filepath.Join(s.cacheDir, fileName)
	log.Infof("Loading auth token from file: %s", fileName)

	rawFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer rawFile.Close()

	bytes, _ := ioutil.ReadAll(rawFile)

	var token oauth2.Token
	err = json.Unmarshal(bytes, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// SaveToken saves the auth token to JSON file
func (s *storage) SaveToken(token *oauth2.Token, fileName string) error {
	jsonData, _ := json.MarshalIndent(token, "", " ")
	fileName = filepath.Join(s.cacheDir, fileName)
	log.Debugf("Saving auth token to file: %s", fileName)
	err := os.WriteFile(fileName, jsonData, 0644)
	return err
}

// LoadTracksFile loads the playlist tracks from JSON file
func (s *storage) LoadTracksFile(playlistName string) ([]spotify.PlaylistTrack, error) {
	fileName := s.getPlaylistFilename(playlistName)
	log.Debugf("Loading playlist %s from file: %s", playlistName, fileName)

	var tracks []spotify.PlaylistTrack

	rawFile, err := os.Open(fileName)
	if _, ok := err.(*os.PathError); ok {
		return tracks, nil
	}
	if err != nil {
		return nil, err
	}
	defer rawFile.Close()

	bytes, _ := ioutil.ReadAll(rawFile)

	err = json.Unmarshal(bytes, &tracks)
	if err != nil {
		return nil, err
	}

	log.Debugf("Loaded %d cached tracks for playlist: %s", len(tracks), playlistName)
	return tracks, nil
}

// SaveTracksFile saves the playlist tracks to JSON file
func (s *storage) SaveTracksFile(playlistName string, tracks []spotify.PlaylistTrack) error {
	jsonData, _ := json.MarshalIndent(tracks, "", " ")
	fileName := s.getPlaylistFilename(playlistName)
	log.Debugf("Saving playlist %s to file %s with %d tracks", playlistName, fileName, len(tracks))

	err := os.WriteFile(fileName, jsonData, 0644)
	return err
}

// getPlaylistFilename returns the full path to a playlist file
func (s *storage) getPlaylistFilename(playlistName string) string {
	return filepath.Join(s.cacheDir, playlistName+".json")
}

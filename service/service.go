package service

import (
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type StorageInterface interface {
	LoadToken(fileName string) (*oauth2.Token, error)
	SaveToken(token *oauth2.Token, fileName string) error
	LoadTracksFile(playlistName string) ([]spotify.PlaylistTrack, error)
	SaveTracksFile(playlistName string, tracks []spotify.PlaylistTrack) error
}

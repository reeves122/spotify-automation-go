package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

var testTracks = []spotify.PlaylistTrack{
	{
		Track: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				Name: "track 1",
				Artists: []spotify.SimpleArtist{
					{
						Name: "artist 1",
					},
				},
			},
			Album: spotify.SimpleAlbum{
				Name: "album1",
			},
		},
	},
}

func cleanUp(cacheDir string) {
	_ = os.RemoveAll(cacheDir)
}

func Test_SaveTracksFile_RelativePath(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)
}

func Test_SaveTracksFile_AbsolutePath(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test")
	defer cleanUp(path)
	s := NewStorage(path, false)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)
}

func Test_LoadTracksFile(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)

	result, err := s.LoadTracksFile("test playlist")
	assert.Nil(t, err)
	assert.Equal(t, testTracks, result)
}

func Test_GetPlaylistFilename(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	playlist := "foo"
	cwd, _ := os.Getwd()
	result := s.getPlaylistFilename(playlist)
	assert.Equal(t, filepath.Join(cwd, "test", "foo.json"), result)
}

func Test_SaveToken(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	token := &oauth2.Token{
		AccessToken:  "",
		RefreshToken: "",
		TokenType:    "",
	}
	assert.NoError(t, s.SaveToken(token, "test.json"))
}

func Test_LoadToken(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	token := &oauth2.Token{
		AccessToken:  "",
		RefreshToken: "",
		TokenType:    "",
	}
	_ = s.SaveToken(token, "test.json")

	result, err := s.LoadToken("test.json")
	assert.NoError(t, err)
	assert.Equal(t, token, result)
}

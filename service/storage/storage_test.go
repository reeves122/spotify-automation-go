package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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

var testToken = &oauth2.Token{
	AccessToken:  "test1",
	RefreshToken: "test2",
	TokenType:    "Bearer",
	Expiry:       time.Unix(1644696995, 0).UTC(),
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

func Test_SaveTracksFile_With_Special_Characters(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test")
	defer cleanUp(path)
	s := NewStorage(path, false)

	// Test forward slash
	err := s.SaveTracksFile("test / playlist", testTracks)
	assert.Nil(t, err)

	// Test backslash
	err = s.SaveTracksFile("test \\ playlist", testTracks)
	assert.Nil(t, err)

	// Test period
	err = s.SaveTracksFile("test . playlist", testTracks)
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

func Test_LoadTracksFile_With_Special_Characters(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)

	// Test forward slash
	err := s.SaveTracksFile("test / playlist", testTracks)
	assert.Nil(t, err)

	result, err := s.LoadTracksFile("test / playlist")
	assert.Nil(t, err)
	assert.Equal(t, testTracks, result)

	// Test backslash
	err = s.SaveTracksFile("test \\ playlist", testTracks)
	assert.Nil(t, err)

	result, err = s.LoadTracksFile("test \\ playlist")
	assert.Nil(t, err)
	assert.Equal(t, testTracks, result)

	// Test period
	err = s.SaveTracksFile("test . playlist", testTracks)
	assert.Nil(t, err)

	result, err = s.LoadTracksFile("test . playlist")
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

func Test_GetPlaylistFilename_With_Special_Characters(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)

	// Test forward slash
	playlist := "test / playlist"
	cwd, _ := os.Getwd()
	result := s.getPlaylistFilename(playlist)
	assert.Equal(t, filepath.Join(cwd, "test", "test - playlist.json"), result)

	// Test backslash
	playlist = "test \\ playlist"
	cwd, _ = os.Getwd()
	result = s.getPlaylistFilename(playlist)
	assert.Equal(t, filepath.Join(cwd, "test", "test - playlist.json"), result)

	// Test period
	playlist = "test . playlist"
	cwd, _ = os.Getwd()
	result = s.getPlaylistFilename(playlist)
	assert.Equal(t, filepath.Join(cwd, "test", "test - playlist.json"), result)
}

func Test_SaveToken(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)
	assert.NoError(t, s.SaveToken(testToken, "test.json"))
}

func Test_LoadToken(t *testing.T) {
	defer cleanUp("test")
	s := NewStorage("test", true)

	_ = s.SaveToken(testToken, "test.json")
	result, err := s.LoadToken("test.json")
	assert.NoError(t, err)
	assert.Equal(t, testToken, result)
}

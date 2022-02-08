package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
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

func Test_SaveTracksFile_RelativePath(t *testing.T) {
	s := NewStorage("test", true)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)
}

func Test_SaveTracksFile_AbsolutePath(t *testing.T) {
	path := filepath.Join(os.TempDir(), "test")
	s := NewStorage(path, false)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)
}

func Test_LoadTracksFile(t *testing.T) {
	s := NewStorage("test", true)
	err := s.SaveTracksFile("test playlist", testTracks)
	assert.Nil(t, err)

	result, err := s.LoadTracksFile("test playlist")
	assert.Nil(t, err)
	assert.Equal(t, testTracks, result)
}

package util

import (
	"github.com/golang/mock/gomock"
	"github.com/reeves122/spotify-automation-go/mocks/mock_adapter"
	"github.com/reeves122/spotify-automation-go/service/storage"
	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify/v2"
	"testing"
)

func Test_util_NewUtil(t *testing.T) {
	assert.NotNil(t, NewUtil(nil, nil, "", ""))
}

func Test_util_GetAllPlaylistsForUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWrapper := mock_adapter.NewMockSpotifyWrapperInterface(ctrl)
	s := storage.NewStorage("test", true)
	u := util{spotify: mockWrapper, storage: s, dislikedPrefix: "disliked_", queueSuffix: " Queue"}

	mockWrapper.EXPECT().GetAllPlaylistsForUser("user").Return([]spotify.SimplePlaylist{
		{
			ID: spotify.ID("123"),
		},
	}, nil)
	result, err := u.GetAllPlaylistsForUser("user")
	assert.NoError(t, err)
	assert.Len(t, result, 1)

}

func Test_util_FindPossibleDuplicateTracks(t *testing.T) {

}

func Test_util_GetAllPlaylistsForUser1(t *testing.T) {

}

func Test_util_LoadAllDislikedTracks(t *testing.T) {

}

func Test_util_ProcessQueuePlaylists(t *testing.T) {

}

func Test_util_ScanPlaylistsForDislikedTracks(t *testing.T) {

}

func Test_util_UpdateLocalCache(t *testing.T) {

}

func Test_util_processQueuePlaylist(t *testing.T) {
}

func Test_util_scanPlaylistForDislikedTracks(t *testing.T) {

}

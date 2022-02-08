package spotifywrapper

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

const state = "spotify-automation-go"

type wrapper struct {
	client *spotify.Client
	auth   *spotifyauth.Authenticator
}

func NewWrapper(client spotify.Client, auth spotifyauth.Authenticator) *wrapper {
	return &wrapper{
		client: &client,
		auth:   &auth,
	}
}

func (w *wrapper) CreateAuthenticator(redirectURL string) {
	var redirect = spotifyauth.WithRedirectURL(redirectURL)
	var scopes = spotifyauth.WithScopes(
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopeUserLibraryModify,
		spotifyauth.ScopeUserReadRecentlyPlayed,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistModifyPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
	)
	w.auth = spotifyauth.New(redirect, scopes)
}

func (w *wrapper) GetAuthURL() string {
	return w.auth.AuthURL(state)
}

func (w *wrapper) GetTokenFromResponseCode(responseCode string) (*oauth2.Token, error) {
	return w.auth.Exchange(context.Background(), responseCode)
}

func (w *wrapper) LoginAndCreateClient(token *oauth2.Token) {
	client := spotify.New(w.auth.Client(context.Background(), token), spotify.WithRetry(true))
	w.client = client
}

func (w *wrapper) GetAllPlaylistsForUser(username string) ([]spotify.SimplePlaylist, error) {
	log.Infof("Getting list of playlists for user: %s", username)
	ctx := context.Background()
	playlists, err := w.client.GetPlaylistsForUser(ctx, username, spotify.Limit(50))
	if err != nil {
		return nil, err
	}

	var allPlaylists []spotify.SimplePlaylist
	allPlaylists = append(allPlaylists, playlists.Playlists...)

	log.Infof("User has %d total playlists", playlists.Total)
	for {
		err = w.client.NextPage(ctx, playlists)
		if err == spotify.ErrNoMorePages {
			break
		}

		allPlaylists = append(allPlaylists, playlists.Playlists...)
		log.Debugf("Retrieved %d playlists", len(allPlaylists))

		if err != nil {
			return nil, err
		}
	}

	log.Debugf("Retrieved %d total playlists", len(allPlaylists))
	return allPlaylists, nil
}

func (w *wrapper) GetAllPlaylistTracks(playlistID spotify.ID) ([]spotify.PlaylistTrack, error) {
	ctx := context.Background()
	tracks, err := w.client.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	var allTracks []spotify.PlaylistTrack
	allTracks = append(allTracks, tracks.Tracks...)

	log.Infof("Getting %d total tracks", tracks.Total)
	for {
		err = w.client.NextPage(ctx, tracks)
		if err == spotify.ErrNoMorePages {
			break
		}

		allTracks = append(allTracks, tracks.Tracks...)
		log.Debugf("Retrieved %d tracks", len(allTracks))

		if err != nil {
			return nil, err
		}
	}

	log.Debugf("Retrieved %d total tracks from playlist", len(allTracks))
	return allTracks, nil
}

func (w *wrapper) RemoveTracksFromPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) error {
	log.Debugf("Removing tracks %s from playlist %s", trackIDs, playlistID)
	ctx := context.Background()
	_, err := w.client.RemoveTracksFromPlaylist(ctx, playlistID, trackIDs...)
	return err
}

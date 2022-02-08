package adapter

import (
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyAuthWrapperInterface interface {
	New(opts ...spotifyauth.AuthenticatorOption) *spotifyauth.Authenticator
	WithRedirectURL(url string) spotifyauth.AuthenticatorOption
	WithScopes(scopes ...string) spotifyauth.AuthenticatorOption
}

type SpotifyWrapperInterface interface {
	GetAllPlaylistsForUser(username string) ([]spotify.SimplePlaylist, error)
	GetAllPlaylistTracks(playlistID spotify.ID) ([]spotify.PlaylistTrack, error)
	RemoveTracksFromPlaylist(playlistID spotify.ID, trackIDs ...spotify.ID) error
	GetAuthURL() string
	GetTokenFromResponseCode(responseCode string) (*oauth2.Token, error)
	LoginAndCreateClient(token *oauth2.Token)
	CreateAuthenticator(redirectURL string)
}

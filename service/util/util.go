package util

import (
	"strings"

	"github.com/reeves122/spotify-automation-go/adapter"
	"github.com/reeves122/spotify-automation-go/service"
	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify/v2"
)

type util struct {
	spotify        adapter.SpotifyWrapperInterface
	storage        service.StorageInterface
	dislikedPrefix string // ex: 'disliked_'
	queueSuffix    string // ex: ' Queue'
}

func NewUtil(spotify adapter.SpotifyWrapperInterface, storage service.StorageInterface, dislikedPrefix string, queueSuffix string) *util {
	return &util{
		spotify:        spotify,
		storage:        storage,
		dislikedPrefix: dislikedPrefix,
		queueSuffix:    queueSuffix,
	}
}

func (u *util) GetAllPlaylistsForUser(username string) ([]spotify.SimplePlaylist, error) {
	return u.spotify.GetAllPlaylistsForUser(username)
}

// UpdateLocalCache saves the contents of all playlists to a file
func (u *util) UpdateLocalCache(playlists []spotify.SimplePlaylist) error {
	log.Info("Updating local cache of playlists")

	for _, playlist := range playlists {

		log.Infof("Getting list of tracks for playlist: %s", playlist.Name)
		cached, err := u.storage.LoadTracksFile(playlist.Name)
		if err != nil {
			return err
		}

		// Compare the count of tracks to detect changes, for lack of a better option
		if len(cached) == int(playlist.Tracks.Total) {
			continue
		}

		log.Infof("Detected changes in playlist: %s", playlist.Name)
		tracks, err := u.spotify.GetAllPlaylistTracks(playlist.ID)
		if err != nil {
			return err
		}

		err = u.storage.SaveTracksFile(playlist.Name, tracks)
		if err != nil {
			return err
		}
		log.Infof("Done updating cache for playlist: %s", playlist.Name)
	}
	return nil
}

// LoadAllDislikedTracks loads tracks from all playlists matching the dislikedPrefix pattern
func (u *util) LoadAllDislikedTracks(playlists []spotify.SimplePlaylist) ([]spotify.PlaylistTrack, error) {
	log.Info("Building list of all disliked tracks")

	var allTracks []spotify.PlaylistTrack

	for _, playlist := range playlists {
		if strings.HasPrefix(playlist.Name, u.dislikedPrefix) {
			tracks, err := u.storage.LoadTracksFile(playlist.Name)
			if err != nil {
				return nil, err
			}
			allTracks = append(allTracks, tracks...)
		}
	}

	log.Infof("Loaded %d total disliked tracks", len(allTracks))
	return allTracks, nil
}

// ScanPlaylistsForDislikedTracks checks all playlists for any disliked tracks
func (u *util) ScanPlaylistsForDislikedTracks(playlists []spotify.SimplePlaylist, disliked []spotify.PlaylistTrack, username string) error {
	log.Infof("Scanning playlists for disliked tracks")
	dislikedHash := createTrackIdHash(disliked)

	for _, playlist := range playlists {
		if strings.HasPrefix(playlist.Name, u.dislikedPrefix) {
			continue
		}

		if playlist.Owner.ID != username {
			continue
		}

		err := u.scanPlaylistForDislikedTracks(playlist, dislikedHash)
		if err != nil {
			return err
		}
	}
	return nil
}

// scanPlaylistForDislikedTracks scans the playlist for any track whose ID
// is in the disliked map and removes it
func (u *util) scanPlaylistForDislikedTracks(playlist spotify.SimplePlaylist, disliked map[string]bool) error {
	log.Infof("Scanning playlist %s for disliked tracks", playlist.Name)

	tracks, err := u.storage.LoadTracksFile(playlist.Name)
	if err != nil {
		return nil
	}
	for _, track := range tracks {
		if _, present := disliked[track.Track.ID.String()]; present {
			log.WithFields(log.Fields{
				"name":   track.Track.Name,
				"artist": track.Track.Artists[0].Name,
				"album":  track.Track.Album.Name,
				"id":     track.Track.ID}).
				Warningf("Disliked track found")

			err := u.spotify.RemoveTracksFromPlaylist(playlist.ID, track.Track.ID)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

// ProcessQueuePlaylists checks all queue playlists
func (u *util) ProcessQueuePlaylists(playlists []spotify.SimplePlaylist, username string) error {
	for _, playlist := range playlists {
		if !strings.HasSuffix(playlist.Name, u.queueSuffix) {
			continue
		}

		if playlist.Owner.ID != username {
			continue
		}

		err := u.processQueuePlaylist(playlist)
		if err != nil {
			return err
		}
	}
	return nil
}

// Scan a "Queue" playlist (playlist of songs yet to be listened to and rated) for songs
// which have been added to the corresponding destination playlist. For example, the user may
// have "Favorites" and "Favorites Queue" playlists. The latter being songs the user has not
// heard and rated before. If the user likes a song, they add it to the "Favorites" list and this
// function will then remove it from the "Favorites Queue" playlist.
func (u *util) processQueuePlaylist(playlist spotify.SimplePlaylist) error {
	log.Infof("Processing queue playlist: %s", playlist.Name)

	destPlaylistTracks, err := u.storage.LoadTracksFile(strings.Replace(playlist.Name, u.queueSuffix, "", 1))
	if err != nil {
		return err
	}
	destPlaylistTracksHash := createTrackIdHash(destPlaylistTracks)

	playlistTracks, err := u.storage.LoadTracksFile(playlist.Name)
	if err != nil {
		return err
	}
	for _, track := range playlistTracks {
		if _, present := destPlaylistTracksHash[track.Track.ID.String()]; present {
			log.WithFields(log.Fields{
				"name":   track.Track.Name,
				"artist": track.Track.Artists[0].Name,
				"album":  track.Track.Album.Name,
				"id":     track.Track.ID}).
				Warningf("Queue track found in destination playlist")

			err := u.spotify.RemoveTracksFromPlaylist(playlist.ID, track.Track.ID)
			if err != nil {
				return nil
			}
		}
	}
	return nil
}

// FindPossibleDuplicateTracks finds tracks which are duplicated in a playlist
func (u *util) FindPossibleDuplicateTracks(playlists []spotify.SimplePlaylist) error {
	log.Warning("FindPossibleDuplicateTracks() not implemented")
	return nil
}

// createTrackIdHash creates a map of track ID for quick lookups
func createTrackIdHash(tracks []spotify.PlaylistTrack) map[string]bool {
	dislikedHash := map[string]bool{}
	for _, track := range tracks {
		dislikedHash[track.Track.ID.String()] = true
	}
	return dislikedHash
}

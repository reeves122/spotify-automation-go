# spotify-automation-go

## Overview

This project uses the Go Spotify wrapper (https://github.com/zmb3/spotify) to automate various
tasks for managing a personal Spotify library. 

The tasks include:

- Backing up all user playlist tracks to JSON file for safe keeping and exporting to other services
- Scanning all user playlists for "disliked" tracks and removing them (details below)
- Processing user "Queue" playlists (details below)
- Scanning a playlist for duplicate tracks (based on name and time length)


## How to Run
Main branch builds are automatically uploaded to Docker Hub:
https://hub.docker.com/repository/docker/reeves122/spotify-automation-go


Example of running with docker:
```
docker run --env-file spotify.env -v "/Users/example/spotify_cache:/spotify_cache" reeves122/spotify-automation-go:latest
```

Use the example environment file below.


## Docker Environment File Example

```
USER_NAME=reeves122
SPOTIFY_ID=...
SPOTIFY_SECRET=...
REDIRECT_URL=http://localhost:8888/callback
CACHE_DIR=/spotify_cache
TOKEN_FILE=auth_token.json
DISLIKED_PREFIX=disliked_
QUEUE_SUFFIX= Queue
```

## Disliked Tracks
Spotify has no concept of a "thumbs down" rating for songs in your library or 
playlist (like YouTube Music or Apple Music). Instead, a user may maintain a playlist of "disliked" 
songs and this program will use that playlist of "disliked" songs to
prune other playlists across your library. 

This program supports disliked playlists named with the `DISLIKED_PREFIX`. For example: `disliked_1`


## Queue Playlists
Scans a "Queue" playlist (playlist of songs yet to be listened to and rated) for songs
which have been added to the corresponding destination playlist. For example, the user may
have "Favorites" and "Favorites Queue" playlists. The latter being songs the user has not
heard and rated before. If the user likes a song, they add it to the "Favorites" list and this
program will then remove it from the "Favorites Queue" playlist.

This program supports Queue playlists named with the `QUEUE_SUFFIX`. For example: `Favorites Queue`

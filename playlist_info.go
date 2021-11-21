package main

import (
	"os"

	"github.com/zmb3/spotify"
)

var playlistID = os.Getenv("SPOTIFY_PLAYLIST_ID")

type playListInfo struct {
	ID             spotify.ID
	Name           string `bson:"name"`
	NumberOfTracks int    `bson:"numberOfTracks"`
	Tracks         spotify.PlaylistTrackPage
}

func createPlaylistInfoFromPlaylist(playlist *spotify.FullPlaylist) playListInfo {
	return playListInfo{
		ID:             playlist.ID,
		Name:           playlist.Name,
		NumberOfTracks: playlist.Tracks.Total,
		Tracks:         playlist.Tracks,
	}
}

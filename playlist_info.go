package main

import (
	"fmt"
	"os"

	"github.com/globalsign/mgo/bson"
	"github.com/zmb3/spotify/v2"
)

var playlistID = os.Getenv("SPOTIFY_PLAYLIST_ID")

type playListInfo struct {
	ID             spotify.ID `bson:"id"`
	Name           string     `bson:"name"`
	NumberOfTracks int        `bson:"numberOfTracks"`
}

func createPlaylistInfoFromPlaylist(playlist *spotify.FullPlaylist) playListInfo {
	return playListInfo{
		ID:             playlist.ID,
		Name:           playlist.Name,
		NumberOfTracks: playlist.Tracks.Total,
	}
}

func getPlaylistInfoByID(id spotify.ID) (playListInfo, error) {
	var pi playListInfo
	if err := getMongoCollection().Find(bson.M{"id": id}).One(&pi); err != nil {
		return playListInfo{}, fmt.Errorf("error getting playlist: %v", err)
	}

	return pi, nil
}

func updatePlaylistInfo(currentPi playListInfo, playlist *spotify.FullPlaylist) error {
	updatedPi := createPlaylistInfoFromPlaylist(playlist)

	if err := getMongoCollection().Update(currentPi, updatedPi); err != nil {
		return fmt.Errorf("error updating playlist info in DB: %v", err)
	}

	return nil
}

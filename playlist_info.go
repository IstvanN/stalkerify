package main

import (
	"fmt"
	"log"
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

func comparePlaylistWithPlaylistInfoInDB(playlist *spotify.FullPlaylist, pi playListInfo) error {
	if playlist.Tracks.Total > pi.NumberOfTracks {
		log.Println("new song found, email has been sent to you!")
		if err := sendMail(playlist.Name); err != nil {
			return fmt.Errorf("error sending mail: ", err)
		}

		if err := updatePlaylistInfo(pi, playlist); err != nil {
			return fmt.Errorf("error updating playlistinfo in DB: ", err)
		}
		return nil
	}
	log.Println("no new song found!")
	return nil
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/zmb3/spotify/v2"
)

var playlistID = os.Getenv("SPOTIFY_PLAYLIST_ID")

type playListInfo struct {
	ID             spotify.ID `bson:"id"`
	Name           string     `bson:"name"`
	NumberOfTracks int        `bson:"numberOfTracks"`
}

type newSongData struct {
	addedBy, artist, title, addedAt string
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

		newSongs := getNewSongDatas(playlist, pi)
		if err := sendMail(playlist.Name, newSongs); err != nil {
			return fmt.Errorf("error sending mail: %v", err)
		}

		if err := updatePlaylistInfo(pi, playlist); err != nil {
			return fmt.Errorf("error updating playlistinfo in DB: %v", err)
		}
		return nil
	}
	log.Println("no new song found!")
	return nil
}

func getNewSongDatas(playlist *spotify.FullPlaylist, pi playListInfo) []newSongData {
	var newSongs []newSongData
	for i, track := range playlist.Tracks.Tracks {
		if i > pi.NumberOfTracks-1 {
			addedAt, _ := time.Parse(spotify.TimestampLayout, track.AddedAt)
			budapest, _ := time.LoadLocation("Europe/Budapest")
			addedAtString := addedAt.In(budapest).Format("2006.01.06 15:03")
			nsd := newSongData{
				addedBy: track.AddedBy.ID,
				artist:  track.Track.Artists[0].Name,
				title:   track.Track.Name,
				addedAt: addedAtString,
			}
			newSongs = append(newSongs, nsd)
		}
	}
	return newSongs
}

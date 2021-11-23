package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/zmb3/spotify/v2"
)

var playlistIDs = os.Getenv("SPOTIFY_PLAYLIST_IDS")

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

func comparePlaylistWithPlaylistInfoInDB(playlist *spotify.FullPlaylist, pi playListInfo, client *spotify.Client) error {
	log.Println("checking playlist: ", playlist.Name)
	if playlist.Tracks.Total < pi.NumberOfTracks {
		log.Println("playlist's song number < number of tracks in DB: maybe someone deleted track(s)")
		log.Println("setting it up to date in DB...")
		if err := updatePlaylistInfo(pi, playlist); err != nil {
			return fmt.Errorf("error updating playlistinfo in DB: %v", err)
		}
		log.Printf("number of tracks set to %v in DB\n", pi.NumberOfTracks)
		return nil
	}

	if playlist.Tracks.Total > pi.NumberOfTracks {
		log.Println("new song found, email has been sent to you!")

		newSongs, err := getNewSongDatas(playlist, pi, client)
		if err != nil {
			return fmt.Errorf("error getting info on new songs: %v", err)
		}
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

func getNewSongDatas(playlist *spotify.FullPlaylist, pi playListInfo, client *spotify.Client) ([]newSongData, error) {
	var newSongs []newSongData
	for i, track := range playlist.Tracks.Tracks {
		if i > pi.NumberOfTracks-1 {
			addedAt, err := transformAddedAt(track.AddedAt)
			if err != nil {
				return nil, err
			}

			user, err := client.GetUsersPublicProfile(context.Background(), spotify.ID(track.AddedBy.ID))
			if err != nil {
				return nil, err
			}

			nsd := newSongData{
				addedBy: user.DisplayName,
				artist:  track.Track.Artists[0].Name,
				title:   track.Track.Name,
				addedAt: addedAt,
			}
			newSongs = append(newSongs, nsd)
		}
	}
	return newSongs, nil
}

func transformAddedAt(old string) (string, error) {
	addedAt, err := time.Parse(spotify.TimestampLayout, old)
	if err != nil {
		return "", err
	}

	budapest, err := time.LoadLocation("Europe/Budapest")
	if err != nil {
		return "", err
	}

	return addedAt.In(budapest).Format("2006.01.06 15:03"), nil
}

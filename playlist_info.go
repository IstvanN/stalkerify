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

type playlistInfo struct {
	ID    spotify.ID              `bson:"id"`
	Name  string                  `bson:"name"`
	Songs []spotify.PlaylistTrack `bson:"songs"`
}

type newSongData struct {
	addedBy, artist, title, addedAt string
}

func createPlaylistInfoFromPlaylist(playlist *spotify.FullPlaylist) playlistInfo {
	return playlistInfo{
		ID:    playlist.ID,
		Name:  playlist.Name,
		Songs: playlist.Tracks.Tracks,
	}
}

func getPlaylistInfoByID(id spotify.ID) (playlistInfo, error) {
	var pi playlistInfo
	if err := getMongoCollection().Find(bson.M{"id": id}).One(&pi); err != nil {
		return playlistInfo{}, fmt.Errorf("error getting playlist: %v", err)
	}

	return pi, nil
}

func updatePlaylistInfo(currentPi playlistInfo, playlist *spotify.FullPlaylist) error {
	updatedPi := createPlaylistInfoFromPlaylist(playlist)

	if err := getMongoCollection().Update(currentPi, updatedPi); err != nil {
		return fmt.Errorf("error updating playlist info in DB: %v", err)
	}

	return nil
}

func comparePlaylistWithPlaylistInfoInDB(playlist *spotify.FullPlaylist, pi playlistInfo, client *spotify.Client) error {
	log.Println("checking playlist:", playlist.Name)
	if err := appendAdditionalSongsToPlaylist(playlist, client); err != nil {
		return err
	}

	nsd, err := getNewSongDatas(playlist, pi, client)
	if err != nil {
		return fmt.Errorf("error getting info on new songs: %v", err)
	}

	if len(nsd) == 0 {
		log.Println("no new song found!")
		log.Println("updating DB because of potential deletes...")
		if err := updatePlaylistInfo(pi, playlist); err != nil {
			return fmt.Errorf("error updating playlistinfo in DB: %v", err)
		}
		return nil
	}

	if err := sendMail(playlist.Name, nsd); err != nil {
		return fmt.Errorf("error sending mail: %v", err)
	}
	log.Println("new song found, email has been sent to you!")

	if err := updatePlaylistInfo(pi, playlist); err != nil {
		return fmt.Errorf("error updating playlistinfo in DB: %v", err)
	}

	return nil
}

func isSongInDB(song *spotify.PlaylistTrack, pi playlistInfo) bool {
	for _, track := range pi.Songs {
		if song.Track.ID == track.Track.ID {
			return true
		}
	}

	return false
}

func appendAdditionalSongsToPlaylist(playlist *spotify.FullPlaylist, client *spotify.Client) error {
	request := spotify.Offset(100)
	songs, err := client.GetPlaylistTracks(context.Background(), playlist.ID, request)
	if err != nil {
		return fmt.Errorf("error returning songs from spotify API: %v", err)
	}

	playlist.Tracks.Tracks = append(playlist.Tracks.Tracks, songs.Tracks...)

	return nil
}

func getNewSongDatas(playlist *spotify.FullPlaylist, pi playlistInfo, client *spotify.Client) ([]newSongData, error) {
	var newSongs []newSongData
	for _, track := range playlist.Tracks.Tracks {
		if !isSongInDB(&track, pi) {
			addedAt, err := transformAddedAt(track.AddedAt)
			if err != nil {
				return nil, err
			}

			spotifyUser, err := client.GetUsersPublicProfile(context.Background(), spotify.ID(track.AddedBy.ID))
			if err != nil {
				return nil, err
			}

			artistName := createArtistNameOfArtists(track.Track.Artists)

			ns := newSongData{
				addedBy: spotifyUser.DisplayName,
				artist:  artistName,
				title:   track.Track.Name,
				addedAt: addedAt,
			}
			newSongs = append(newSongs, ns)
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

	return addedAt.In(budapest).Format("Monday, 2006 Jan 2 15:04"), nil
}

func createArtistNameOfArtists(artists []spotify.SimpleArtist) string {
	var artistName string
	for i, artist := range artists {
		if i == len(artists)-1 {
			artistName += artist.Name
		} else {
			artistName += artist.Name + ", "
		}
	}

	return artistName
}

package main

import (
	"context"
	"log"

	"github.com/zmb3/spotify/v2"
)

func main() {
	db, err := startupMongo()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	log.Println("successfully connected to MongoDB!")

	client, err := getSpotifyClient()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("successfully created Spotify client!")

	playlistIDsSlice, err := convertCSVtoSliceOfStrings(playlistIDs)
	if err != nil {
		log.Fatal(err)
	}

	for _, plid := range playlistIDsSlice {
		playlist, err := client.GetPlaylist(context.Background(), spotify.ID(plid))
		if err != nil {
			log.Fatalln(err)
		}

		currentPi, err := getPlaylistInfoByID(playlist.ID)
		if err != nil {
			log.Fatalln(err)
		}

		if err := comparePlaylistWithPlaylistInfoInDB(playlist, currentPi, client); err != nil {
			log.Fatalln(err)
		}
	}
}

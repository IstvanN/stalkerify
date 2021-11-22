package main

import (
	"context"
	"fmt"
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

	playlist, err := client.GetPlaylist(context.Background(), spotify.ID(playlistID))
	if err != nil {
		log.Fatalln(err)
	}

	currentPi, err := getPlaylistInfoByID(playlist.ID)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(playlist.Tracks)

	if err := comparePlaylistWithPlaylistInfoInDB(playlist, currentPi); err != nil {
		fmt.Println(err)
	}
}

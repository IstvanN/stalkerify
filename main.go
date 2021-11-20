package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
	spotifyauth "github.com/zmb3/spotify/auth"
	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	client, err := getSpotifyClient()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(client)
}

func getSpotifyClient() (*spotify.Client, error) {
	config := clientcredentials.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error retrieving token for spotify client: %v", err)
	}

	httpClient := spotifyauth.New().Client(context.Background(), token)
	return spotify.New(httpClient), nil
}

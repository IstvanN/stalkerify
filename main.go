package main

import (
	"fmt"
	"log"
)

func main() {
	client, err := getSpotifyClient()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(client)
}

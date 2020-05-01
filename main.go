package main

import (
	log "github.com/sirupsen/logrus"
	"sync-sandbox/data/spotify"
)

func main() {
	spotifyAuth, err := spotify.NewSpotifyDAO()
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"spotify_auth": spotifyAuth,
	}).Info("fetched auth token")
}

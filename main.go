package main

import (
	data "github.com/cshealy/sync-sandbox/data/spotify"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
)

// 12 factor config for establishing connections and other various attributes with env vars
type Config struct {
	ServerPort int `required:"true" default:"50051" split_words:"true"` // SERVER_PORT
}

func main() {

	// process environment variables for config
	var config Config
	err := envconfig.Process("sync-sandbox", &config)

	// check for any errors while parsing environment variables
	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithFields(log.Fields{
		"config": config,
	}).Info("initialized config")

	// create dao for spotify
	spotifyAuth, err := data.NewDAO()

	// check for any errors while fetching the spotify bearer token
	if err != nil {
		log.Fatal(err.Error())
	}

	// confirm spotify dao is generated
	log.WithFields(log.Fields{
		"spotify_auth": spotifyAuth,
	}).Info("fetched auth token")
}

package main

import (
	"context"
	"fmt"
	"net"

	"github.com/cshealy/sync-sandbox/rpc"

	data "github.com/cshealy/sync-sandbox/data/spotify"
	pb "github.com/cshealy/sync-sandbox/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"net/http"

	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// 12 factor config for establishing connections and other various attributes with env vars
type Config struct {
	GatewayPort int `required:"true" default:"8080" split_words:"true"`  // GATEWAY_PORT
	ServerPort  int `required:"true" default:"50051" split_words:"true"` // SERVER_PORT

	SpotifyPassword   string `required:"true" split_words:"true"` // SPOTIFY_PASSWORD
	SpotifyPlaylistId string `required:"true" split_words:"true"` // SPOTIFY_PLAYLIST_ID
	SpotifyUsername   string `required:"true" split_words:"true"` // SPOTIFY_USERNAME
}

func init() {
	// TODO: set log level from config
}

func main() {

	// process environment variables for config
	var config Config
	err := envconfig.Process("", &config)

	// check for any errors while parsing environment variables
	if err != nil {
		log.Fatal(err.Error())
	}

	// create dao for spotify
	spotifyDAO, err := data.NewDAO(config.SpotifyUsername, config.SpotifyPassword, config.SpotifyPlaylistId)

	// check for any errors while fetching the spotify bearer token
	if err != nil {
		log.Fatal(err.Error())
	}

	log.WithFields(log.Fields{
		"grpc_server": fmt.Sprintf("http://localhost:%d", config.ServerPort),
	}).Info("starting gRPC server")

	// start gRPC server
	go runServer(config, spotifyDAO)

	log.WithFields(log.Fields{
		"grpc_gateway": fmt.Sprintf("http://localhost:%d", config.GatewayPort),
	}).Info("starting gRPC gateway")

	// start our gRPC gateway
	if err = runGateway(config); err != nil {
		log.Fatal(err)
	}
}

// RunServer will start a new server and begin listening using the port provided by our config
func runServer(config Config, spotifyDAO *data.SpotifyDAO) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
	if err != nil {
		return err
	}

	// create a new server
	grpcServer := grpc.NewServer()

	// register TestService
	pb.RegisterTestsServer(grpcServer, &rpc.TestService{
		SpotifyDAO: spotifyDAO,
	})

	// begin listening
	grpcServer.Serve(listen)

	return nil
}

// RunGateway starts to translate RESTful HTTP API into gRPC
func runGateway(config Config) error {

	// create our context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create the multiplexer
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterTestsHandlerFromEndpoint(ctx, mux, fmt.Sprintf(":%d", config.ServerPort), opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", config.GatewayPort), mux)
}

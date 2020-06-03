package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

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
	LogLevel string `required:"true" default:"info" split_words:"true"` // LOG_LEVEL

	GatewayPort int `required:"true" default:"8080" split_words:"true"`  // GATEWAY_PORT
	ServerPort  int `required:"true" default:"50051" split_words:"true"` // SERVER_PORT

	SpotifyPassword   string `required:"true" split_words:"true"` // SPOTIFY_PASSWORD
	SpotifyPlaylistId string `required:"true" split_words:"true"` // SPOTIFY_PLAYLIST_ID
	SpotifyUsername   string `required:"true" split_words:"true"` // SPOTIFY_USERNAME
}

var config Config

func init() {

	// process environment variables for config
	err := envconfig.Process("", &config)

	// check for any errors while parsing environment variables
	if err != nil {
		log.Fatal(err.Error())
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the severity set in env var
	// https://github.com/sirupsen/logrus/blob/39a5ad12948d094ddd5d5a6a4a4281f453d77562/logrus.go#L25
	logLevel, err := log.ParseLevel(strings.ToLower(config.LogLevel))

	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(logLevel)

	log.WithFields(log.Fields{
		"log_level": logLevel,
	}).Info("set log level")
}

func main() {

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
	// TODO: take this mux returned and serve the swagger docs
	_, err = runGateway(config)

	if err != nil {
		log.Fatal(err)
	}

	// TODO: server swagger docs
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
func runGateway(config Config) (*runtime.ServeMux, error) {

	// create our context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create the multiplexer
	rmux := runtime.NewServeMux()

	// Serve the swagger-ui and swagger file
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	mux.HandleFunc("/sync.swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "proto/sync.swagger.json")
	})
	fs := http.FileServer(http.Dir("proto"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := pb.RegisterTestsHandlerFromEndpoint(ctx, rmux, fmt.Sprintf(":%d", config.ServerPort), opts)
	if err != nil {
		return nil, err
	}

	return rmux, http.ListenAndServe(fmt.Sprintf(":%d", config.GatewayPort), rmux)
}

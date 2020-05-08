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
}

func init() {
	// TODO: set log level from config
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

	// start gRPC server
	go RunServer(config)

	// start our gRPC gateway
	if err = RunGateway(config); err != nil {
		log.Fatal(err)
	}
}

// RunServer will start a new server and begin listening using the port provided by our config
func RunServer(config Config) error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", config.ServerPort))
	if err != nil {
		return err
	}

	// TODO: add HTTP mux

	// create a new server
	grpcServer := grpc.NewServer()

	// register TestService
	pb.RegisterTestsServer(grpcServer, &rpc.TestService{})

	// begin listening
	grpcServer.Serve(listen)

	return nil
}

// RunGateway starts to translate RESTful HTTP API into gRPC
func RunGateway(config Config) error {

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

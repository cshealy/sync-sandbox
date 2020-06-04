package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	pb "github.com/cshealy/sync-sandbox/proto"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// 12 factor config for establishing connections and other various attributes with env vars
type Config struct {
	LogLevel   string `required:"true" default:"info" split_words:"true"`  // LOG_LEVEL
	ServerPort string `required:"true" default:"50051" split_words:"true"` // SERVER_PORT
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

	// dial our api microservice internally
	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(fmt.Sprintf("sync-sandbox-api:%s", config.ServerPort), opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// create context with a reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create a client to communicate via gRPC
	client := NewClient(ctx, pb.NewTestsClient(conn))

	// call unary RPC test
	unaryResponse := client.GetUnaryTest("This is a test")

	log.WithFields(log.Fields{
		"GetTest": unaryResponse,
	}).Info("received test")

	// call GetSpotifyPlaylistStream via gRPC
	serverStreamResponse := client.GetServerStreamPlaylist()

	log.WithFields(log.Fields{
		"streamed_playlist": serverStreamResponse,
	}).Info("received GetSpotifyPlaylistStream")

	// client-side streaming
	clientStreamResponse := client.GetClientStream()

	// log the results
	log.WithFields(log.Fields{
		"tests": clientStreamResponse,
	}).Info("received GetClientStream")

	// bidirectional streaming
	client.GetBidirectionalStream()
}

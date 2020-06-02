package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	pb "github.com/cshealy/sync-sandbox/proto"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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

	// create a client to communicate via gRPC
	client := pb.NewTestsClient(conn)

	// create context with a reasonable timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// call GetTest via gRPC
	echo, err := client.GetTest(ctx, &pb.Test{
		Name: "this is a test",
	})

	if err != nil {
		log.Fatalf("failed to call GetTest: %v", err)
	}

	log.WithFields(log.Fields{
		"GetTest": echo,
	}).Info("received test")

	// call GetSpotifyPlaylistStream via gRPC
	spotifyPlaylistStream, err := client.GetSpotifyPlaylistStream(ctx, &emptypb.Empty{})

	if err != nil {
		log.Fatalf("failed to call GetTest: %v", err)
	}

	var tracks []*pb.Track

	// read from the stream
	for {
		track, err := spotifyPlaylistStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to call GetSpotifyPlaylistStream: %v", err)
		}
		tracks = append(tracks, track)
	}

	// compose the playlist
	streamedPlaylist := &pb.SpotifyPlaylist{
		Tracks: tracks,
	}

	log.WithFields(log.Fields{
		"streamed_playlist": streamedPlaylist,
	}).Info("received GetSpotifyPlaylistStream")

	// client-side streaming
	var tests []*pb.Test
	for i := 0; i < 10; i++ {
		tests = append(tests, &pb.Test{
			Name: fmt.Sprintf("Test - %d", i),
		})
	}
	stream, err := client.GetClientStream(context.Background())
	if err != nil {
		log.Fatalf("failed to call GetClientStream: %v", err)
	}

	// send off those tests
	for _, test := range tests {
		if err := stream.Send(test); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("failed to send %v to GetClientStream - %v", test, err)
		}
	}

	// close the client stream
	testEcho, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("failed to close stream for GetClientStream")
	}

	// log the results
	log.WithFields(log.Fields{
		"tests": testEcho,
	}).Info("received GetClientStream")

	// bidirectional streaming
	bidirectionalStream, err := client.GetBidirectionalStream(context.Background())

	// wait channel
	waitc := make(chan struct{})

	// go routine to receive response from server stream
	go func() {
		for {
			test, err := bidirectionalStream.Recv()

			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("unable to get test: %s", err)
			}
			log.WithFields(log.Fields{
				"test": test,
			}).Info("received test from server")
		}
	}()

	// stream tests to server
	for _, test := range tests {
		if err := bidirectionalStream.Send(test); err != nil {
			log.Fatalf("failed to send test: %s", err)
		}
	}

	// close client stream
	bidirectionalStream.CloseSend()
	<-waitc
}

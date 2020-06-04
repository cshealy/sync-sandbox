package main

import (
	"context"
	"fmt"
	"io"

	pb "github.com/cshealy/sync-sandbox/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Client is the configuration for our api microservice client
type Client struct {
	testClient pb.TestsClient  // client used to communicate with api microservice
	ctx        context.Context // context used for communicating with client
}

// NewClient generates a new client for communicating with our API
func NewClient(ctx context.Context, testClient pb.TestsClient) *Client {
	return &Client{
		testClient: testClient,
		ctx:        ctx,
	}
}

// GetUnaryTest calls GetTest using unary RPC
func (client *Client) GetUnaryTest(name string) *pb.Test {

	// call GetTest via gRPC
	response, err := client.testClient.GetTest(client.ctx, &pb.Test{
		Name: name,
	})

	// report error if one occurs
	if err != nil {
		log.Fatalf("failed to call GetTest: %v", err)
	}

	// return the response from GetTest
	return response
}

// GetServerStreamPlaylist calls the test client to receive the spotify playlist through a server stream
func (client *Client) GetServerStreamPlaylist() *pb.SpotifyPlaylist {

	spotifyPlaylistStream, err := client.testClient.GetSpotifyPlaylistStream(client.ctx, &emptypb.Empty{})

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
	return &pb.SpotifyPlaylist{
		Tracks: tracks,
	}
}

// GetClientStream calls the client API's GetClientStream and returns the response
func (client *Client) GetClientStream() *pb.MultiTest {

	// generate tests
	tests := generateTests()

	stream, err := client.testClient.GetClientStream(client.ctx)

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

	return testEcho
}

// GetBidirectionalStream calls the client API's GetBidirectionalStream and logs the response
func (client *Client) GetBidirectionalStream() {

	// generate tests
	tests := generateTests()

	bidirectionalStream, err := client.testClient.GetBidirectionalStream(client.ctx)

	if err != nil {
		log.Fatalf("failed to call GetBidirectionalStream")
	}

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

// generateTests creates some tests to stream to our client API
func generateTests() []*pb.Test {

	// generate tests
	var tests []*pb.Test
	for i := 0; i < 10; i++ {
		tests = append(tests, &pb.Test{
			Name: fmt.Sprintf("Test - %d", i),
		})
	}

	return tests
}

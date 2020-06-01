package rpc

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc/codes"

	data "github.com/cshealy/sync-sandbox/data/spotify"
	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/cshealy/sync-sandbox/proto"
	"google.golang.org/grpc/status"
)

type TestService struct {
	*data.SpotifyDAO
}

// GetSpotifyPlaylist retrieves our spotify playlist and it's metadata, then returns it to the client
func (svc *TestService) GetSpotifyPlaylist(_ context.Context, _ *empty.Empty) (*pb.SpotifyPlaylist, error) {

	// fetch the spotify playlist from spotify
	spotifyPlaylist, err := svc.SpotifyDAO.GetPlaylist()

	// if any error occurs during our spotify processing, send it back to the client
	if err != nil {
		// TODO: return improved errors
		return nil, status.Error(codes.Unknown,
			fmt.Sprintf("unable to retrieve spotify playlist -- %s", svc.SpotifyDAO.SpotifyPlaylist))
	}

	return spotifyPlaylist, nil
}

// GetTest is an echo endpoint until I put more interesting logic into here
func (svc *TestService) GetTest(ctx context.Context, test *pb.Test) (*pb.Test, error) {
	// TODO: fill this out with logic
	return &pb.Test{
		Name: test.GetName(),
	}, nil
}

// GetSpotifyPlaylistStream gets a spotify playlist and streams back the tracks to our client
func (svc *TestService) GetSpotifyPlaylistStream(_ *empty.Empty, stream pb.Tests_GetSpotifyPlaylistStreamServer) error {

	// fetch the spotify playlist from spotify
	spotifyPlaylist, err := svc.SpotifyDAO.GetPlaylist()

	// if any error occurs during our spotify processing, send it back to the client
	if err != nil {
		// TODO: return improved errors
		return status.Error(codes.Unknown,
			fmt.Sprintf("unable to retrieve spotify playlist -- %s", svc.SpotifyDAO.SpotifyPlaylist))
	}

	// iterate over each track and stream it back to the client
	for _, track := range spotifyPlaylist.GetTracks() {
		// stream the track back to our client
		if err := stream.Send(track); err != nil {
			return err
		}
	}
	return nil
}

// GetClientStream acts as GetTest does, but as a client stream
func (svc *TestService) GetClientStream(stream pb.Tests_GetClientStreamServer) error {

	// create a slice of client streamed tests
	var multiTest []string

	// loop until the client is finished streaming
	for {

		// get the streamed test
		test, err := stream.Recv()

		// if the client is done streaming, return our multi tests
		if err == io.EOF {
			return stream.SendAndClose(&pb.MultiTest{
				Name: multiTest,
			})
		}

		// if we hit an error while reading from our client stream, return the error
		if err != nil {
			return err
		}

		// append a test to our response
		multiTest = append(multiTest, test.Name)
	}
}

// GetBidirectionalStream will act the same as GetTest does, but with client and server streaming
func (svc *TestService) GetBidirectionalStream(stream pb.Tests_GetBidirectionalStreamServer) error {
	// TODO: fill out
	return nil
}

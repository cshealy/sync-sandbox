package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	data "github.com/cshealy/sync-sandbox/data/spotify"
	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/cshealy/sync-sandbox/proto"
	log "github.com/sirupsen/logrus"
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
	log.Infof("%+v", ctx)
	log.Infof("%+v", test)
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

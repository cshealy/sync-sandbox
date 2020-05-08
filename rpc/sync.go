package rpc

import (
	"context"

	"github.com/cshealy/sync-sandbox/data"
	pb "github.com/cshealy/sync-sandbox/proto"
	log "github.com/sirupsen/logrus"
)

type TestService struct {
	data.DAO
}

// GetTest is an echo endpoint until I put more interesting logic into here
func (svc TestService) GetTest(ctx context.Context, test *pb.Test) (*pb.Test, error) {
	log.Infof("%+v", ctx)
	log.Infof("%+v", test)
	// TODO: fill this out with logic
	return &pb.Test{
		Name: test.GetName(),
	}, nil
}

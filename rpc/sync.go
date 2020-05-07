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

func (svc TestService) GetTest(ctx context.Context, test *pb.Test) (*pb.Test, error) {
	// TODO: fill this out with logic
	log.Info("made it to get test")
	return nil, nil
}
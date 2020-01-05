package grpc

import (
	"context"
	"errors"
	"github.com/reinkrul/go-keyvalue-db/api"
	"github.com/reinkrul/go-keyvalue-db/server/spi"
	"google.golang.org/grpc"
	"log"
	"net"
)

type impl struct {
	dataStore *spi.DataStore
	server    *grpc.Server
}

func Start(addr string, store spi.DataStore) (spi.API, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	result := impl{dataStore: &store, server: s}
	api.RegisterKeyValueDatabaseServer(s, &result)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return &result, nil
}

func (s impl) Set(ctx context.Context, in *api.SetValueRequest) (*api.SetValueResponse, error) {
	if in.Bucket == "" {
		return nil, errors.New("missing 'bucket'")
	}
	if in.Key == "" {
		return nil, errors.New("missing 'key'")
	}
	log.Printf("Writing %s/%s", in.Bucket, in.Key)
	err := (*s.dataStore).Set(in.Bucket, in.Key, in.Value)
	return &api.SetValueResponse{}, err
}

func (s impl) Get(ctx context.Context, in *api.GetValueRequest) (*api.GetValueResponse, error) {
	if in.Bucket == "" {
		return nil, errors.New("missing 'bucket'")
	}
	if in.Key == "" {
		return nil, errors.New("missing 'key'")
	}
	log.Printf("Reading %s/%s", in.Bucket, in.Key)
	value, err := (*s.dataStore).Get(in.Bucket, in.Key)
	return &api.GetValueResponse{Value: value}, err
}

func (s *impl) Close() error {
	log.Println("Stopping gRPC server")
	s.server.Stop()
	return nil
}

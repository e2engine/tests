package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	okv1 "github.com/e2engine/tests/services/ok-grpc/okv1"
)

const address = "127.0.0.1:9000"

type server struct {
	okv1.UnimplementedOKServiceServer
}

func (s *server) OK(
	context.Context,
	*okv1.OKRequest,
) (*okv1.OKResponse, error) {
	return &okv1.OKResponse{}, nil
}

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	okv1.RegisterOKServiceServer(
		grpcServer,
		&server{},
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

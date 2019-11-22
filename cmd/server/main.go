package main

import (
	"log"
	"net"

	"os"

	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/gcputil/env"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	logger   = log.New(os.Stdout, "", 0)
	grpcPort = env.MustGetEnvVar("PORT", "8080")
)

type pingService struct{}

func (s *pingService) Ping(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	if req == nil {
		return nil, errors.New("nil request")
	}

	if req.GetId() == "" {
		return nil, errors.New("nil Id")
	}

	if req.Target == nil {
		return nil, errors.New("nil Target")
	}

	logger.Printf("request: %+v", req)

	return &pb.Response{
		Request: req,
	}, nil

}

func startGRPCServer(hostPort string) error {
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen on %s: %v", hostPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, &pingService{})
	return grpcServer.Serve(listener)
}

func main() {
	grpcHostPort := net.JoinHostPort("0.0.0.0", grpcPort)
	go func() {
		err := startGRPCServer(grpcHostPort)
		if err != nil {
			logger.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()
	logger.Println("Server started...")
	select {}
}

package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"time"

	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// PingClient pings rundom endpont in collection of targets and returns it resp
func PingClient(req *pb.Request) (resp *pb.Response, err error) {
	if req == nil {
		return nil, errors.New("nil request")
	}
	targetNum := len(req.GetTargets())
	if targetNum == 0 {
		return nil, errors.New("zero targets")
	}
	rand.Seed(time.Now().UnixNano())
	target := req.GetTargets()[rand.Intn(targetNum)]
	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	return client.Ping(ctx, req)
}

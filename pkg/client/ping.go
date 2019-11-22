package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// PingClient pings rundom endpont in collection of targets and returns it resp
func PingClient(target *pb.Target) (resp *pb.Response, err error) {
	if target == nil {
		return nil, errors.New("nil target")
	}
	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", target.Uri, target.Port)
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", uri, err)
	}
	defer conn.Close()
	client := pb.NewPingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	req := &pb.Request{
		Id:     uuid.New().String(),
		Sent:   ptypes.TimestampNow(),
		Target: target,
	}

	return client.Ping(ctx, req)
}

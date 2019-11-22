package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"time"

	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	// EchoTimeout is the max number of seconds echo service will wait for response
	EchoTimeout = 300 * time.Second
)

// Ping pings request target
func Ping(req *pb.RequestMessage) (resp *pb.ResponseMessage, err error) {

	// validate
	if req == nil {
		return nil, errors.New("nil request")
	}
	if req.GetTarget() == nil {
		return nil, errors.New("nil target")
	}
	if req.GetSource() == nil {
		return nil, errors.New("nil source")
	}

	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", req.GetTarget().GetUri(), req.GetTarget().GetPort())
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %v", uri, err)
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), EchoTimeout)
	defer cancel()

	return client.Echo(ctx, req)
}

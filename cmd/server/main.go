package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/gcputil/env"
	"github.com/mchmarny/gcputil/metric"

	"crypto/tls"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger   = log.New(os.Stdout, "", 0)
	grpcPort = env.MustGetEnvVar("PORT", "8080")
	dbPath   = env.MustGetEnvVar("DB_PATH", "")
)

type echoService struct{}

func (s *echoService) Broadcast(ctx context.Context, in *pb.BroadcastMessage) (*pb.BroadcastResults, error) {

	if in == nil {
		return nil, errors.New("nil BroadcastMessage")
	}

	out := &pb.BroadcastResults{
		Results: make([]*pb.BroadcastResult, len(in.GetTargets())),
	}

	for i, t := range in.GetTargets() {
		out.Results[i] = execEcho(ctx, in.GetSelf(), t)
	}

	return out, nil

}

func (s *echoService) Echo(ctx context.Context, req *pb.EchoMessage) (resp *pb.EchoMessage, err error) {
	if req == nil {
		return req, errors.New("nil request")
	}
	logger.Printf("request: %+v", req)
	return req, nil
}

func execEcho(ctx context.Context, self *pb.Node, target *pb.Node) *pb.BroadcastResult {

	result := &pb.BroadcastResult{
		Source:   self,
		Target:   target,
		Duration: 0,
	}

	// setup connection
	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", target.GetUri(), target.GetPort())
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		result.Error = fmt.Sprintf("failed to dial %s: %v", uri, err)
		return result
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)

	// client message
	msgIn := &pb.EchoMessage{
		From: self.GetRegion(),
		Sent: ptypes.TimestampNow(),
	}

	// call remote service
	logger.Printf("submitting echo: %+v", msgIn)
	started := time.Now()
	msgOut, err := client.Echo(ctx, msgIn)
	if err != nil {
		result.Error = fmt.Sprintf("error while executing echo %s: %v", uri, err)
		return result
	}
	finished := time.Now()
	result.Duration = finished.Sub(started).Milliseconds()

	// make sure echo returned the same message
	if msgOut.GetFrom() != msgIn.GetFrom() {
		result.Error = fmt.Sprintf("unexpected echo from %s (want %s, got %s)",
			uri, msgIn.GetFrom(), msgOut.GetFrom())
		return result
	}

	// debug into server side
	logger.Printf("echo-ping from: %s to: %s (duration: %v)\n ",
		msgIn.GetFrom(), msgOut.GetFrom(), result.Duration)

	// save to db
	if err = save(ctx, dbPath, uuid.New().String(), msgIn.GetFrom(), msgOut.GetFrom(),
		started, finished, result.Duration); err != nil {
		result.Error = fmt.Sprintf("error while saving request %s: %v", uri, err)
		return result
	}

	// metrics
	labels := map[string]string{
		"source": msgIn.GetFrom(),
		"target": msgOut.GetFrom(),
	}
	if err = metric.MetricClient(ctx).Publish(ctx, "echo-duration", result.Duration, labels); err != nil {
		result.Error = fmt.Sprintf("error while saving echo-duration metric %s: %v", uri, err)
		return result
	}

	return result

}

func startGRPCServer(hostPort string) error {
	listener, err := net.Listen("tcp", hostPort)
	if err != nil {
		return errors.Wrapf(err, "Failed to listen on %s: %v", hostPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterEchoServiceServer(grpcServer, &echoService{})
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

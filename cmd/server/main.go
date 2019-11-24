package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"os"

	ptypes "github.com/golang/protobuf/ptypes"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/gcputil/env"
	"github.com/mchmarny/gcputil/metric"

	"crypto/tls"

	"github.com/google/uuid"
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

func (s *echoService) Broadcast(ctx context.Context, in *pb.BroadcastMessage) (*pb.BroadcastResult, error) {

	if in == nil {
		return nil, errors.New("nil BroadcastMessage")
	}

	for _, t := range in.GetTargets() {
		if err := execEcho(ctx, in.GetSelf(), t); err != nil {
			return nil, fmt.Errorf("error on echo: %v", err)
		}
	}
	
	return &pb.BroadcastResult{
		Count: int32(len(in.GetTargets())),
	}, nil

}

func (s *echoService) Echo(ctx context.Context, req *pb.EchoMessage) (resp *pb.EchoMessage, err error) {
	if req == nil {
		return req, errors.New("nil request")
	}
	logger.Printf("request: %+v", req)
	return req, nil
}

func execEcho(ctx context.Context, self *pb.Node, target *pb.Node) error {

	msg := &pb.EchoMessage{
		Id:     uuid.New().String(),
		Sent:   ptypes.TimestampNow(),
		Source: self,
		Target: target,
	}

	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", target.GetUri(), target.GetPort())
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %v", uri, err)
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)

	logger.Printf("submitting echo: %+v", msg)
	_, err = client.Echo(ctx, msg)
	if err != nil {
		logger.Printf("error while executing echo: %v", err)
		return err
	}

	completedOn := time.Now()
	sentOn, _ := ptypes.Timestamp(msg.GetSent())
	dur := completedOn.Sub(sentOn)

	logger.Printf("echo-ping: %s from %s (Duration: %v)\n ",
		msg.GetId(), msg.GetTarget().GetRegion(), dur)

	if err = save(ctx, dbPath, msg.GetId(), target.GetRegion(),
		msg.GetSource().GetRegion(), sentOn, completedOn, dur); err != nil {
		return fmt.Errorf("error while saving request: %v", err)
	}
	
	labels := map[string]string{
		"source": msg.GetSource().GetRegion(),
		"targrt": msg.GetTarget().GetRegion(),
	}
	
	if err = metric.MetricClient(ctx).Publish(ctx, "echo-duration", dur.Milliseconds(), labels); err != nil {
	  return fmt.Errorf("error while saving echo-duration metric: %v", err)	
	} 
	
	return nil

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

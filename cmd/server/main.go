package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"os"

	ptypes "github.com/golang/protobuf/ptypes"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/gcputil/env"

	"github.com/google/uuid"
	"github.com/mchmarny/gcputil/metric"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	logger   = log.New(os.Stdout, "", 0)
	grpcPort = env.MustGetEnvVar("PORT", "8080")
	dbName   = env.MustGetEnvVar("DB_NAME", "")
)

type echoService struct{}

func (s *echoService) Echo(ctx context.Context, req *pb.RequestMessage) (*pb.ResponseMessage, error) {

	// validation
	if req == nil {
		return nil, errors.New("nil request")
	}

	if req.GetTarget() == nil {
		return nil, errors.New("nil target")
	}
	if req.GetNodes() == nil {
		return nil, errors.New("nil Id")
	}
	logger.Printf("request: %+v", req)

	//parse sent time
	// TODO: create utility for this in pkg
	sentOn, err := ptypes.Timestamp(req.GetSent())
	if err != nil {
		return nil, fmt.Errorf("invalid request sent on: %v", err)
	}

	// initial ping to start the conversation
	if req.GetSource() == nil {
		rm := &pb.RequestMessage{
			Source: req.GetTarget(),
			Nodes:  req.GetNodes(),
		}

		// set target to a random node that is not self
		randNode, err := getRundomNode(req.GetNodes(), req.GetTarget(), nil)
		if err != nil {
			return nil, fmt.Errorf("invalid selecting random node: %v", err)
		}
		rm.Target = randNode

		return &pb.ResponseMessage{
			Id:      uuid.New().String(),
			Request: rm,
		}, nil
	} // done checking if new

	// save
	err = savePing(ctx, dbName, uuid.New().String(), req.GetTarget().GetRegion(),
		req.GetSource().GetRegion(), sentOn)
	if err != nil {
		return nil, fmt.Errorf("error while saving request: %v", err)
	}

	// metrics
	c, err := metric.NewClient(ctx)
	if err = c.Publish(ctx, req.GetTarget().GetRegion(), "ping", 1); err != nil {
		return nil, fmt.Errorf("error while publishing metrics: %v", err)
	}

	// response
	return &pb.ResponseMessage{
		Request: req,
	}, nil

}

func getRundomNode(nodes []*pb.EchoNode, self *pb.EchoNode, source *pb.EchoNode) (node *pb.EchoNode, err error) {
	maxLoops := 10
	i := 0
	for {
		i++
		rand.Seed(time.Now().UnixNano())
		randIndex := rand.Intn(len(nodes))
		randNode := nodes[randIndex]
		if randNode.GetRegion() != self.GetRegion() &&
			source != nil && randNode.GetRegion() != source.GetRegion() {
			return randNode, nil
		}
		if i >= maxLoops {
			return nil, errors.New("max number of rundom node selections reached")
		}
	}
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

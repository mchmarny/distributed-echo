package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/distributed-echo/pkg/client"
	"gopkg.in/yaml.v2"
)

var (
	logger = log.New(os.Stdout, "", 0)
	path   = flag.String("targets", "", "Path to server targets file")
	dbName = flag.String("db", "", "Database path")
	source = flag.String("source", "client", "Name of the invoking client ['client']")
)

func main() {

	flag.Parse()

	data, err := ioutil.ReadFile(*path)
	if err != nil {
		logger.Printf("error reading file: %v", err)
	}

	nodes := []*pb.EchoNode{}
	err = yaml.Unmarshal([]byte(data), &nodes)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	logger.Printf("Targets: %d", len(nodes))
	ctx := context.Background()
	for _, t := range nodes {
		startEcho(ctx, t, nodes)
	}

}

func startEcho(ctx context.Context, target *pb.EchoNode, nodes []*pb.EchoNode) {

	req := &pb.RequestMessage{
		Sent:   ptypes.TimestampNow(),
		Target: target,
		Nodes:  nodes,
	}

	logger.Printf("Pinging:\n   %s", req)
	resp, err := client.Ping(req)
	if err != nil {
		logger.Fatalf("error while executing Ping: %v", err)
	}

	sentOn, err := ptypes.Timestamp(resp.GetRequest().GetSent())
	if err != nil {
		logger.Fatalf("invalid response sent on: %v", err)
	}
	now := time.Now()
	dur := now.Sub(sentOn)

	logger.Printf("Echo: %s from %s (Duration: %v)\n ",
		resp.GetId(), resp.GetRequest().GetTarget().GetRegion(), dur)

}

package main

import (
	"flag"
	"log"
	"os"

	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"github.com/mchmarny/distributed-echo/pkg/client"
)

var (
	logger = log.New(os.Stdout, "", 0)
	target = flag.String("target", "", "Server address (host:port)")
	source = flag.String("source", "client", "Name of the invoking client [client]")
)

func main() {
	flag.Parse()

	req := &pb.Request{
		RequestId:  uuid.New().String(),
		SourceName: *source,
		SentOn:     ptypes.TimestampNow(),
		Targets:    []string{*target},
	}

	resp, err := client.PingClient(req)

	if err != nil {
		logger.Fatalf("Error while executing Ping: %v", err)
	}
	logger.Printf("REquest:\n  %+v", req)
	logger.Printf("Response:\n  %+v", resp)

	sentOn, _ := ptypes.Timestamp(req.GetSentOn())
	recdOn, _ := ptypes.Timestamp(resp.GetProcessedOn())
	reqDur := recdOn.Sub(sentOn)

	logger.Printf("Duration:\n  %v", reqDur)

}

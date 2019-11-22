package main

import (
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
	source = flag.String("source", "client", "Name of the invoking client ['client']")
)

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*path)
	if err != nil {
		logger.Printf("error reading file: %v", err)
	}

	targets := []pb.Target{}
	err = yaml.Unmarshal([]byte(data), &targets)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	logger.Printf("Targets: %d", len(targets))
	for _, t := range targets {
		ping(t)
	}

}

func ping(target pb.Target) {

	logger.Printf("Ping:\n   %s", target.GetRegion())
	resp, err := client.PingClient(&target)
	if err != nil {
		logger.Fatalf("error while executing Ping: %v", err)
	}

	logger.Printf("Response:\n  %+v", resp)

	sentOn, err := ptypes.Timestamp(resp.GetRequest().GetSent())
	if err != nil {
		logger.Fatalf("invalid response sent on: %v", err)
	}
	dur := time.Now().Sub(sentOn)

	logger.Printf("Duration:\n  %v", dur)

}

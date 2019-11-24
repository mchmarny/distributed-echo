package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"crypto/tls"
	"fmt"

	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"gopkg.in/yaml.v2"
)

var (
	logger = log.New(os.Stdout, "", 0)

	// ConfigFilePath contains all targets
	ConfigFilePath = flag.String("file", "", "Path to server targets file")

	// AppVersion is set at compile
	AppVersion = "0.0.0-default"

	// EchoTimeout is the max number of seconds echo service will wait for response
	EchoTimeout = 300 * time.Second
)

func main() {

	flag.Parse()

	logger.Printf("distributed-echo (v%s)", AppVersion)

	data, err := ioutil.ReadFile(*ConfigFilePath)
	if err != nil {
		logger.Printf("error reading file: %v", err)
	}

	nodes := []*pb.Node{}
	err = yaml.Unmarshal([]byte(data), &nodes)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	logger.Printf("broadcast targets: %d\n", len(nodes))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, t := range nodes {
		submitBroadcast(ctx, t, nodes)
	}

}

func submitBroadcast(ctx context.Context, from *pb.Node, to []*pb.Node) {

	logger.Println()

	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", from.GetUri(), from.GetPort())
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		logger.Printf("   failed to dial %s: %v", uri, err)
		return
	}
	defer conn.Close()

	client := pb.NewEchoServiceClient(conn)
	br, err := client.Broadcast(ctx, &pb.BroadcastMessage{
		Self:    from,
		Targets: to,
	})

	if err != nil {
		logger.Printf("   broadest error[%s]: %v", from, err)
		return
	}

	for _, nr := range br.GetResults() {
		if nr.GetError() == "" {
			logger.Printf("   %s -> %s: %dms",
				nr.GetSource().GetRegion(),
				nr.GetTarget().GetRegion(),
				nr.GetDuration(),
			)
		} else {
			logger.Printf("   %s -> %s error: %s",
				nr.GetSource().GetRegion(),
				nr.GetTarget().GetRegion(),
				nr.GetError(),
			)
		}
	}

}

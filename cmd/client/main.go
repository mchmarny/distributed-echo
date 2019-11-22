package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"os"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	pb "github.com/mchmarny/distributed-echo/pkg/api/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	serverAddr = flag.String("server", "", "Server address (host:port)")
	serverHost = flag.String("server-host", "", "Host name to which server IP should resolve")
	insecure   = flag.Bool("insecure", false, "Skip SSL validation? [false]")
	skipVerify = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
	streamSize = flag.Int("stream", 0, "Number of messages to stream [0]")
	message    = flag.String("message", "Hi there", "The body of the content sent to server")
)

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *serverHost != "" {
		opts = append(opts, grpc.WithAuthority(*serverHost))
	}
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *skipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPingServiceClient(conn)

	if *streamSize == 0 {
		send(client)
	} else {
		sendStream(client)
	}
}

func send(client pb.PingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	resp, err := client.Send(ctx, &pb.Request{

	})
	if err != nil {
		logger.Fatalf("Error while executing Send: %v", err)
	}
	logger.Printf("Unary Request/Unary Response\n Sent:\n  %s\n Response:\n  %+v", *message, resp)
}

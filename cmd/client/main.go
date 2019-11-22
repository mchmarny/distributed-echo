package main

import (
	"crypto/tls"
	"flag"
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
	source     = flag.String("source", "client", "Name of the invoking client (source:client)")
	insecure   = flag.Bool("insecure", false, "Skip SSL validation? [false]")
	skipVerify = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
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
	ping(client)

}

func ping(client pb.PingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	req := &pb.Request{
		RequestId:  "1234",
		SourceName: *source,
		SentOn:     ptypes.TimestampNow(),
		Targets:    []string{"http://localhost"}, //implement
	}
	resp, err := client.Ping(ctx, req)
	if err != nil {
		logger.Fatalf("Error while executing Ping: %v", err)
	}
	logger.Printf("Response:\n  %+v", resp)
}

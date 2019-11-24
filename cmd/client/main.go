package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"

	"crypto/tls"
	"errors"
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

	logger.Printf("version: %s", AppVersion)

	data, err := ioutil.ReadFile(*ConfigFilePath)
	if err != nil {
		logger.Printf("error reading file: %v", err)
	}

	nodes := []*pb.Node{}
	err = yaml.Unmarshal([]byte(data), &nodes)
	if err != nil {
		logger.Fatalf("error: %v", err)
	}

	logger.Printf("targets: %d", len(nodes))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan *broadcastResult, len(nodes))
		
	for _, t := range nodes {
		go func(){
			logger.Printf("broadcasting: %+v", t)
			r := submitBroadcast(ctx, &pb.BroadcastMessage{
				Self:    t,
				Targets: nodes,
			}) 
			results <- r
		}()
	}
	
	//loop:
	for {
	  select {
	  case r := <-results:
	    logger.Printf("error[%s]: %v", r.node.GetRegion(), r.err)
// 	  case <-cancel:
// 	    logger.Println("closing down...")
// 			close(results)
// 	    break loop
	  }
	}
	
}

type broadcastResult struct {
	node *pb.Node
	err  error
}

func submitBroadcast(ctx context.Context, msg *pb.BroadcastMessage) *broadcastResult {

	if msg == nil {
		return &broadcastResult{
			node: msg.GetSelf(),
			err:  errors.New("nil BroadcastMessage"),
		}
	}

	var opts []grpc.DialOption
	cred := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: false,
	})
	opts = append(opts, grpc.WithTransportCredentials(cred))
	uri := fmt.Sprintf("%s:%s", msg.GetSelf().GetUri(), msg.GetSelf().GetPort())
	conn, err := grpc.Dial(uri, opts...)
	if err != nil {
		return &broadcastResult{
			node: msg.GetSelf(),
			err:  fmt.Errorf("failed to dial %s: %v", uri, err),
		}
	}
	defer conn.Close()
	client := pb.NewEchoServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), EchoTimeout)
	defer cancel()
	_, err = client.Broadcast(ctx, msg)
	return &broadcastResult{
		node: msg.GetSelf(),
		err:  err,
	}

}

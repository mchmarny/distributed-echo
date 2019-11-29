package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/gcputil/metric"
	"gopkg.in/yaml.v2"
)

var (
	poster ContentPoster = &EchoContentPoster{}
)

// ContentPoster posts content to provided URL
type ContentPoster interface {
	Post(url string, in []byte) (out []byte, err error)
}

func pingNode(ctx context.Context, target *EchoNode) (dur int64, err error) {

	in := EchoMessage{
		From: nodeRegion,
		To:   target.Region,
		Sent: time.Now().UTC().Unix(),
	}

	dataIn, err := yaml.Marshal(in)
	if err != nil {
		return 0, fmt.Errorf("Error marshaling echo message: %v", err)
	}

	// ping
	logger.Printf("Posting echo to: %s", target.URL)
	started := time.Now()
	dataOut, err := poster.Post(target.URL, dataIn)
	finished := time.Now()
	if err != nil {
		return 0, fmt.Errorf("Error posting echo message: %v", err)
	}

	// convert
	var out EchoMessage
	if err := yaml.Unmarshal(dataOut, &out); err != nil {
		return 0, fmt.Errorf("Error decoding echo response: %v", err)
	}

	// validate
	if in != out {
		return 0, fmt.Errorf("Unexpected echo response (wanted: %v, got: %v)", in, out)
	}
	echoDuration := finished.Sub(started).Milliseconds()
	logger.Printf("echo-ping from: %s to: %s (duration: %v)\n ",
		nodeRegion, target.Region, echoDuration)

	// save
	if err := save(ctx, dbPath, uuid.New().String(), nodeRegion, target.Region,
		started, finished, echoDuration); err != nil {
		return echoDuration, fmt.Errorf("Error while saving results: %v", err)
	}

	// metrics
	labels := map[string]string{
		"source": nodeRegion,
		"target": target.Region,
	}
	if err := metric.MakeClient(ctx).Publish(ctx, metricName, echoDuration, labels); err != nil {
		// more then 1 metric per sec will cause an error
		logger.Printf("Non-fatal error while publishing metrics: %v", err)
	}

	return echoDuration, nil

}

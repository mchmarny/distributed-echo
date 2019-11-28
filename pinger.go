package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/gcputil/metric"
	"gopkg.in/yaml.v2"
)

var (
	poster ContentPoster = &EchoContentPoster{}
)

func pingNode(ctx context.Context, target *EchoNode) error {

	in := EchoMessage{
		From: nodeRegion,
		To:   target.Region,
		Sent: time.Now().UTC().Unix(),
	}

	dataIn, err := yaml.Marshal(in)
	if err != nil {
		logger.Printf("Error marshaling echo message: %v", err)
		return err
	}

	// ping
	logger.Printf("Posting echo to: %s", target.URL)
	started := time.Now()
	dataOut, err := poster.Post(target.URL, dataIn)
	finished := time.Now()
	if err != nil {
		logger.Printf("Error posting echo message: %v", err)
		return err
	}

	// convert
	var out EchoMessage
	if err := yaml.Unmarshal(dataOut, &out); err != nil {
		logger.Printf("Error decoding echo response: %v", err)
		return err
	}

	// validate
	if in != out {
		logger.Printf("Unexpected echo response (wanted: %v, got: %v)", in, out)
		return err
	}
	echoDuration := finished.Sub(started).Milliseconds()
	logger.Printf("echo-ping from: %s to: %s (duration: %v)\n ",
		nodeRegion, target.Region, echoDuration)

	// save
	if err := save(ctx, dbPath, uuid.New().String(), nodeRegion, target.Region,
		started, finished, echoDuration); err != nil {
		logger.Printf("Error while saving results: %v", err)
		return err
	}

	// metrics
	labels := map[string]string{
		"source": nodeRegion,
		"target": target.Region,
	}
	if err := metric.MakeClient(ctx).Publish(ctx, "echo-duration", echoDuration, labels); err != nil {
		logger.Printf("Error while publishing metrics: %v", err)
		return err
	}

	return nil

}

// ContentPoster posts content to provided URL
type ContentPoster interface {
	Post(url string, in []byte) (out []byte, err error)
}

// EchoContentPoster posts to echo endpoint
type EchoContentPoster struct{}

// Post posts to echo endpoint
func (p *EchoContentPoster) Post(url string, in []byte) (out []byte, err error) {

	resp, err := http.Post(url, "text/x-yaml", bytes.NewBuffer(in))
	if err != nil {
		logger.Printf("Error posting echo message: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid post response code: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("Error reading response body content: %v", err)
		return nil, err
	}

	return data, nil

}

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/google/uuid"
	"github.com/mchmarny/gcputil/metric"
	"gopkg.in/yaml.v2"
)

var (
	poster ContentPoster = &EchoContentPoster{}
)

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
	if err := metric.MakeClient(ctx).Publish(ctx, "echo-duration", echoDuration, labels); err != nil {
		// more then 1 metric per sec will cause an error
		logger.Printf("Non-fatal error while publishing metrics: %v", err)
	}

	return echoDuration, nil

}

// ContentPoster posts content to provided URL
type ContentPoster interface {
	Post(url string, in []byte) (out []byte, err error)
}

// EchoContentPoster posts to echo endpoint
type EchoContentPoster struct{}

// Post posts to echo endpoint
func (p *EchoContentPoster) Post(url string, in []byte) (out []byte, err error) {

	// get auth token from metadata server
	tokenURL := fmt.Sprintf("/instance/service-accounts/default/identity?audience=%s", url)
	idToken, err := metadata.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("Error getting metadata: %v", err)
	}

	// create request
	logger.Printf("HTTP Post to %s with %d bytes", url, len(in))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(in))
	if err != nil {
		return nil, fmt.Errorf("Error creating posting request: %v", err)
	}
	req.Header.Add("Content-Type", "text/x-yaml")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", idToken))

	// process response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error posting echo message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Invalid post response code: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body content: %v", err)
	}
	logger.Printf("HTTP Post to %s returned %d bytes", url, len(data))

	return data, nil

}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/gcputil/metric"
	"gopkg.in/yaml.v2"
)

func pingNode(ctx context.Context, target *EchoNode) error {

	msg := &EchoMessage{
		From: nodeRegion,
		To:   target.Region,
		Sent: time.Now(),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Printf("Error marshaling echo message: %v", err)
		return err
	}

	logger.Printf("Posting echo to: %s", target.URL)
	started := time.Now()
	resp, err := http.Post(target.URL, "text/x-yaml", bytes.NewBuffer(data))
	finished := time.Now()
	if err != nil {
		logger.Printf("Error posting echo message: %v", err)
		return err
	}

	var out EchoMessage
	if err := yaml.NewDecoder(resp.Body).Decode(&out); err != nil {
		logger.Printf("Error decoding echo response: %v", err)
		return err
	}

	if *msg != out {
		logger.Printf("Unexpected echo response (wanted: %v, got: %v)", *msg, out)
		return err
	}

	echoDuration := finished.Sub(started).Milliseconds()
	logger.Printf("echo-ping from: %s to: %s (duration: %v)\n ",
		nodeRegion, target.Region, echoDuration)

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

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/compute/metadata"
)

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

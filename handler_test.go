package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestEchoHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)
	in := EchoMessage{
		From: nodeRegion,
		To:   "http://localhost:8080/",
		Sent: time.Now().UTC().Unix(),
	}

	b, err := yaml.Marshal(in)
	assert.Nil(t, err)

	r := gin.Default()
	r.GET("/", echoHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", bytes.NewBuffer(b))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var out EchoMessage
	err = yaml.NewDecoder(w.Result().Body).Decode(&out)
	assert.Nilf(t, err, "Error decoding body: %v", err)

	assert.Equal(t, in, out)

}

func TestBroadcastHandler(t *testing.T) {

	gin.SetMode(gin.ReleaseMode)
	testFile := "etc/test.yaml"
	data, err := ioutil.ReadFile(testFile)
	assert.Nilf(t, err, "Error reading file (%s): %v", testFile, err)

	var msg BroadcastMessage
	err = yaml.Unmarshal([]byte(data), &msg)
	assert.Nilf(t, err, "Error decoding file content: %v", err)

	poster = &TestEchoContentPoster{}

	r := gin.Default()
	r.POST("/", broadcastHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

}

// TestEchoContentPoster is a test implementation of the poster
type TestEchoContentPoster struct{}

// Post posts to echo endpoint
func (p *TestEchoContentPoster) Post(url string, in []byte) (out []byte, err error) {
	if url == "" {
		return nil, errors.New("nil URL")
	}
	return in, nil
}

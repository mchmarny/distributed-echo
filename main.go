package main

import (
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mchmarny/gcputil/env"
	"github.com/mchmarny/gcputil/project"
)

var (
	logger = log.New(os.Stdout, "[ECHO] ", 0)

	projectID = project.GetIDOrFail()

	port       = env.MustGetEnvVar("PORT", "8080")
	dbPath     = env.MustGetEnvVar("DBP", "")
	nodeRegion = env.MustGetEnvVar("REG", "")
	metricName = env.MustGetEnvVar("MET", "echo-latency")
	release    = env.MustGetEnvVar("REL", "v0.0.0-default")
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// simple routes
	r.GET("/", defaultHandler)

	// api
	v1 := r.Group("/v1")
	{
		v1.POST("/broadcast", broadcastHandler)
		v1.POST("/echo", echoHandler)
	}

	// server
	hostPort := net.JoinHostPort("0.0.0.0", port)
	logger.Printf("Server starting: %s \n", hostPort)
	if err := r.Run(hostPort); err != nil {
		logger.Fatal(err)
	}
}

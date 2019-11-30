package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func defaultHandler(c *gin.Context) {

	dump, err := httputil.DumpRequest(c.Request, true)
	if err != nil {
		c.YAML(http.StatusInternalServerError, gin.H{
			"code":   string(http.StatusInternalServerError),
			"status": err.Error(),
		})
		return
	}

	c.YAML(http.StatusOK, gin.H{
		"release": ReleaseVersion,
		"request": string(dump),
	})
}

func broadcastHandler(c *gin.Context) {

	var msg BroadcastMessage
	if err := c.ShouldBindYAML(&msg); err != nil {
		logger.Printf("Error binding broadcast YAML: %v", err)
		c.YAML(http.StatusBadRequest, gin.H{
			"message": "Invalid broadcast message format",
			"status":  "Error",
		})
		return
	}

	results := make(map[string]string, len(msg.Targets))
	for _, n := range msg.Targets {
		logger.Printf("Target: %+v", n)
		dur, err := pingNode(c.Request.Context(), n)
		if err != nil {
			results[n.Region] = err.Error()
			logger.Printf("Error pinging %s: %v", n.Region, err)
		} else {
			results[n.Region] = fmt.Sprintf("OK: %dms", dur)
		}
	}

	c.YAML(http.StatusOK, results)

}

func echoHandler(c *gin.Context) {

	var msg EchoMessage
	if err := c.ShouldBindYAML(&msg); err != nil {
		logger.Printf("Error binding echo YAML: %v", err)
		c.YAML(http.StatusBadRequest, gin.H{
			"message": "Invalid echo message format",
			"status":  "Error",
		})
		return
	}

	logger.Printf("Echo message: %v", msg)
	c.YAML(http.StatusOK, msg)

}

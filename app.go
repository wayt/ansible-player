package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	envBindAddress     = "BIND_ADDRESS"
	defaultBindAddress = ":8080"
)

func main() {

	r := gin.Default()
	r.Use(AuthMiddleware())

	// Setup seed
	rand.Seed(time.Now().UnixNano())

	// Register actions
	r.POST("/job", postJobAction)
	r.GET("/job/:id", getJobAction)

	bind := defaultBindAddress
	if e := os.Getenv(envBindAddress); e != "" {
		bind = e
	}

	r.Run(bind)
}

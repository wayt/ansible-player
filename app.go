package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"math/rand"
	"time"
)

func main() {

	r := gin.Default()
	r.Use(AuthMiddleware())

	// Setup seed
	rand.Seed(time.Now().UnixNano())

	// Register actions
	r.POST("/job", postJobAction)
	r.GET("/job/:id", getJobAction)

	r.Run(env.GetDefault("BIND_ADDRESS", ":8080"))
}

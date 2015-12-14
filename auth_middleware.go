package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gotoolz/env"
	"log"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if !validateAuth(username, password) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func validateAuth(username, password string) bool {

	// Password sha1
	password = fmt.Sprintf("%x", sha1.Sum([]byte(password)))

	// Re-open auth file each time, to avoid reloading it
	file, err := os.Open(env.GetDefault("AUTH_FILE", "access"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Scan line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0] == '#' {
			continue
		}

		// Split format: username:sha1_password
		splited := strings.SplitN(line, `:`, 2)
		if len(splited) != 2 {
			continue
		}

		log.Println(splited, username, password)
		// Check credentials
		if splited[0] == username && splited[1] == password {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	log.Println("Bad password")

	return false

}

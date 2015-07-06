package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/wayt/happyngine"
	"github.com/wayt/happyngine/env"
	"os"
	"strings"
)

type AuthMiddleware struct {
	happyngine.Middleware
}

func NewAuthMiddleware() happyngine.MiddlewareHandler {

	return func(context *happyngine.Context) happyngine.MiddlewareInterface {

		this := &AuthMiddleware{happyngine.Middleware{context}}
		return this
	}
}

func (this *AuthMiddleware) HandleBefore() (err error) {

	username, password, ok := this.Context.Request.BasicAuth()
	if !ok {
		this.Context.Send(403, `Missing http auth`)
		return errors.New("Missing http auth")
	}

	if !this.validate(username, password) {
		this.Context.Send(403, `Unauthorized`)
		return errors.New("Unauthorized")
	}

	return nil
}

func (this *AuthMiddleware) HandleAfter() error {

	return nil
}

func (this *AuthMiddleware) validate(username, password string) bool {

	// Password sha1
	password = fmt.Sprintf("%x", sha1.Sum([]byte(password)))

	// Re-open auth file each time, to avoid reloading it
	file, err := os.Open(env.Get("AUTH_FILE"))
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

		// Check credentials
		if splited[0] == username && splited[1] == password {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return false

}

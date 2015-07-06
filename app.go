package main

import (
	"github.com/wayt/happyngine"
	"github.com/wayt/happyngine/log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func main() {

	app := happyngine.NewAPI()

	// Setup seed
	rand.Seed(time.Now().UnixNano())

	// Register actions
	app.AddRoute("POST", "/job", newPostJobAction, NewAuthMiddleware())

	// Setup custuom 404 handler
	app.Error404Handler = func(ctx *happyngine.Context, err interface{}) {

		ctx.Send(http.StatusNotFound, `not found 404`)
	}

	// Setup custuom panic handler
	app.PanicHandler = func(ctx *happyngine.Context, err interface{}) {

		ctx.Send(500, `internal error`)

		trace := make([]byte, 1024)
		runtime.Stack(trace, true)

		ctx.Criticalln(err, string(trace))
	}

	log.Debugln("Running...")
	if err := app.Run(":8080"); err != nil {
		log.Criticalln("app.Run:", err)
	}
}

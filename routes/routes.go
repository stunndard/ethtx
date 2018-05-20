package routes

import (
	"log"
	"time"

	"gopkg.in/kataras/iris.v6"
)

var RedactLogs bool = false

func Headers(ctx *iris.Context) {
	ctx.SetHeader("Server", "X-Files")
	ctx.Next()
}

func Log(ctx *iris.Context) {
	log.Println("[B]", ctx.RemoteAddr(), ctx.Method(), ctx.Path())
	timeStarted := time.Now()
	ctx.Next()
	timeElapsed := time.Now().Sub(timeStarted)

	log.Println("[E]", ctx.RemoteAddr(), ctx.Method(), ctx.StatusCode(), ctx.Path(), timeElapsed)
	log.Println("---")
}

func NotFound(ctx *iris.Context ) {
	log.Println(ctx.RemoteAddr(), ctx.Method(), ctx.StatusCode(), ctx.Path())
	renderJSONError(ctx, "404", "not found")
}

func renderJSONError(ctx *iris.Context, code, message string) {
	ctx.JSON(
		iris.StatusOK,
		&iris.Map{
			"error": &iris.Map{
				"code":    code,
				"message": message,
			},
		},
	)
	ctx.StopExecution()
}

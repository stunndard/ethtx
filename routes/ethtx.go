package routes

import (
	"time"

	"gopkg.in/kataras/iris.v6"
)

var (
	Started time.Time
	Version string
)

// show health
func ApiHealth(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, &iris.Map{"health": "OK"})
}

// show version
func ApiVersion(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, &iris.Map{
		"success": true,
		"status":  200,
		"version": Version,
		"started": Started,
	})
}
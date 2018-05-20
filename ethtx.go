package main

import (
	"github.com/stunndard/ethtx/routes"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

func startWebApp() {
	// load config
	app := iris.New(iris.YAML("ethtx.yml"))

	routes.RedactLogs = app.Config.Other["redactLogs"].(bool) == true

	app.Adapt(iris.DevLogger())
	app.Adapt(httprouter.New())

	// first handlers
	app.UseFunc(routes.Log)
	app.UseFunc(routes.Headers)

	// 404 handler
	app.OnError(iris.StatusNotFound, routes.NotFound)

	// app handlers
	app.Get("/health", routes.ApiHealth)
	app.Get("/", routes.ApiVersion)

	// sign tx
	app.Post("/signTx", routes.CheckTx, routes.SignTx)
	// generate address
	app.Post("generateAddress", routes.GenerateAddress)

	// start the web app
	app.Listen(app.Config.Other["listen"].(string))
}

func main() {



	startWebApp()
}

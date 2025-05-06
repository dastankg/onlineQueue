package main

import (
	"net/http"
	"onlineQueue/internal/app"
)

// @title OnlineQueue documentation
// @version 1.0.1
// @host 127.0.0.1:8080
// @BasePath
func main() {
	app := app.App()
	server := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

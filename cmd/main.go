package main

import (
	"net/http"
	"onlineQueue/internal/app"
)

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

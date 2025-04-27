package app

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/internal/auth"
	"onlineQueue/internal/operators"
	db2 "onlineQueue/pkg/db"
	"onlineQueue/pkg/middleware"
)

func App() http.Handler {
	conf, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := db2.NewDb(conf)
	if err != nil {
		panic(err)
	}
	router := http.NewServeMux()

	operatorsRepositoru := operators.NewOperatorRepository(db)
	authService := auth.NewAuthService(operatorsRepositoru)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		conf,
		authService,
	})
	stack := middleware.Chain(middleware.CORS, middleware.Logging)
	return stack(router)
}

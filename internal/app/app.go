package app

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/internal/auth"
	"onlineQueue/internal/offices"
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

	operatorsRepository := operators.NewOperatorRepository(db)
	registersRepository := offices.NewOfficeRepository(db)
	authService := auth.NewAuthService(operatorsRepository)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		conf,
		authService,
	})
	offices.NewOfficeHandler(router, offices.OfficeHandlerDeps{
		registersRepository,
		conf,
	})
	stack := middleware.Chain(middleware.CORS, middleware.Logging)
	return stack(router)
}

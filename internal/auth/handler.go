package auth

import (
	"net/http"
	"onlineQueue/configs"
	"onlineQueue/pkg/jwt"
	"onlineQueue/pkg/req"
	"onlineQueue/pkg/res"
	"time"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}

		login, err := handler.AuthService.Register(body.Name, body.Login, body.Password, false, body.RegisterID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)
		tokens, err := jwtService.CreateTokenPair(
			login,
			15*time.Minute,
			24*7*time.Hour,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Name:         login,
			Login:        login,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		res.Json(w, data, 201)
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		login, err := handler.AuthService.Login(body.Login, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)
		tokens, err := jwtService.CreateTokenPair(
			login,
			15*time.Minute, // access token на 15 минут
			24*7*time.Hour, // refresh token на 7 дней
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Login:        login,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		res.Json(w, data, 201)
	}
}

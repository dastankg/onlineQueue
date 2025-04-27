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
	//router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		var registerID *uint
		if body.RegisterID != nil {
			registerID = body.RegisterID
		} else {
			registerID = nil
		}
		login, err := handler.AuthService.Register(body.Name, body.Login, body.Password, false, registerID)
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

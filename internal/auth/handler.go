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
	router.HandleFunc("POST /login", handler.Login())
	router.HandleFunc("POST /register", handler.Register())
}

// Register регистрация пользователя и возвращает пару токенов.
//
// @Summary 	Регистрация в систему
// @Description Регистрация по логину и паролю, возвращает access и refresh токены
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param 		body body RegisterRequest true "Данные для регистрации"
// @Success     201    {object}  RegisterResponse
// @Failure     400    {string}  string  "Неверный логин или пароль"
// @Failure     500    {string}  string  "Ошибка генерации токена"
// @Router      /register [post]
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}

		login, err := handler.AuthService.Register(body.Name, body.Login, body.Password1, body.Password2, false, body.RegisterID)
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

// Login аутентифицирует пользователя и возвращает пару токенов.
//
// @Summary 	Вход в систему
// @Description Авторизация пользователя по логину и паролю, возвращает access и refresh токены
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       body  body      LoginRequest  true  "Данные для входа"
// @Success     201    {object}  LoginResponse
// @Failure     400    {string}  string  "Неверный логин или пароль"
// @Failure     500    {string}  string  "Ошибка генерации токена"
// @Router      /login [post]
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

// @Summary Обновление токена доступа
// @Description Обновляет access token используя refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh токен"
// @Success 200 {object} RefreshResponse "Новая пара токенов"
// @Router /auth/refresh [post]
func (handler *AuthHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RefreshRequest](&w, r)
		if err != nil {
			return
		}

		jwtService := jwt.NewJWT(
			handler.Config.Auth.AccessSecret,
			handler.Config.Auth.RefreshSecret,
		)

		flag, claims := jwtService.ParseRefreshToken(body.RefreshToken)
		if !flag {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Unix(claims.ExpiresAt.Unix(), 0)
		remainingTime := time.Until(expirationTime)

		if remainingTime <= 0 {
			http.Error(w, "Refresh token has expired", http.StatusUnauthorized)
			return
		}

		accessToken, err := jwtService.Create(jwt.JWTData{
			Login:     claims.Login,
			ExpiresAt: time.Now().Add(15 * time.Minute),
			TokenType: "access",
		}, jwtService.AccessSecret)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := RefreshResponse{
			AccessToken: accessToken,
		}

		res.Json(w, data, http.StatusOK)
	}
}

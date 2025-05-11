package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	AccessToken  string = "access"
	RefreshToken string = "refresh"
)

type JWTData struct {
	Login     string
	ExpiresAt time.Time
	TokenType string
}

type JWT struct {
	AccessSecret  string
	RefreshSecret string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func NewJWT(accessSecret, refreshSecret string) *JWT {
	return &JWT{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
	}
}

func (j *JWT) Create(data JWTData, secret string) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":      data.Login,
		"exp":        data.ExpiresAt.Unix(),
		"token_type": data.TokenType,
	})
	return t.SignedString([]byte(secret))
}

func (j *JWT) CreateTokenPair(login string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	accessToken, err := j.Create(JWTData{
		Login:     login,
		ExpiresAt: time.Now().Add(accessTTL),
		TokenType: AccessToken,
	}, j.AccessSecret)
	if err != nil {
		return nil, err
	}
	refreshToken, err := j.Create(JWTData{
		Login:     login,
		ExpiresAt: time.Now().Add(refreshTTL),
		TokenType: RefreshToken,
	}, j.RefreshSecret)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWT) parse(token string, secret string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return false, nil
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}
	login, ok := claims["login"].(string)
	if !ok {
		return false, nil
	}
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, nil
	}
	tokenType, ok := claims["token_type"].(string)
	if !ok {
		return false, nil
	}
	return t.Valid, &JWTData{
		Login:     login,
		ExpiresAt: time.Unix(int64(exp), 0),
		TokenType: tokenType,
	}
}

func (j *JWT) ParseAccessToken(token string) (bool, *JWTData) {
	return j.parse(token, j.AccessSecret)
}

func (j *JWT) ParseRefreshToken(token string) (bool, *JWTData) {
	return j.parse(token, j.RefreshSecret)
}

func (j *JWT) Refresh(refreshToken string, accessTTL, refreshTTL time.Duration) (*TokenPair, error) {
	valid, date := j.ParseRefreshToken(refreshToken)
	if !valid || date == nil {
		return nil, jwt.ErrSignatureInvalid
	}

	if date.TokenType != RefreshToken {
		return nil, jwt.ErrTokenInvalidId
	}

	if time.Now().After(date.ExpiresAt) {
		return nil, fmt.Errorf("токен истек")
	}
	return j.CreateTokenPair(date.Login, accessTTL, refreshTTL)
}

package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"onlineQueue/internal/operators"
)

type AuthService struct {
	OperatorRepository *operators.OperatorRepository
}

func NewAuthService(operatorRepository *operators.OperatorRepository) *AuthService {
	return &AuthService{
		OperatorRepository: operatorRepository,
	}
}

func (service *AuthService) Register(name, login, password string, isAdmin bool, registerID *uint) (string, error) {
	existedOperator, err := service.OperatorRepository.FindByLogin(login)
	if err != nil {
		return "", err
	}
	if existedOperator != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	operator := &operators.Operator{
		Name:       name,
		Login:      login,
		Password:   string(hashedPassword),
		IsAdmin:    isAdmin,
		RegisterID: registerID,
	}
	_, err = service.OperatorRepository.CreateOperator(operator)
	if err != nil {
		return "", err
	}
	return operator.Login, nil
}

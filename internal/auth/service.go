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

	if existedOperator != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	operator := &operators.Operator{
		Name:     name,
		Login:    login,
		Password: string(hashedPassword),
		IsAdmin:  isAdmin,
		OfficeID: registerID,
	}
	_, err = service.OperatorRepository.CreateOperator(operator)
	if err != nil {
		return "", err
	}
	return operator.Login, nil
}

func (service *AuthService) Login(login, password string) (string, error) {
	existedOperator, err := service.OperatorRepository.FindByLogin(login)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(existedOperator.Password), []byte(password))
	if err != nil {
		return "", err
	}
	return existedOperator.Login, nil
}

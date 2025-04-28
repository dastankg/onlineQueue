package operators

import (
	"onlineQueue/pkg/db"
)

type OperatorRepository struct {
	database *db.DB
}

func NewOperatorRepository(database *db.DB) *OperatorRepository {
	return &OperatorRepository{database: database}
}

func (repo *OperatorRepository) CreateOperator(operator *Operator) (*Operator, error) {
	result := repo.database.DB.Create(operator)
	if result.Error != nil {
		return nil, result.Error
	}
	return operator, nil
}

func (repo *OperatorRepository) FindByLogin(login string) (*Operator, error) {
	var operator Operator
	result := repo.database.DB.Where("login = ?", login).First(&operator)

	if result.Error != nil {
		return nil, result.Error
	}
	return &operator, nil
}

package registers

import "onlineQueue/pkg/db"

type RegisterRepository struct {
	Database *db.DB
}

//func NewRegisterRepository(db *db.DB) *RegisterRepository {
//	return &RegisterRepository{
//		Database: db,
//	}
//}

func (repo *RegisterRepository) CreateRegister(register *Register) (*Register, error) {
	result := repo.Database.DB.Create(register)
	if result.Error != nil {
		return nil, result.Error
	}
	return register, nil
}

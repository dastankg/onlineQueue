package offices

import "onlineQueue/pkg/db"

type OfficeRepository struct {
	Database *db.DB
}

func NewOfficeRepository(db *db.DB) *OfficeRepository {
	return &OfficeRepository{
		Database: db,
	}
}

func (repo *OfficeRepository) CreateOffice(register *Offices) (*Offices, error) {
	result := repo.Database.DB.Create(register)
	if result.Error != nil {
		return nil, result.Error
	}
	return register, nil
}

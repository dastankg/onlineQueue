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

func (repo *OfficeRepository) CreateOffice(register *Office) (*Office, error) {
	result := repo.Database.DB.Create(register)
	if result.Error != nil {
		return nil, result.Error
	}
	return register, nil
}

func (repo *OfficeRepository) GetOffices() []Office {
	var offices []Office
	repo.Database.Table("offices").Find(&offices)
	return offices
}

func (repo *OfficeRepository) GetOfficeById(id uint) (*Office, error) {
	var office Office
	result := repo.Database.DB.First(&office, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &office, nil
}

func (repo *OfficeRepository) DeleteOffice(id uint) error {
	result := repo.Database.DB.Delete(&Office{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

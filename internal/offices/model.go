package offices

import (
	"gorm.io/gorm"
)

type Office struct {
	gorm.Model
	Name         string
	Address      string
	WorkingHours string
}

func NewOffice(name string, address string, workingHours string) *Office {
	return &Office{
		Name:         name,
		Address:      address,
		WorkingHours: workingHours,
	}
}

package offices

import (
	"gorm.io/gorm"
)

type Offices struct {
	gorm.Model
	Name         string
	Address      string
	WorkingHours string
}

func NewOffice(name string, address string, workingHours string) *Offices {
	return &Offices{
		Name:         name,
		Address:      address,
		WorkingHours: workingHours,
	}
}

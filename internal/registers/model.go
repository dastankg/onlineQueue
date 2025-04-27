package registers

import (
	"gorm.io/gorm"
)

type Register struct {
	gorm.Model
	Name         string
	Address      string
	WorkingHours string
}

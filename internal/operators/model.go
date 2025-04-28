package operators

import (
	"gorm.io/gorm"
	"onlineQueue/internal/offices"
)

type Operator struct {
	gorm.Model
	Name        string
	Login       string `gorm:"uniqueIndex"`
	Password    string
	IsActive    bool `gorm:"default:true"`
	IsAdmin     bool `gorm:"default:false"`
	TableNumber int
	RegisterID  *uint
	Register    offices.Office `gorm:"foreignKey:RegisterID"`
}

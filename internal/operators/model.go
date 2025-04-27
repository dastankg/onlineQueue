package operators

import (
	"gorm.io/gorm"
	"onlineQueue/internal/registers"
)

type Operator struct {
	gorm.Model
	Name       string
	Login      string `gorm:"uniqueIndex"`
	Password   string
	IsActive   bool `gorm:"default:true"`
	IsAdmin    bool `gorm:"default:false"`
	RegisterID *uint
	Register   registers.Register `gorm:"foreignKey:RegisterID"`
}

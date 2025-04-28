package queuee

import "gorm.io/gorm"

type QueueClient struct {
	gorm.Model
	OfficeID uint
	Number   uint
	Status   string
}

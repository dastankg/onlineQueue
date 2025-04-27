package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"onlineQueue/configs"
)

type DB struct {
	*gorm.DB
}

func NewDb(conf *configs.Config) (*DB, error) {
	db, err := gorm.Open(postgres.Open(conf.Db.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

package repository

import (
	"book-management/iface"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) iface.Respository {
	return &repository{
		db: db,
	}
}

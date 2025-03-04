package models

import (
	"github.com/go-playground/validator/v10"
)

type Book struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Title  string `gorm:"type:varchar(255);not null" json:"title" validate:"required"`
	Author string `gorm:"type:varchar(255);not null" json:"author" validate:"required"`
	Year   int    `gorm:"not null" json:"year" validate:"gte=1000,lte=2100"`
}

type BookEvent struct {
	EventType string `json:"event_type"`
	Book      Book   `json:"book"`
}

func (Book) TableName() string {
	return "books"
}

// ValidateBook
func (b *Book) ValidateCreateBook() error {
	validate := validator.New()
	return validate.Struct(b)
}

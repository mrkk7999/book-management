package iface

import (
	models "book-management/request_response/books"
)

type Service interface {
	GetBookByID(id int) (*models.Book, error)
	GetAllBooks(page, limit int) ([]models.Book, int64, error)
	CreateBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

type ConsumerService interface {
	CreateBook(req models.BookEvent)
	UpdateBook(req models.BookEvent)
	DeleteBook(req models.BookEvent)
}

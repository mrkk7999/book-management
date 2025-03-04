package iface

import models "book-management/request_response/books"

type Respository interface {
	GetAllBooks(page, limit int) ([]models.Book, int64, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book *models.Book) error
	UpdateBook(book *models.Book) error
	DeleteBook(id int) error
}

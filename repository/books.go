package repository

import models "book-management/request_response/books"

func (r *repository) GetAllBooks(page, limit int) ([]models.Book, int64, error) {
	var books []models.Book
	var total int64

	// Count total books
	r.db.Model(&models.Book{}).Count(&total)

	// Pagination logic
	offset := (page - 1) * limit
	result := r.db.Order("id ASC").Limit(limit).Offset(offset).Find(&books)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return books, total, nil
}

func (r *repository) GetBookByID(id int) (*models.Book, error) {
	var book models.Book
	result := r.db.First(&book, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &book, nil
}

func (r *repository) CreateBook(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *repository) UpdateBook(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *repository) DeleteBook(id int) error {
	return r.db.Delete(&models.Book{}, id).Error
}

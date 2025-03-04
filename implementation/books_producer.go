package implementation

import (
	"book-management/kafka/producer"
	models "book-management/request_response/books"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// GetBookByID
func (s *service) GetBookByID(id int) (*models.Book, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("book:%d", id)

	cachedBook, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedBook != "" {
		var book models.Book
		json.Unmarshal([]byte(cachedBook), &book)
		return &book, nil
	}

	// Fetch from database if not found in cache
	book, err := s.repo.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	bookJSON, _ := json.Marshal(book)
	s.cache.Set(ctx, cacheKey, string(bookJSON), 10*time.Minute)

	return book, nil
}

// GetAllBooks
func (s *service) GetAllBooks(page, limit int) ([]models.Book, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("books:page%d_limit%d", page, limit)

	cachedBooks, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedBooks != "" {
		var books []models.Book
		json.Unmarshal([]byte(cachedBooks), &books)
		total := int64(len(books))
		return books, total, nil
	}

	books, total, err := s.repo.GetAllBooks(page, limit)
	if err != nil {
		return nil, 0, err
	}

	booksJSON, _ := json.Marshal(books)
	s.cache.Set(ctx, cacheKey, string(booksJSON), 10*time.Minute)

	return books, total, nil
}

// CreateBook
func (s *service) CreateBook(book *models.Book) error {
	// Cache Invalidation
	s.cache.Delete(context.Background(), "books:*")

	bookEvent := models.BookEvent{
		Book:      *book,
		EventType: "create_book",
	}
	message, err := json.Marshal(bookEvent)
	if err != nil {
		s.log.Error("Error marshalling create book request:", err)
		return err
	}
	producer.PublishMessageAsynchronous(s.asyncProducer, s.topic, string(message))
	return nil
}

// UpdateBook
func (s *service) UpdateBook(book *models.Book) error {
	// Cache Invalidation
	s.cache.Delete(context.Background(), fmt.Sprintf("book:%d", book.ID))
	s.cache.Delete(context.Background(), "books:*")

	bookEvent := models.BookEvent{
		Book:      *book,
		EventType: "update_book",
	}
	message, err := json.Marshal(bookEvent)
	if err != nil {
		s.log.Error("Error marshalling update book request", err)
		return err
	}
	producer.PublishMessageAsynchronous(s.asyncProducer, s.topic, string(message))
	return nil
}

// DeleteBook
func (s *service) DeleteBook(id int) error {
	// Cache Invalidation
	s.cache.Delete(context.Background(), fmt.Sprintf("book:%d", id))
	s.cache.Delete(context.Background(), "books:*")

	bookEvent := models.BookEvent{
		Book: models.Book{
			ID: uint(id),
		},
		EventType: "delete_book",
	}
	message, err := json.Marshal(bookEvent)
	if err != nil {
		s.log.Error("Error marshalling book deletion request", err)
		return err
	}
	producer.PublishMessageAsynchronous(s.asyncProducer, s.topic, string(message))
	return nil
}

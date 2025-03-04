package implementation

import (
	models "book-management/request_response/books"
	"log"
)

func (s *consumerService) CreateBook(req models.BookEvent) {
	err := s.repository.CreateBook(&req.Book)
	if err != nil {
		log.Println("Error creating book", err)
	}
}

func (s *consumerService) UpdateBook(req models.BookEvent) {
	err := s.repository.UpdateBook(&req.Book)
	if err != nil {
		log.Println("Error updaing book", err)
	}
}

func (s *consumerService) DeleteBook(req models.BookEvent) {
	err := s.repository.DeleteBook(int(req.Book.ID))
	if err != nil {
		log.Println("Error deleting book", err)
	}
}

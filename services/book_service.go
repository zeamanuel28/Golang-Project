package services

import (
	"gocheck/models"

	"gorm.io/gorm"
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

// CreateBook adds a new book to the database

func (s *BookService) CreateBook(book *models.Book) (*models.Book, error) {
	err := s.db.Create(book).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetBooksByUserID returns all books for a specific user
func (s *BookService) GetBooksByUserID(userID uint) ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Where("user_id = ?", userID).Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// GetBookByID gets a single book
func (s *BookService) GetBookByID(id uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (bs *BookService) UpdateBook(book *models.Book) (*models.Book, error) {
	err := bs.db.Save(book).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

// DeleteBook deletes a single book
func (s *BookService) DeleteBook(id uint) error {
	if err := s.db.Delete(&models.Book{}, id).Error; err != nil {
		return err
	}
	return nil
}
func (bs *BookService) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := bs.db.Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

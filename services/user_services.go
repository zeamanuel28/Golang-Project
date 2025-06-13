package services

import (
	"errors"
	"gocheck/models"
	"gocheck/utils"

	"gorm.io/gorm"
)

// UserService provides business logic for user operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new UserService
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(user *models.User) error {
	// Hash the password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	if err := s.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsersPaginated retrieves a limited set of users based on pagination parameters
func (s *UserService) GetUsersPaginated(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total users for frontend pagination info
	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch the limited users based on limit & offset
	if err := s.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *models.User) (*models.User, error) {
	// First, check if the user exists
	existingUser, err := s.GetUserByID(user.ID)
	if err != nil {
		return nil, err // User not found or other DB error
	}

	// Update only the fields that are provided (e.g., username, email).
	// Password update should be handled separately if it's a security concern.
	existingUser.Username = user.Username
	existingUser.Email = user.Email
	existingUser.Name.FirstName = user.Name.FirstName
	existingUser.Name.LastName = user.Name.LastName

	if err := s.db.Save(existingUser).Error; err != nil {
		return nil, err
	}
	return existingUser, nil
}

// DeleteUser deletes a user by their ID
func (s *UserService) DeleteUser(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

// AuthenticateUser authenticates a user by email and password
// It returns the authenticated user if successful, or an error.
func (s *UserService) AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User

	// 1. Find the user by email
	// Use .Limit(1) to ensure only one user is returned, though email is unique
	if err := s.db.Where("email = ?", email).Limit(1).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials") // Return a generic error for security
		}
		return nil, err // Other database error
	}

	// 2. Check if the user was found (should be caught by gorm.ErrRecordNotFound, but good to be explicit)
	if user.ID == 0 {
		return nil, errors.New("invalid credentials")
	}

	// 3. Compare the provided password with the hashed password from the database
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials") // Return a generic error for security
	}

	// Authentication successful
	return &user, nil
}

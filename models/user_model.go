package models

import "gorm.io/gorm"

// User represents the user model
type User struct {
	gorm.Model        // Provides ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"unique;not null" json:"username" binding:"required"`
	Email      string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
}

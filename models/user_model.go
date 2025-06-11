package models

import "gorm.io/gorm"

// Name struct holds first and last names
type Name struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// User represents the user model
type User struct {
	gorm.Model
	Name     *Name  `gorm:"embedded;embeddedPrefix:name_" json:"name" binding:"required"`
	Username string `gorm:"unique;not null" json:"username" binding:"required"`
	Email    string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

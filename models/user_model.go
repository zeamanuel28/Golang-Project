package models

type Name struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     *Name  `gorm:"embedded;embeddedPrefix:name_" json:"name" binding:"required"`
	Username string `gorm:"unique;not null" json:"username" binding:"required"`
	Email    string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Role     string `gorm:"type:varchar(20);default:'user'" json:"role"`
	Password string `json:"password" binding:"required,min=6"`

	// One-to-Many relationship with Book
	Books []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"books"`
}

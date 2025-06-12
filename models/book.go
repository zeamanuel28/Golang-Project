package models

type Book struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	UserID uint   `json:"user_id"` // Foreign key

	// Establishing the relationship to User with proper cascading
	//User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

package models

type Role struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex" validate:"required"`
	Description string `json:"description" validate:"required"`
}

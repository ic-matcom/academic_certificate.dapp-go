package models

type User struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Username   string `json:"username" gorm:"uniqueIndex" validate:"required"`
	Passphrase string `json:"passphrase" validate:"required"`
	FirstName  string `json:"firstname" validate:"required"`
	LastName   string `json:"lastname" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Role        string `json:"rol_id" validate:"required"`
}

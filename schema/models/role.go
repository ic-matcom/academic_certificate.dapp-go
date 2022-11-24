package models

type Role struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"uniqueIndex" validate:"required"`
	Description string `json:"description" validate:"required"`
}

const (
	Role_Invalid          = "Inválido"
	Role_SystemAdmin      = "Administrador de Sistema"
	Role_CertificateAdmin = "Administrador de Certificados"
	Role_Secretary        = "Secretario"
	Role_Dean             = "Decano"
	Role_Rector           = "Rector"
)

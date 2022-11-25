package models

type Role struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Label       string `json:"label" gorm:"uniqueIndex" validate:"required"`
	Name        string `json:"name"  validate:"required"`
	Description string `json:"description" validate:"required"`
}

const (
	Role_Invalid          = "invalid"
	Role_SystemAdmin      = "sysadmin"
	Role_CertificateAdmin = "certadmin"
	Role_Secretary        = "secretary"
	Role_Dean             = "dean"
	Role_Rector           = "rector"
)

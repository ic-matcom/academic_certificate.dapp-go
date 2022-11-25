package repo

import (
	"dapp/lib"
	"dapp/schema/dto"
	"dapp/schema/models"
	"dapp/service/utils"
	"fmt"
	"log"
	"math"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// region ======== SETUP =================================================================

type RepoUser struct {
	DBLocation string
	DB         *gorm.DB
}

var singletonRU *RepoUser

// using Go sync package to invoke a method exactly only once
var onceRU sync.Once

// endregion =============================================================================

func NewRepoUser(svcConf *utils.SvcConfig) *RepoUser {
	onceRU.Do(func() {
		singletonRU = &RepoUser{DBLocation: svcConf.StoreDBPath}
		singletonRU.InitDB()
		singletonRU.PopulateDB()
	})
	return singletonRU
}

// GetUser get the user from the DB
func (r *RepoUser) GetUser(userID int) (models.User, error) {
	var modelUser models.User
	if result := r.DB.First(&modelUser, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	return modelUser, nil
}

// GetUserByUsername GetUser get the user from the DB
func (r *RepoUser) GetUserByUsername(username string) (models.User, error) {
	var modelUser models.User
	if result := r.DB.First(&modelUser, models.User{Username: username}); result.Error != nil {
		return models.User{}, result.Error
	}
	return modelUser, nil
}

// GetUsers return a list of dto.User
func (r *RepoUser) GetUsers(pagination *dto.Pagination) (*dto.Pagination, error) {
	var users []models.User
	result := r.DB.Scopes(paginate(users, pagination, r.DB)).Find(&users)
	if result.Error != nil {
		return &dto.Pagination{}, result.Error
	}
	pagination.Rows = users
	return pagination, nil
}

// AddUser Add the user to database
// Returns nil if user was added correctly, otherwise return error found
func (r *RepoUser) AddUser(user models.User) (models.User, error) {
	result := r.DB.Create(&user)
	return user, result.Error
}

// UpdateUser Update user with id UserID to new data in database
// Returns nil if user was updated correctly, otherwise return error found
func (r *RepoUser) UpdateUser(userID int, user models.User) (models.User, error) {
	var userInDB models.User
	if result := r.DB.First(&userInDB, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	result := r.DB.Save(&user)
	return user, result.Error
}

// RemoveUser Remove user from database
// Returns nil if user was removed correctly, otherwise return error found
func (r *RepoUser) RemoveUser(userID int) (models.User, error) {
	var modelUser models.User
	if result := r.DB.First(&modelUser, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	result := r.DB.Delete(&modelUser)
	return modelUser, result.Error
}

func (r *RepoUser) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	result := r.DB.Find(&roles)
	return roles, result.Error
}

// InvalidateUser Invalidate user, remove all access privileges
func (r *RepoUser) InvalidateUser(userID int) (models.User, error) {
	var modelUser models.User
	if result := r.DB.First(&modelUser, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	modelUser.Role = models.Role_Invalid
	result := r.DB.Save(&modelUser)
	return modelUser, result.Error
}

func (r *RepoUser) InitDB() {
	// TODO: move dbURL to configuration files conf.sample.unix, conf.sample.windows and conf.yaml
	dbURL := "postgres://pg:pass@localhost:5432/users"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&models.User{}, &models.Role{})
	r.DB = db
}
func (r *RepoUser) PopulateDB() {
	r.PopulateUserTable()
	r.PopulateRolTable()
}

func (r *RepoUser) PopulateUserTable() {
	var usersInDB []models.User
	if result := r.DB.Find(&usersInDB); result.Error != nil {
		fmt.Println(result.Error)
	}
	if len(usersInDB) != 0 {
		return
	}

	p1, _ := lib.Checksum("SHA256", []byte("password1"))
	users := []models.User{
		{
			Username:   "richard",
			Passphrase: p1,
			FirstName:  "Richard",
			LastName:   "Sargon",
			Email:      "richard.sargon@meinermail.com",
			Role:       models.Role_Secretary,
		},
		{
			Username:   "tom",
			Passphrase: p1,
			FirstName:  "Tom",
			LastName:   "Carter",
			Email:      "tom.carter@meinermail.com",
			Role:       models.Role_Rector,
		},
		{
			Username:   "Ariel",
			Passphrase: p1,
			FirstName:  "Ariel",
			LastName:   "Huerta",
			Email:      "ariel@gmail.com",
			Role:       models.Role_SystemAdmin,
		},
		{
			Username:   "ALab",
			Passphrase: p1,
			FirstName:  "Alejandro",
			LastName:   "Labourdette",
			Email:      "alab@gmail.com",
			Role:       models.Role_Dean,
		},
	}
	for i := 0; i < 100; i++ {
		users = append(users, models.User{
			Username:   fmt.Sprintf("UserName%d", i),
			Passphrase: p1,
			FirstName:  fmt.Sprintf("Name%d", i),
			LastName:   fmt.Sprintf("Last%d", i),
			Email:      fmt.Sprintf("bot%d@gmail.com", i),
			Role:       models.Role_CertificateAdmin,
		})
	}
	for _, user := range users {
		if result := r.DB.Create(&user); result.Error != nil {
			fmt.Println(result.Error)
		}
	}
}

func (r *RepoUser) PopulateRolTable() {
	var rolInDB []models.Role
	if result := r.DB.Find(&rolInDB); result.Error != nil {
		fmt.Println(result.Error)
	}
	if len(rolInDB) != 0 {
		return
	}
	r.DB.Create(&models.Role{
		Label:       models.Role_Invalid,
		Name:        "Usuario Invalidado",
		Description: "Usuario que le fueron quitados sus privilegios.",
	})
	r.DB.Create(&models.Role{
		Label:       models.Role_SystemAdmin,
		Name:        "Administrador de Sistemas",
		Description: "Usuario que puede gestionar los usuarios de la dapp.",
	})
	r.DB.Create(&models.Role{
		Label:       models.Role_CertificateAdmin,
		Name:        "Administrador de Certificados",
		Description: "Usuario que puede gestionar los certificados almacenados.",
	})
	r.DB.Create(&models.Role{
		Label:       models.Role_Secretary,
		Name:        "Secretario General",
		Description: "Usuario que valida los certificados emitidos.",
	})
	r.DB.Create(&models.Role{
		Label:       models.Role_Dean,
		Name:        "Decano de Facultad",
		Description: "Usuario que valida los certificados emitidos.",
	})
	r.DB.Create(&models.Role{
		Label:       models.Role_Rector,
		Name:        "Rector de Universidad",
		Description: "Usuario que valida los certificados emitidos.",
	})
}

func paginate(value interface{}, pagination *dto.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}

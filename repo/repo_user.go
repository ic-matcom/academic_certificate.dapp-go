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
}

var singletonRU *RepoUser

// using Go sync package to invoke a method exactly only once
var onceRU sync.Once

// endregion =============================================================================

func NewRepoUser(svcConf *utils.SvcConfig) *RepoUser {
	onceRU.Do(func() {
		singletonRU = &RepoUser{DBLocation: svcConf.StoreDBPath}
		InitDB()
		PopulateDB()
	})
	return singletonRU
}

// DB to treat users persistence
var UsersDB *gorm.DB

// GetUser get the user from the DB
func (r *RepoUser) GetUser(userID int) (models.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	return modelUser, nil
}

// GetUser get the user from the DB
func (r *RepoUser) GetUserByUsername(username string) (models.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, models.User{Username: username}); result.Error != nil {
		return models.User{}, result.Error
	}
	return modelUser, nil
}

// GetUsers return a list of dto.User
func (r *RepoUser) GetUsers(pagination *dto.Pagination) (*dto.Pagination, error) {
	var users []models.User
	result := UsersDB.Scopes(paginate(users, pagination, UsersDB)).Find(&users)
	if result.Error != nil {
		return &dto.Pagination{}, result.Error
	}
	pagination.Rows = users
	return pagination, nil
}

// AddUser Add the user to database
// Returns nil if user was added correctly, otherwise return error found
func (r *RepoUser) AddUser(user models.User) (models.User, error) {
	result := UsersDB.Create(&user)
	return user, result.Error
}

// UpdateUser Update user with id UserID to new data in database
// Returns nil if user was updated correctly, otherwise return error found
func (r *RepoUser) UpdateUser(userID int, user models.User) (models.User, error) {
	var userInDB models.User
	if result := UsersDB.First(&userInDB, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	result := UsersDB.Save(&user)
	return user, result.Error
}

// RemoveUser Remove user from database
// Returns nil if user was removed correctly, otherwise return error found
func (r *RepoUser) RemoveUser(userID int) (models.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, userID); result.Error != nil {
		return models.User{}, result.Error
	}
	result := UsersDB.Delete(&modelUser)
	return modelUser, result.Error
}

func InitDB() {
	dbURL := "postgres://pg:pass@localhost:5432/users"
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&models.User{})
	UsersDB = db
}

func PopulateDB() {
	var usersInDB []models.User
	if result := UsersDB.Find(&usersInDB); result.Error != nil {
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
		},
		{
			Username:   "tom",
			Passphrase: p1,
			FirstName:  "Tom",
			LastName:   "Carter",
			Email:      "tom.carter@meinermail.com",
		},
		{
			Username:   "Ariel",
			Passphrase: p1,
			FirstName:  "Ariel",
			LastName:   "Huerta",
			Email:      "ariel@gmail.com",
		},
		{
			Username:   "ALab",
			Passphrase: p1,
			FirstName:  "Alejandro",
			LastName:   "Labourdette",
			Email:      "alab@gmail.com",
		},
	}
	for i := 0; i < 100; i++ {
		users = append(users, models.User{
			Username:   fmt.Sprintf("UserName%d", i),
			Passphrase: p1,
			FirstName:  fmt.Sprintf("Name%d", i),
			LastName:   fmt.Sprintf("Last%d", i),
			Email:      fmt.Sprintf("bot%d@gmail.com", i),
		})
	}
	for _, user := range users {
		if result := UsersDB.Create(&user); result.Error != nil {
			fmt.Println(result.Error)
		}
	}
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

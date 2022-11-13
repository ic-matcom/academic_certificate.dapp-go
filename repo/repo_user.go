package repo

import (
	"dapp/lib"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"dapp/schema/models"
	"dapp/service/utils"
	"fmt"
	"log"
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
func (r *RepoUser) GetUser(userID int) (dto.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, userID); result.Error != nil {
		return dto.User{}, result.Error
	}
	return mapper.MapModelUser2DtoUser(modelUser), nil
}

func (r *RepoUser) GetUserByUserName(username string) (dto.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, models.User{Username: username}); result.Error != nil {
		return dto.User{}, result.Error
	}
	return mapper.MapModelUser2DtoUser(modelUser), nil
}

// GetUsers return a list of dto.User
func (r *RepoUser) GetUsers() ([]dto.User, error) {
	var users []models.User
	if result := UsersDB.Find(&users); result.Error != nil {
		return []dto.User{}, result.Error
	}
	dtoUsers := []dto.User{}
	for _, u := range users {
		dtoUsers = append(dtoUsers, mapper.MapModelUser2DtoUser(u))
	}
	return dtoUsers, nil
}

// AddUser Add the user to database
// Returns nil if user was added correctly, otherwise return error found
func (r *RepoUser) AddUser(user dto.User) (dto.User, error) {
	modelUser := mapper.MapDtoUser2ModelUser(user)
	result := UsersDB.Create(&modelUser)
	return mapper.MapModelUser2DtoUser(modelUser), result.Error
}

// UpdateUser Update user with id UserID to new data in database
// Returns nil if user was updated correctly, otherwise return error found
func (r *RepoUser) UpdateUser(userID int, user dto.UserUpdate) (dto.User, error) {
	modelUser := mapper.MapUserUpd2ModelUser(userID, user)
	var userInDB models.User
	if result := UsersDB.First(&userInDB, userID); result.Error != nil {
		return dto.User{}, result.Error
	}
	result := UsersDB.Save(&modelUser)
	return mapper.MapModelUser2DtoUser(modelUser), result.Error
}

// RemoveUser Remove user from database
// Returns nil if user was removed correctly, otherwise return error found
func (r *RepoUser) RemoveUser(userID int) (dto.User, error) {
	var modelUser models.User
	if result := UsersDB.First(&modelUser, userID); result.Error != nil {
		return dto.User{}, result.Error
	}
	result := UsersDB.Delete(&modelUser)
	return mapper.MapModelUser2DtoUser(modelUser), result.Error
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
			Username:   "ALab",
			Passphrase: p1,
			FirstName:  "Alejandro",
			LastName:   "Labourdette",
			Email:      "alab@gmail.com",
		},
		{
			Username:   "Ariel",
			Passphrase: p1,
			FirstName:  "Ariel",
			LastName:   "Huerta",
			Email:      "ariel@gmail.com",
		},
	}
	for _, user := range users {
		if result := UsersDB.Create(&user); result.Error != nil {
			fmt.Println(result.Error)
		}
	}
}

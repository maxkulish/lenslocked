package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"lenslocked/database"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database
	ErrNotFound = errors.New("models: resource not found")
)

type User struct {
	Name  string
	Email string `gorm:"not null;unique_index"`
	gorm.Model
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService() (*UserService, error) {
	db, err := database.NewDBConn()
	if err != nil {
		fmt.Println("UserService error")
	}

	return &UserService{DB: db.Conn}, nil
}

// ByID will look up user by the id
// 1 - user, nil
// 2 - nil, errorNotFound
// 3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.DB.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (us *UserService) FullReset() {
	us.DB.DropTableIfExists(&User{})
	us.DB.AutoMigrate(&User{})
}

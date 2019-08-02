package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked/database"
)

const userPWPepper = "k$cUXbp!WY&vfGyhY64#UdeGesqz"

type User struct {
	gorm.Model
	Name         string `gorm:"index:user_name"`
	Email        string `gorm:"not null;unique_index:user_email"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService(env string) (*UserService, error) {
	db, err := database.NewDBConn(env)
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
	db := us.DB.Where("id = ?", id).First(&user)

	err := database.HandleDBError(db)
	return &user, err

}

// ByEmail looks up a user with the given email address and
// returns that user
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.DB.Where("email = ?", email).First(&user)

	err := database.HandleDBError(db)
	return &user, err
}

// Create will create the provided user
func (us *UserService) Create(user *User) error {
	// Add pepper to the user password
	pwBytes := []byte(user.Password + userPWPepper)
	// Hashing without validation
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	// Clear raw password from log, memory etc.
	user.Password = ""
	return us.DB.Create(user).Error
}

// Update will update the provided user with all of the data
func (us *UserService) Update(user *User) error {
	return us.DB.Save(user).Error
}

func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return database.ErrInvalidID
	}
	user := User{
		Model: gorm.Model{ID: id},
	}
	return us.DB.Delete(&user).Error
}

func (us *UserService) FullReset() error {
	if err := us.DB.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.DB.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}

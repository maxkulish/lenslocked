package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked/database"
	"lenslocked/hash"
	"lenslocked/rand"
)

const (
	userPWPepper  = "k$cUXbp!WY&vfGyhY64#UdeGesqz"
	hmacSecretKey = "ujY4n%wnUBD#cAyQh4VXqJk*imr"
)

type User struct {
	gorm.Model
	Name         string `gorm:"index:user_name"`
	Email        string `gorm:"not null;unique_index:user_email"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type UserService struct {
	DB   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(env string) (*UserService, error) {
	db, err := database.NewDBConn(env)
	if err != nil {
		fmt.Println("UserService error")
	}

	hmac := hash.NewHMAC(hmacSecretKey)

	return &UserService{
		DB:   db.Conn,
		hmac: hmac,
	}, nil
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

// ByRemember looks up a user with the given remember token
// and returns that user. This method will handle
// hashing the token for us
// Errors the same as ByEmail
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User

	hashedToken := us.hmac.Hash(token)
	err := us.DB.Where("remember_hash = ?", hashedToken).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPWPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, database.ErrInvalidPass
		default:
			return nil, err
		}
	}

	return foundUser, nil
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

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	user.RememberHash = us.hmac.Hash(user.Remember)

	return us.DB.Create(user).Error
}

// Update will update the provided user with all of the data
func (us *UserService) Update(user *User) error {

	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}

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

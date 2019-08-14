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

// UserDB is used to interact with the users database
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Used to close a DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	FullReset() error
}

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
	UserDB
}

type UserValidator struct {
	UserDB
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

func NewUserService(env string) (*UserService, error) {
	ug, err := NewUserGorm(env)
	if err != nil {
		return nil, err
	}

	return &UserService{
		ug,
	}, nil
}

func NewUserGorm(env string) (*userGorm, error) {
	db, err := database.NewDBConn(env)
	if err != nil {
		fmt.Println("UserService error")
	}

	hmac := hash.NewHMAC(hmacSecretKey)

	return &userGorm{
		db:   db.Conn,
		hmac: hmac,
	}, nil
}

// ByID will look up user by the id
// 1 - user, nil
// 2 - nil, errorNotFound
// 3 - nil, otherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id).First(&user)

	err := database.HandleDBError(db)
	return &user, err

}

// ByEmail looks up a user with the given email address and
// returns that user
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email).First(&user)

	err := database.HandleDBError(db)
	return &user, err
}

// ByRemember looks up a user with the given remember token
// and returns that user. This method will handle
// hashing the token for us
// Errors the same as ByEmail
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User

	hashedToken := ug.hmac.Hash(token)
	err := ug.db.Where("remember_hash = ?", hashedToken).First(&user).Error

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
func (ug *userGorm) Create(user *User) error {
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

	user.RememberHash = ug.hmac.Hash(user.Remember)

	return ug.db.Create(user).Error
}

// Update will update the provided user with all of the data
func (ug *userGorm) Update(user *User) error {

	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}

	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return database.ErrInvalidID
	}
	user := User{
		Model: gorm.Model{ID: id},
	}
	return ug.db.Delete(&user).Error
}

func (ug *userGorm) FullReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}

	return nil
}

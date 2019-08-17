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

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address
	// and password are correct
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(env string) (UserService, error) {
	ug, err := NewUserGorm(env)
	if err != nil {
		return nil, err
	}

	hmac := hash.NewHMAC(hmacSecretKey)

	uv := &userValidator{
		UserDB: ug,
		hmac:   hmac,
	}

	return &userService{
		uv,
	}, nil
}

func (us *userService) Authenticate(email, password string) (*User, error) {
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

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

// ByRemember will hash the remember token and then call
// ByRemember on the subsequent UserDB layer
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *userValidator) Create(user *User) error {

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	err := runUserValFuncs(user,
		uv.bcryptPass,
		uv.setRememberIfUnset,
		uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// bcryptPass will hash a user's password with a
// predefined pepper (userPwPepper) and bcypt if the
// Password field is not the empty string
func (uv *userValidator) bcryptPass(user *User) error {
	if user.Password == "" {
		return nil
	}

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

	return nil
}

func (uv *userValidator) hmacRemember(user *User) error {

	if user.Remember == "" {
		return nil
	}

	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember == "" {
		return nil
	}

	token, err := rand.RememberToken()
	if err != nil {
		return err
	}

	user.Remember = token
	return nil
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {

	err := runUserValFuncs(user, uv.bcryptPass, uv.hmacRemember)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	if id == 0 {
		return database.ErrInvalidID
	}
	return uv.UserDB.Delete(id)
}

type userService struct {
	UserDB
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

func NewUserGorm(env string) (*userGorm, error) {
	db, err := database.NewDBConn(env)
	if err != nil {
		fmt.Println("UserService error")
	}

	return &userGorm{
		db: db.Conn,
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
func (ug *userGorm) ByRemember(hashedToken string) (*User, error) {
	var user User

	err := ug.db.Where("remember_hash = ?", hashedToken).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create will create the provided user
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update the provided user with all of the data
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

func (ug *userGorm) Delete(id uint) error {
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

func (ug *userGorm) Close() error {

	if err := ug.db.Close(); err != nil {
		return err
	}

	return nil
}

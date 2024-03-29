package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"lenslocked/database"
	"lenslocked/hash"
	"lenslocked/rand"
	"regexp"
	"strings"
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
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address
	// and password are correct
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB) UserService {

	ug := &userGorm{db}

	hmac := hash.NewHMAC(hmacSecretKey)

	uv := newUserValidator(ug, hmac)

	return &userService{
		uv,
	}
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
			return nil, ErrPasswordIncorrect
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

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB:     udb,
		hmac:       hmac,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}

	err := runUserValFuncs(&user, uv.emailNormalizer)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
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
		uv.passwordRequired,
		uv.passwordMinLenght,
		uv.bcryptPass,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.emailNormalizer,
		uv.emailRequired,
		uv.emailFormat,
		uv.emailIsAvail)
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

func (uv *userValidator) passwordMinLenght(user *User) error {
	if user.Password == "" {
		return nil
	}

	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}

	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}

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

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}

	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}

	if n < 32 {
		return ErrRememberTooShort
	}

	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}

	return nil
}

// Update will hash a remember token if it is provided
func (uv *userValidator) Update(user *User) error {

	err := runUserValFuncs(user,
		uv.passwordMinLenght,
		uv.bcryptPass,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.emailNormalizer,
		uv.emailRequired,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id

	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrInvalidID
		}
		return nil
	})
}

func (uv *userValidator) emailNormalizer(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) emailRequired(user *User) error {
	if user.Email == "" {
		return ErrEmailRequited
	}

	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}

	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	exist, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}

	if err != nil {
		return err
	}

	// We found a user with this email address
	// If the found user has the same ID as this user, it is
	// an update and this is the same user
	if user.ID != exist.ID {
		return ErrEmailTaken
	}

	return nil
}

type userService struct {
	UserDB
}

type userGorm struct {
	db *gorm.DB
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

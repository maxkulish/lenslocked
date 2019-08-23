package models

import "strings"

var (
	// ErrEmailRequited is returned when an email address is not provided
	ErrEmailRequited modelError = "email address is required"

	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid modelError = "email address is not valid"

	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user password that is less than 8 characters
	ErrPasswordTooShort modelError = "password must be at least 8 characters"

	// ErrPasswordRequired is returned when an user password field is not provided
	ErrPasswordRequired modelError = "password is required"

	// ErrRememberTooShort is returned when a remember token is
	// not at least 32 bytes
	ErrRememberTooShort modelError = "remember token must be at least 32 bytes"

	// ErrRememberRequired
	ErrRememberRequired modelError = "remember hash is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
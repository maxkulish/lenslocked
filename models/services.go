package models

import (
	"fmt"
	"lenslocked/database"
)

func NewServices(env string) (*Services, error) {
	db, err := database.NewDBConn(env)
	if err != nil {
		fmt.Println("UserService error")
	}
	db.Conn.LogMode(true)
	return &Services{}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
}

package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"lenslocked/database"
)

func NewServices(env string) (*Services, error) {
	db, err := database.NewDBConn(env)
	if err != nil {
		fmt.Println("UserService error")
	}
	db.Conn.LogMode(true)
	return &Services{
		User:    NewUserService(db.Conn),
		Gallery: NewGalleryService(db.Conn),
		db:      db.Conn,
	}, nil
}

type Services struct {
	Gallery GalleryService
	User    UserService
	db      *gorm.DB
}

// FullReset drops all tables and rebuilds them
func (s *Services) FullReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the all tables
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

// Closes the database connection
func (s *Services) Close() error {
	return s.db.Close()
}

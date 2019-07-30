package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	Name  string
	Email string `gorm:"not null;unique_index"`
	gorm.Model
}

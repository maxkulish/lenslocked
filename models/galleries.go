package models

import "github.com/jinzhu/gorm"

// Gallery is our image container resources that
// visitor view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface{}

type GalleryDB interface {
	Create(gallery *Gallery)
}

package models

import (
	"fmt"
	"io"
	"os"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	makeImagePath(galleryID uint) (string, error)
	//ByGalleryID(galleryID uint) []string
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	r.Close()
	path, err := is.makeImagePath(galleryID)
	if err != nil {
		return err
	}

	// Create a destination file
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}

func (is *imageService) makeImagePath(galleryID uint) (string, error) {
	// Create directory using gallery ID
	galleryPath := fmt.Sprintf("images/galleries/%v/", galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	return galleryPath, nil
}

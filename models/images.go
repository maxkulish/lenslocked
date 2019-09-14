package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Image is NOT stored in the database
type Image struct {
	GalleryID uint
	Filename  string
}

func (i *Image) Path() string {
	return "/" + i.RelativePath()
}

func (i *Image) RelativePath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	makeImagePath(galleryID uint) (string, error)
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
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
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}

	return galleryPath, nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)

	imgStr, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}

	ret := make([]Image, len(imgStr))
	for i := range imgStr {
		imgStr[i] = strings.Replace(imgStr[i], path, "", 1)
		ret[i] = Image{
			GalleryID: galleryID,
			Filename:  imgStr[i],
		}
	}

	return ret, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) Delete(image *Image) error {
	return os.Remove(image.RelativePath())
}

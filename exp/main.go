package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"lenslocked/models"
)

type DB struct {
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type User struct {
	gorm.Model
	Name   string
	Email  string `gorm:"not null;unique_index"`
	Color  string
	Orders []Order
}

type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func (d *DB) loadConf(path string) *DB {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, d)
	if err != nil {
		panic(err)
	}

	return d
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	err := db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	}).Error
	if err != nil {
		panic(err)
	}
}

func main() {

	us, err := models.NewUserService()
	if err != nil {
		panic(err)
	}
	defer us.DB.Close()

	//us.FullReset()
	us.DB.AutoMigrate(&User{})

	user, err := us.ByID(2)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}

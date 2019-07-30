package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	var d DB
	d.loadConf("/Users/mk/Code/lenslocked/config.yaml")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Database)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	db.LogMode(true)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Order{})

	var u User
	if err := db.Preload("Orders").First(&u, "id = ?", 4).Error; err != nil {
		panic(err)
	}
	//createOrder(db, u, 101, "Fake description #1")
	//createOrder(db, u, 9999, "Fake description #2")
	//createOrder(db, u, 6666, "Fake description #3")

	if db.RecordNotFound() {
		fmt.Println("no user found!")
	} else if db.Error != nil {
		panic(db.Error)
	}
	fmt.Println(u)

}

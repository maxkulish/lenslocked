package main

import (
	"fmt"
	"lenslocked/models"
)

func main() {

	us, err := models.NewUserService("dev")
	if err != nil {
		panic(err)
	}
	defer us.DB.Close()
	us.AutoMigrate()

	user := models.User{
		Name:     "Max Kul",
		Email:    "max@g.com",
		Password: "1234",
		Remember: "abc1234",
	}

	err = us.Create(&user)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user)

	user2, err := us.ByRemember("abc1234")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", *user2)
}

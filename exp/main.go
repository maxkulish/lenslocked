package main

import (
	"html/template"
	"log"
	"os"
)

type User struct {
	Name string
	Dog
	Int   int
	Float float64
	Slice []string
	Map   map[string]string
}

type Dog struct {
	Name string
}

func main() {
	t, err := template.ParseFiles("./exp/hello.gohtml")
	if err != nil {
		log.Fatal("Can't open gohtml template file", err)
	}

	data := User{
		Name:  "John Smith",
		Int:   15,
		Float: 67.0,
		Slice: []string{"nestor", "igor"},
		Map: map[string]string{
			"key1": "some key 1",
			"key2": "some key 2",
		},
	}

	data.Dog.Name = "Jessie"

	err = t.Execute(os.Stdout, data)
	if err != nil {
		log.Fatal("Can't execute template. Error: ", err)
	}
}

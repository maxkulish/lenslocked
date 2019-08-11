package main

import (
	"fmt"
	"lenslocked/rand"
)

func main() {
	fmt.Println(rand.String(10))
	fmt.Println(rand.RememberToken())
}

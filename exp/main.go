package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"lenslocked/hash"
)

func main() {
	toHash := []byte("This is my string to hash")
	h := hmac.New(sha256.New, []byte("my-secret-key"))
	h.Write(toHash)
	b := h.Sum(nil)
	fmt.Println(base64.URLEncoding.EncodeToString(b))
	h.Reset()

	myhmac := hash.NewHMAC("my-secret-key")
	fmt.Println(myhmac.Hash("This is my string to hash"))
}

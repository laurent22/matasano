package main

import (
	"log"
	"./cryptoutil"
)

func main() {
	s := []byte("YELLOW SUBMARINE")
	s = cryptoutil.Pkcs7padding(s, 20)
	log.Println(string(s))
}
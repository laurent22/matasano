package main

import (
	"log"
	"io/ioutil"
	"encoding/base64"
	"./cryptoutil"
)

func main() {
	key := []byte("YELLOW SUBMARINE")
	iv := []byte{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0}
	content, _ := ioutil.ReadFile("q10_data.txt")
	data, _ := base64.StdEncoding.DecodeString(string(content))
	dec := cryptoutil.AES128CBCDecrypt(data, key, iv)
	log.Println(string(dec))
}
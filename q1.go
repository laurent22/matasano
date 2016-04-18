package main

import (
	"encoding/hex"
	"encoding/base64"
	"fmt"
)

func HexToBase64(s string) string {
	bytes, _ := hex.DecodeString(s)
	output := base64.StdEncoding.EncodeToString(bytes)
	return output
}

func main() {
	r := HexToBase64("49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d")
	fmt.Println(r)
}
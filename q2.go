package main

import (
	"encoding/hex"
	"fmt"
)

func XorHexStrings(s1 string, s2 string) string {
	if len(s1) != len(s2) {
		panic("Strings must have the same length")
	}
	
	// Decode the two strings into byte slices
	bytes1, _ := hex.DecodeString(s1)
	bytes2, _ := hex.DecodeString(s2)
	var output []byte
	for i := 0; i < len(bytes1); i++ {
		// XOR each byte to byte and build the encrypted string
		b := bytes1[i] ^ bytes2[i]
		output = append(output, b)
	}
	
	return hex.EncodeToString(output)
}

func main() {
	r := XorHexStrings("1c0111001f010100061a024b53535009181c", "686974207468652062756c6c277320657965")
	fmt.Println(r)
}
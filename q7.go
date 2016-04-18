package main

import (
	"crypto/aes"
	"fmt"
	"io/ioutil"
	"encoding/base64"
)

func AES128ECBDecrypt(encrypted []byte, key []byte) []byte {
	cypher, _ := aes.NewCipher(key)
	blockSize := 16
	output := make([]byte, len(encrypted))
	for i := 0; i < len(output); i += blockSize {
		cypher.Decrypt(output[i:i+blockSize], encrypted[i:i+blockSize]) 
	}
	return output
}

func main() {
	// Use Go built-in library to do the AES-128-ECB decryption.
	content, _ := ioutil.ReadFile("q7_data.txt")
	data, _ := base64.StdEncoding.DecodeString(string(content))
	key := []byte("YELLOW SUBMARINE")
	decrypted := AES128ECBDecrypt(data, key)
	fmt.Println(string(decrypted))
}
package main

import (
	"log"
	"math/rand"	
	"./cryptoutil"
	"time"
)

func prependBytes(dest []byte, source []byte) []byte {
	var output []byte
	output = appendBytes(output, source)
	output = appendBytes(output, dest)
	return output
}

func appendBytes(dest []byte, source []byte) []byte {
	output := dest
	for _, b := range source {
		output = append(output, b)
	}
	return output
}

// Returns the randomly encrypted bytes along with the 
// mode that was used (for testing purposes)
func encryption_oracle(ptext []byte) ([]byte, string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Prepend 5-10 random bytes
	ptext = prependBytes(ptext, cryptoutil.RandomBytes(5 + r.Intn(6)))
	// Append 5-10 random bytes
	ptext = appendBytes(ptext, cryptoutil.RandomBytes(5 + r.Intn(6)))
	// Pad the data to 16 bytes so that it can be encrypted
	ptext = cryptoutil.Pkcs7padding(ptext, len(ptext) + cryptoutil.Pkcs7paddingCount(ptext))
	// Create the random key
	key := cryptoutil.RandomBytes(16)
	if r.Intn(2) == 0 {
		return cryptoutil.AES128ECBEncrypt(ptext, key), "ecb"
	} else {
		iv := cryptoutil.RandomBytes(16)
		return cryptoutil.AES128CBCEncrypt(ptext, key, iv), "cbc"
	}
}

func main() {
	// Use a string with the same characters, long enough so that we
	// can produce a repeating block.
	input := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	for i := 0; i < 100; i++ {
		ciphertext, mode := encryption_oracle([]byte(input))
		var detectedMode string
		if cryptoutil.IsECBEncrypted(ciphertext) {
			detectedMode = "ecb"
		} else {
			detectedMode = "cbc"
		}
		log.Printf("Detected: %s. Real: %s. OK: %t", detectedMode, mode, detectedMode == mode)
	}
}
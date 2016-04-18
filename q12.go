package main

import (
	"log"
	"./cryptoutil"
	"encoding/base64"
)

func appendBytes(dest []byte, source []byte) []byte {
	output := dest
	for _, b := range source {
		output = append(output, b)
	}
	return output
}

func fillBytes(b byte, length int) []byte {
	var output []byte
	for len(output) < length {
		output = append(output, b)
	}
	return output
} 

var encryption_oracle_key []byte

func encryption_oracle(plaintext []byte) []byte {
	secret, _ := base64.StdEncoding.DecodeString("Um9sbGluJyBpbiBteSA1LjAKV2l0aCBteSByYWctdG9wIGRvd24gc28gbXkgaGFpciBjYW4gYmxvdwpUaGUgZ2lybGllcyBvbiBzdGFuZGJ5IHdhdmluZyBqdXN0IHRvIHNheSBoaQpEaWQgeW91IHN0b3A/IE5vLCBJIGp1c3QgZHJvdmUgYnkK")
	plaintext = appendBytes(plaintext, secret)
	plaintext = cryptoutil.Pkcs7padding(plaintext, len(plaintext) + cryptoutil.Pkcs7paddingCount(plaintext))
	return cryptoutil.AES128ECBEncrypt(plaintext, encryption_oracle_key)
}

func main() {
	var ciphertext []byte
	encryption_oracle_key = cryptoutil.RandomBytes(16)
	
	// Find the block size by feeding the oracle function an input of
	// increasing length. The block size is revealed when
	// the cypher text length increases between two iteration.
	
	blockSize := 0
	previousLength := 0
	s := ""
	for i := 0; true; i++ {
		ciphertext = encryption_oracle([]byte(s))
		currentLength := len(ciphertext)
		if previousLength != 0 {
			if previousLength != currentLength {
				blockSize = currentLength - previousLength
				break
			}
		}
		s += "a"
		previousLength = currentLength
	}
	
	log.Println("Block size:", blockSize)
	
	// Detect that the function is using ECB
	
	ciphertext = encryption_oracle(fillBytes(222, 32))
	isECB := cryptoutil.IsECBEncrypted(ciphertext)
	
	log.Println("Is ECB:", isECB)
	
	// Decrypt the data
	
	var plaintext []byte
	// Start with an input string of (blockSize - 1) length
	input := fillBytes(1, blockSize - 1)
	blockIndex := 0
	// For each block of the cypher text
	for blockIndex < len(ciphertext) / blockSize {
		// Get the cypher text for our current string
		ciphertext = encryption_oracle(input)
		// Then discover the missing byte by looping through them and feeding our guess in the
		// oracle function.
		for i := 0; i < 256; i++ {
			var guess []byte
			// if we haven't decrypted a complete block yet, the input is something like below:
			// Assuming a block size of 4, "#" is what we're looking for, "XXX" is the encrypted string.
			// 111#XXXXXXXX
			// 11A#XXXXXXX
			// 1AB#XXXXXX
			// ABC#XXXX
			if len(plaintext) < blockSize {
				guess = input
				guess = appendBytes(guess, plaintext)
			} else {
				// Now that we've encrypted a full block, we feed this into the oracle function, minus the
				// first byte, and we append the byte we are trying to guess.
				// BCD#XXXX
				// CDE#XXX
				// DEF#XX
				// EFG#
				guess = plaintext[len(plaintext) - (blockSize - 1):len(plaintext)]
			}
			guess = append(guess, byte(i))
			
			// Encrypt our guess
			b1 := encryption_oracle(guess)[0:blockSize]
			// Get the nth block
			b2 := ciphertext[blockSize * blockIndex:blockSize * blockIndex + blockSize]
			// If the blocks match, we found the right byte. Add it to the plain text.
			if cryptoutil.SliceEquals(b1, b2) {
				plaintext = append(plaintext, byte(i))
				break
			}
		}	
		
		// Once we processed a full block, we start over with the next one.
		if len(input) == 0 {
			input = fillBytes(1, blockSize)
			blockIndex++
		}
		
		input = input[0:len(input)-1]
	}
	
	log.Println("Plain text:")
	log.Println(string(plaintext))
}
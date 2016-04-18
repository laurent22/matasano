package main

import (
	"log"
	"math/rand"	
	"time"
	"./cryptoutil"
)

var randomKey []byte

func removePkcs7padding(data []byte) ([]byte, bool) {
	if len(data) == 0 {
		return []byte{}, false
	}
	paddingSize := int(data[len(data) - 1])
	if len(data) - paddingSize < 0 {
		return []byte{}, false
	}
	for i := len(data) - 1; i >= len(data) - paddingSize; i-- {
		if int(data[i]) != paddingSize {
			return []byte{}, false
		}
	}
	if paddingSize == 16 {
		return []byte{}, true
	}
	return data[0:len(data) - paddingSize], true
}

func encrypt(message []byte, iv []byte) []byte {
	message = cryptoutil.Pkcs7padding(message, len(message) + cryptoutil.Pkcs7paddingCount(message))
	return cryptoutil.AES128CBCEncrypt(message, randomKey, iv)
}

func decrypt(ciphertext []byte, iv []byte) ([]byte, bool) {
	output, ok := removePkcs7padding(cryptoutil.AES128CBCDecrypt(ciphertext, randomKey, iv))
	if !ok {
		return ciphertext, false
	}
	return output, true
}

// Decrypt one byte of the given block from the given ciphertext.
// Returns the decrypted byte and the input that was used to generate valid padding.
// In some cases, it's not possible to find a byte that would produce a valid padding, usually because
// the previously discovered plaintext is incorrect. In that case, we return `false`. It's then up to
// the caller (decryptRange) to backtrack.
func decryptByte(ciphertext []byte, bi int, byteIndex int, inputStart byte, plaintext []byte, iv []byte) (byte, byte, bool) {
	bs := 16
	
	// Block 2 is the block we want to decrypt.
	// Block 1 is the one we will modify.
	// "p" for "plaintext"
	// "c" for "ciphertext"
	// c'1 is the block we manually create.
	
	// # To decrypt byte 15:
	
	// We create c'1 so that p2 has valid padding. 
	// If the padding is valid, the last byte of p2 is 0x01

	// Therefore we have this:
	// c'1[15] ^ (c1[15] ^ p2[15]) = 0x01
	
	// Which can be simplified to this:
	// p2[15] = 0x01 ^ c'1[15] ^ c1[15]
	
	// # To decrypt byte 14:
	
	// Set second to last byte so that p2 has valid padding
	// Last byte of p2 is 0x02
	
	// We already know that:
	// 0x01 = c'1[15] ^ c1[15] ^ p2[15]
	// c'1[15] = 0x01 ^ c1[15] ^ p2[15]
	
	// So we can calculate:
	// p2[14] = 0x02 ^ c'1[14] ^ c1[14]
		
	// Create c1
	var c1 []byte
	if bi == 0 {
		c1 = iv
	} else {
		c1 = ciphertext[(bi - 1) * bs : (bi - 1) * bs + bs]	
	}
	
	// Create c2
	c2 := ciphertext[bi * bs : bi * bs + bs]
	
	for i := int(inputStart); i < 256; i++ {
		// Create c'1
		
		// Start with 0s
		c1m := cryptoutil.FillBytes(0, byteIndex)
		
		// Append our guess
		c1m = append(c1m, byte(i))
		
		// Append what we already know
		if len(plaintext) >= 1 {
			padding := byte(bs - byteIndex)
			var knownBytes []byte
			for j := 0; j < len(plaintext); j++ {
				// c'1[15] = 0x01 ^ c1[15] ^ p2[15]
				knownBytes = append(knownBytes, padding ^ c1[bs - len(plaintext) + j] ^ plaintext[j])
			}
			c1m = cryptoutil.AppendBytes(c1m, knownBytes)
		}
		// Put c1 and c2 together and check if the padding is valid
		newCiphertext := c1m
		newCiphertext = cryptoutil.AppendBytes(newCiphertext, c2)
		_, ok := decrypt(newCiphertext, iv)
		if ok {
			return byte(bs - byteIndex) ^ byte(i) ^ c1[byteIndex], byte(i), true
		}
	}
	
	// Couldn't not find any combination that would generate a valid padding.
	return 0, 0, false
}

// Recursively decrypt the given block from the given ciphertext.
// The basic algorithm is like this:
// - P = Decrypt byte n
// - R = Decrypt range 0 to (n-1)
// - Append P to R
// - Return the result
// There are then some additional checks in case a range cannot be decrypted. In that case,
// we try to get P again starting from where we previously stopped. If P still cannot be decrypted,
// we return false and the previous function will also try to find a new value for its own P, and so on, recursively.
func decryptRange(fromIndex int, ciphertext []byte, blockIndex int, plaintext []byte, iv []byte) ([]byte, bool) {
	var inputStart byte = 0
	for {
		// Decrypt byte n
		result, inputByte, ok := decryptByte(ciphertext, blockIndex, fromIndex, inputStart, plaintext, iv)
		if !ok {
			return plaintext, false
		}
		// If it's the last byte, just return it
		if fromIndex == 0 {
			return []byte{result}, true
		}
		// Generate the new plaintext
		temp := cryptoutil.PrependByte(plaintext, result)
		// Decrypt range 0 to n-1
		before, ok := decryptRange(fromIndex - 1, ciphertext, blockIndex, temp, iv)
		if !ok {
			// We've tried all possible combination - return false
			if inputByte == 255 {
				return plaintext, false
			}
			// Try to decrypt byte n again, starting from where we previously stopped
			inputStart = inputByte + 1
			continue
		}
		// Append byte n to range
		before = append(before, result)
		return before, true
	}
}

func main() {
	bs := 16 // block size
	randomKey = cryptoutil.RandomBytes(bs)
	
	randomStrings := []string{
		"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
		"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
		"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
		"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
		"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
		"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
		"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
		"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
		"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
		"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
	}
	
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := randomStrings[r.Intn(len(randomStrings))]
	iv := cryptoutil.RandomBytes(16)
	ciphertext := encrypt([]byte(randomString), iv)
	blockCount := len(ciphertext) / bs
	
	for blockIndex := 0; blockIndex < blockCount; blockIndex++ {
		var plaintext []byte
		t, _ := decryptRange(15, ciphertext, blockIndex, plaintext, iv)
		log.Println(string(t))
	}
}
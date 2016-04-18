package main

import (
	"log"
	"strings"
	"./cryptoutil"
)

var randomKey []byte

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

func decodeQueryString(s string) map[string]string {
	output := make(map[string]string)
	pieces := strings.Split(s, "&")
	for _, v := range pieces {
		keyValue := strings.Split(v, "=")
		output[keyValue[0]] = keyValue[1]
	}
	return output
}

func encodeQueryString(m map[string]string) string {
	output := ""
	for key, value := range m {
		if output != "" {
			output += "&"
		}
		// Remove & and = from keys and values
		output += strings.Replace(strings.Replace(key, "&", "", -1), "=", "", -1) + "=" + strings.Replace(strings.Replace(value, "&", "", -1), "=", "", -1)
	}
	return output
}

func profile_for(email string) []byte {
	user := make(map[string]string)
	user["email"] = email
	user["uid"] = "10"
	user["role"] = "user"
	queryString := encodeQueryString(user)
		
	plaintext := []byte(queryString)
	plaintext = cryptoutil.Pkcs7padding(plaintext, len(plaintext) + cryptoutil.Pkcs7paddingCount(plaintext))
	return cryptoutil.AES128ECBEncrypt(plaintext, randomKey)
}

func decryptProfile(ciphertext []byte) map[string]string {
	plaintext := cryptoutil.AES128ECBDecrypt(ciphertext, randomKey)
	plaintext = cryptoutil.RemovePkcs7padding(plaintext)
	return decodeQueryString(string(plaintext))
}

func main() {
	bs := 16 // block size
	randomKey = cryptoutil.RandomBytes(bs)
	
	// First, find out how many input characters we need to fill up "n" full blocks. We can find
	// this because we know that when the plaintext is going to go from, say, 15 bytes to 16 bytes,
	// the cypthertext is going to bump from 16 to 32 bytes.
	
	previousLength := 0
	input := []byte{'a'}
	fullBlockInputLength := 0
	for {
		ciphertext := profile_for(string(input))
		if previousLength == 0 {
			previousLength = len(ciphertext)
			continue
		}
		if len(ciphertext) != previousLength {
			fullBlockInputLength = len(input)
			break
		}
		input = append(input, byte('a'))
	}
			
	// Since "&" characters are encoded, trying to decrypt the ciphertext starting from
	// the first character won't work. If we input a "&" character, it's going to be encoded
	// then encrypted, so we cannot compare it to anything else.
	
	// Instead, we try to decrypt starting from the the last characters
	
	var plaintext []byte	
	
	for i := 1; true; i++ {
		// Prepare an input so that as to isolate the last plaintext character(s). On the first iteration, 
		// the last block is going to be: "XPPPPPPPPPPPPPPP" - "X" being what we're looking for, and P being the padding.
	
		input = fillBytes('a', fullBlockInputLength + i)
		ciphertext := profile_for(string(input))
		bc := len(ciphertext) / bs // block count
		ciphertextLastBlock := ciphertext[(bc - 1) * bs : bc * bs]
				
		// Now prepare an input so that the second cypthertext block contains the same content as ciphertextLastBlock
		// except for the "X" that we are trying to guess.
		
		found := false
		for b := 0; b < 255; b++ {
			// Enough bytes so that the second block contains what we need
			guess := fillBytes('a', bs - len("email="))
			// Append the byte we are trying to guess
			guess = append(guess, byte(b))
			// Append the plaintext we've already found
			if len(plaintext) > 0 {
				guess = appendBytes(guess, plaintext)
			}
			// Finally, add the padding
			paddingSize := bs - i
			guess = appendBytes(guess, fillBytes(byte(paddingSize), paddingSize))
			
			ciphertext := profile_for(string(guess))
			secondBlock := ciphertext[bs : bs + bs]
			if cryptoutil.SliceEquals(secondBlock, ciphertextLastBlock) {
				// Prepend the byte to the plaintext
				temp := []byte{byte(b)}
				temp = appendBytes(temp, plaintext)
				plaintext = temp
				found = true
				break
			}
		}
		
		// If nothing was found, we might have hit an escaped character ("=" or "&") so exit.
		if !found {
			break
		}
	}
	
	// At this point, we know that the last value is "user",
	// so we can try to "replace" it with admin.
	
	log.Println("Plaintext:", string(plaintext))
	
	// To do so, we first prepare an input so that the second block of the plaintext starts with "admin"
	// followed by the padding. This block can then be inserted at the end of a specially crafted ciphertext.
	
	adminInput := fillBytes('a', bs - len("email="))
	adminInput = appendBytes(adminInput, []byte("admin"))
	paddingSize := bs - len("admin")
	adminInput = appendBytes(adminInput, fillBytes(byte(paddingSize), paddingSize))
	
	adminBlock := profile_for(string(adminInput))
	adminBlock = adminBlock[bs : bs + bs]
	
	// Now prepare a second input so that the last plaintext block start with "user" (followed by the padding)
	
	adminUserciphertextInput := fillBytes('a', fullBlockInputLength + len("user"))
	adminUserciphertext := profile_for(string(adminUserciphertextInput))
	
	// Finally replace the last block by the admin block we prepared above
	
	adminUserciphertext = adminUserciphertext[0:len(adminUserciphertext) - bs]
	adminUserciphertext = appendBytes(adminUserciphertext, adminBlock)
	
	// Check that the decrypted data is correct
	
	log.Println(decryptProfile(adminUserciphertext))
}
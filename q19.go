package main

import (
	"log"
	"./cryptoutil"
	"./charfreq"
	"crypto/aes"
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

// Convert an int to a byte slice
func Int64ToBytes(i int64, byteOrder binary.ByteOrder) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, byteOrder, i)
	return buffer.Bytes()
}

// Create the CTR mode keystream, which is AES(nonce || counter, key)
func CreateKeystream(key []byte, nonce int64, counter int64) []byte {
	cypher, _ := aes.NewCipher(key)
	keystreamKey := Int64ToBytes(nonce, binary.LittleEndian)
	keystreamKey = cryptoutil.AppendBytes(keystreamKey, Int64ToBytes(counter, binary.LittleEndian))
	output := make([]byte, 16)
	cypher.Encrypt(output, keystreamKey)
	return output
}

func AESCTRDecrypt(encrypted []byte, key []byte, nonce int64) []byte {
	output := make([]byte, len(encrypted))
	keystreamIndex := 0
	var counter int64 = 0
	var keystream []byte
	// Loop through each byte of the encrypted string and XOR it against the
	// keystream. Generate new keystreams as needed, incrementing the counter
	// every time.
	for i := 0; i < len(encrypted); i++ {
		if i == 0 || keystreamIndex >= 16 {
			keystream = CreateKeystream(key, nonce, counter)
			counter++
			keystreamIndex = 0
		}
		b := encrypted[i]
		output[i] = b ^ keystream[keystreamIndex]
		keystreamIndex++
	}
	
	return output
}

func AESCTREncrypt(plaintext []byte, key []byte, nonce int64) []byte {
	return AESCTRDecrypt(plaintext, key, nonce)
}

func mostFrequentByte(data []byte) (byte, int) {
	results := make(map[byte]int)
	for _, b := range data {
		_, ok := results[b]
		if !ok {
			results[b] = 1
		} else {
			results[b]++
		}
	}
	
	max := 0
	var output byte
	for k, v := range results {
		if v > max {
			max = v
			output = k
		}
	}
	
	return output, max
}

// Calculate the "score" of a key for the given ciphertexts, based on the
// number of common English letters, bigrams or trigrams. Using charfreq
// lib from Question 6.
func keyScore(ciphertexts [][]byte, key []byte) float64 {
	f := charfreq.NewCharFrequencies()
	var output float64 = 0
	for _, d := range ciphertexts {
		temp := cryptoutil.RepeatingKeyXor(d, key)
		s := f.ScorePlainText(temp)
		output += s
	}
	return output
}

// Given a list of ciphertexts, find the key.
// If a key is provided to the function, it will be used as
// a base to generate the new key. It means that each call
// to the function refines the key a bit more (higher keyScore). Not
// very smart (or fast) but does the job.

func findKey(ciphertexts [][]byte, guessedKey []byte) []byte {
	
	// Make a sorted list of the most common English characters (just a few should be enough)
	
	charsToTry := []byte{
		' ',
		'e',
		't',
		'a',
		'o',
		'i',
		'n',
		'h',
		's',
		'r',
		'd',
		'l',
		'u',
		'm',
		'c',
		'w',
		'g',
		'f',
		'y',
		'p',
		',',
		'.',
		'b',
		'k',
		'v',
	}
	
	// Get the length of the longest ciphertext
	
	maxLength := 0
	for _, ciphertext := range ciphertexts {
		if len(ciphertext) > maxLength {
			maxLength = len(ciphertext)
		}
	}
	
	// If we have an empty key, start by building one	
	
	buildKey := len(guessedKey) == 0
	
	// For each charToTry, build a temporary array containing the i-th bytes
	// of each ciphertext. Then get the most common byte in the array.
	//
	// Now we assume that this cipher-byte decrypts to the charToTry we are looking for.
	// So cipher-byte XOR charToTry = key-byte.
	
	// From this we create a temporary key based on the current key. If the score
	// of this temporary key is higher than the current one, it becomes the current key.
	
	// The latest bytes of some ciphertexts might not be properly decrypted since
	// there is less and less bytes to work on.
	
	for charToTryIndex, charToTry := range charsToTry {
		for i := 0; i < maxLength; i++ {
			var temp []byte
			for _, ciphertext := range ciphertexts {
				if i >= len(ciphertext) { continue }
				temp = append(temp, ciphertext[i])
			}
			b, _ := mostFrequentByte(temp)
			if charToTryIndex == 0 && buildKey {
				guessedKey = append(guessedKey, b ^ charToTry)
			} else {
				newKey := make([]byte, len(guessedKey))
				copy(newKey, guessedKey)
				newKey[i] = b ^ charToTry
				if keyScore(ciphertexts, newKey) > keyScore(ciphertexts, guessedKey) {
					guessedKey = newKey
				}
			}
		}
	}
	
	return guessedKey
}

func main() {
	data := []string{
		"SSBoYXZlIG1ldCB0aGVtIGF0IGNsb3NlIG9mIGRheQ==",
		"Q29taW5nIHdpdGggdml2aWQgZmFjZXM=",
		"RnJvbSBjb3VudGVyIG9yIGRlc2sgYW1vbmcgZ3JleQ==",
		"RWlnaHRlZW50aC1jZW50dXJ5IGhvdXNlcy4=",
		"SSBoYXZlIHBhc3NlZCB3aXRoIGEgbm9kIG9mIHRoZSBoZWFk",
		"T3IgcG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
		"T3IgaGF2ZSBsaW5nZXJlZCBhd2hpbGUgYW5kIHNhaWQ=",
		"UG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
		"QW5kIHRob3VnaHQgYmVmb3JlIEkgaGFkIGRvbmU=",
		"T2YgYSBtb2NraW5nIHRhbGUgb3IgYSBnaWJl",
		"VG8gcGxlYXNlIGEgY29tcGFuaW9u",
		"QXJvdW5kIHRoZSBmaXJlIGF0IHRoZSBjbHViLA==",
		"QmVpbmcgY2VydGFpbiB0aGF0IHRoZXkgYW5kIEk=",
		"QnV0IGxpdmVkIHdoZXJlIG1vdGxleSBpcyB3b3JuOg==",
		"QWxsIGNoYW5nZWQsIGNoYW5nZWQgdXR0ZXJseTo=",
		"QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
		"VGhhdCB3b21hbidzIGRheXMgd2VyZSBzcGVudA==",
		"SW4gaWdub3JhbnQgZ29vZCB3aWxsLA==",
		"SGVyIG5pZ2h0cyBpbiBhcmd1bWVudA==",
		"VW50aWwgaGVyIHZvaWNlIGdyZXcgc2hyaWxsLg==",
		"V2hhdCB2b2ljZSBtb3JlIHN3ZWV0IHRoYW4gaGVycw==",
		"V2hlbiB5b3VuZyBhbmQgYmVhdXRpZnVsLA==",
		"U2hlIHJvZGUgdG8gaGFycmllcnM/",
		"VGhpcyBtYW4gaGFkIGtlcHQgYSBzY2hvb2w=",
		"QW5kIHJvZGUgb3VyIHdpbmdlZCBob3JzZS4=",
		"VGhpcyBvdGhlciBoaXMgaGVscGVyIGFuZCBmcmllbmQ=",
		"V2FzIGNvbWluZyBpbnRvIGhpcyBmb3JjZTs=",
		"SGUgbWlnaHQgaGF2ZSB3b24gZmFtZSBpbiB0aGUgZW5kLA==",
		"U28gc2Vuc2l0aXZlIGhpcyBuYXR1cmUgc2VlbWVkLA==",
		"U28gZGFyaW5nIGFuZCBzd2VldCBoaXMgdGhvdWdodC4=",
		"VGhpcyBvdGhlciBtYW4gSSBoYWQgZHJlYW1lZA==",
		"QSBkcnVua2VuLCB2YWluLWdsb3Jpb3VzIGxvdXQu",
		"SGUgaGFkIGRvbmUgbW9zdCBiaXR0ZXIgd3Jvbmc=",
		"VG8gc29tZSB3aG8gYXJlIG5lYXIgbXkgaGVhcnQs",
		"WWV0IEkgbnVtYmVyIGhpbSBpbiB0aGUgc29uZzs=",
		"SGUsIHRvbywgaGFzIHJlc2lnbmVkIGhpcyBwYXJ0",
		"SW4gdGhlIGNhc3VhbCBjb21lZHk7",
		"SGUsIHRvbywgaGFzIGJlZW4gY2hhbmdlZCBpbiBoaXMgdHVybiw=",
		"VHJhbnNmb3JtZWQgdXR0ZXJseTo=",
		"QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
	}
	
	// Encrypt each line and save it to an array
	
	var ciphertexts [][]byte
	key := cryptoutil.RandomBytes(16)
	for _, line := range data {
		decoded, _ := base64.StdEncoding.DecodeString(line)
		ciphertext := AESCTREncrypt(decoded, key, 0)
		ciphertexts = append(ciphertexts, ciphertext)
	}
	
	// Each successive call to findKey refines it a bit more. It
	// seems three times always give the highest keyScore.
	
	guessedKey := findKey(ciphertexts, []byte{})
	guessedKey = findKey(ciphertexts, guessedKey)
	guessedKey = findKey(ciphertexts, guessedKey)
	
	for _, ciphertext := range ciphertexts {
		temp := cryptoutil.RepeatingKeyXor(ciphertext, guessedKey)
		log.Println(string(temp))
	}
	
	log.Println(keyScore(ciphertexts, guessedKey))
}
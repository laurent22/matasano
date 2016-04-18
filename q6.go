package main

import (
	"io/ioutil"
	"fmt"
	"encoding/base64"
	"./charfreq"
)

func XorBytes(bytes []byte, key byte) []byte {
	var output []byte 
	for i := 0; i < len(bytes); i++ {
		output = append(output, bytes[i] ^ key)
	}
	return output
}

func RepeatingKeyXor(bytes []byte, key []byte) []byte {
	keyIndex := 0
	var output []byte
	for i := 0; i < len(bytes); i++ {
		output = append(output, bytes[i] ^ key[keyIndex])
		keyIndex++
		if keyIndex >= len(key) { keyIndex = 0 }
	}
	return output
}

func HammingDistance(b1 []byte, b2 []byte) int {
	length := len(b1)
	if len(b2) < length { length = len(b2) }
	output := 0
	for i := 0; i < length; i++ {
		val := b1[i] ^ b2[i]
		for val > 0 {
			output++
			val &= val - 1
		}
	}
	return output
}

func main() {
	_ = fmt.Println
	
	content, _ := ioutil.ReadFile("q6_data.txt")
	data, _ := base64.StdEncoding.DecodeString(string(content))
	
	var winner struct {
		dist float64
		keySize int
	}
	
	winner.dist = 99999

	// Find the key size that yields the lowest Hamming distance. For this, divide
	// the data into chunks of KEYSIZE length. Then calculate the Hamming distance
	// between each successive chunks and add this to a sum.
	// Finally, average this sum by the total number of chunks.
	
	for keySize := 1; keySize <= 40; keySize++ {
		var distSum float64 = 0
		count := 0 
		for i := 0; i < len(data) + 1 - keySize; i++ {
			chunk1 := data[i:i+keySize]
			chunk2 := data[i+keySize:i + keySize + keySize]
			dist := float64(HammingDistance(chunk1, chunk2)) / float64(keySize)
			distSum += dist
			count++
		}
		total := distSum / float64(count)
		if total < winner.dist {
			winner.dist = total
			winner.keySize = keySize
		}
	}
	
	// Divice the data into block of KEYSIZE length
	
	keySize := winner.keySize
				
	var temp []([]byte)
	for i := 0; true; i += keySize {
		if i + keySize >= len(data) { break }
		temp = append(temp, data[i:i+keySize])
	}
	
	// Transpose the blocks: make a block that is the first byte of
	// every block, and a block that is the second byte of every block, etc.
	
	var transposed []([]byte)
	for keyIndex := 0; keyIndex < keySize; keyIndex++ {
		var current []byte
		for i := 0; i < len(temp); i++ {
			current = append(current, temp[i][keyIndex])
		}
		transposed = append(transposed, current)
	}
	
	// Apply each key to each block, and calculate the plain text score based
	// on the English letter frequency.
	
	freq := charfreq.NewCharFrequencies()
	
	var key byte
	var fullKey []byte
	for i := 0; i < len(transposed); i++ {
		line := transposed[i]
	
		var guess struct {
			score float64
			text []byte
			key byte
		}
		guess.score = 0
				
		for key = 0; key < 255; key++ {
			plainText := XorBytes(line, key)
			score := freq.ScorePlainTextByItem(plainText, charfreq.LETTER)
			if score > guess.score {
				guess.score = score
				guess.text = plainText
				guess.key = key
			}
		}
		fullKey = append(fullKey, guess.key)
	}
	
	fmt.Println("Key size:", keySize)
	fmt.Println("Key:", string(fullKey))
	
	decrypted := RepeatingKeyXor(data, fullKey)
	fmt.Println(string(decrypted))
}
package main

import (
	"io/ioutil"
	"encoding/base64"
	"strings"
	"log"
	"./cryptoutil"
	"./charfreq"
	"crypto/aes"
	"bytes"
	"encoding/binary"
)

func Int64ToBytes(i int64, byteOrder binary.ByteOrder) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, byteOrder, i)
	return buffer.Bytes()
}

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

func main() {
	// Load the data and encrypt each line
	
	content, _ := ioutil.ReadFile("q20_data.txt")
	lines := strings.Split(string(content), "\n")
	var ciphertexts [][]byte
	key := cryptoutil.RandomBytes(16)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" { continue }
		decoded, _ := base64.StdEncoding.DecodeString(line)
		ciphertext := AESCTREncrypt(decoded, key, 0)
		ciphertexts = append(ciphertexts, ciphertext)
	}
	
	// Get the length of the shortest ciphertext
	
	minLength := 9999
	for _, ciphertext := range ciphertexts {
		if len(ciphertext) < minLength {
			minLength = len(ciphertext)
		}
	}
	
	// Then do the same as in Question 6
		
	var transposed []([]byte)
	for i := 0; i < minLength; i++ {
		var current []byte
		for j := 0; j < len(ciphertexts); j++ {
			current = append(current, ciphertexts[j][i])
		}
		transposed = append(transposed, current)
	}
	
	freq := charfreq.NewCharFrequencies()
	
	var keyByte byte
	var fullKey []byte
	for i := 0; i < len(transposed); i++ {
		line := transposed[i]
	
		var guess struct {
			score float64
			text []byte
			keyByte byte
		}
		guess.score = 0
				
		for keyByte = 0; keyByte < 255; keyByte++ {
			plainText := cryptoutil.XorBytes(line, keyByte)
			score := freq.ScorePlainTextByItem(plainText, charfreq.LETTER)
			if score > guess.score {
				guess.score = score
				guess.text = plainText
				guess.keyByte = keyByte
			}
		}
		fullKey = append(fullKey, guess.keyByte)
	}

	for _, ciphertext := range ciphertexts {
		decrypted := cryptoutil.RepeatingKeyXor(ciphertext[0:len(fullKey)], fullKey)
		log.Println(string(decrypted))
	}
}
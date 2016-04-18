package main

import (
	"log"
	"./cryptoutil"
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

func AESCTREncrypt(encrypted []byte, key []byte, nonce int64) []byte {
	return AESCTRDecrypt(encrypted, key, nonce)
}

func main() {
	data, _ := base64.StdEncoding.DecodeString("L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ==")
	plaintext := AESCTRDecrypt(data, []byte("YELLOW SUBMARINE"), 0)
	log.Println(string(plaintext))
	
	enc := AESCTREncrypt([]byte("testing"), []byte("YELLOW SUBMARINE"), 0)
	log.Println("%x", enc)
	dec := AESCTRDecrypt(enc, []byte("YELLOW SUBMARINE"), 0)
	log.Println(string(dec))
}
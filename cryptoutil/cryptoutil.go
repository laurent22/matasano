package cryptoutil

import (
	"crypto/aes"
	"math/rand"	
	"time"
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

func AES128ECBEncrypt(plain []byte, key []byte) []byte {
	cypher, _ := aes.NewCipher(key)
	blockSize := 16
	output := make([]byte, len(plain))
	for i := 0; i < len(output); i += blockSize {
		cypher.Encrypt(output[i:i+blockSize], plain[i:i+blockSize]) 
	}
	return output
}

func AES128CBCEncrypt(plain []byte, key []byte, iv []byte) []byte {
	cypher, _ := aes.NewCipher(key)
	blockSize := 16
	output := make([]byte, len(plain))
	previousCypherText := iv
	// Iterate over each block, starting from the first one
	for i := 0; i < len(output); i += blockSize {
		// Xor the current block with the previous cyphertext
		plainBlock := RepeatingKeyXor(plain[i:i+blockSize], previousCypherText)
		// Encrypt the current block
		cypher.Encrypt(output[i:i+blockSize], plainBlock) 
		previousCypherText = output[i:i+blockSize]
	}
	return output
}

func AES128CBCDecrypt(encrypted []byte, key []byte, iv []byte) []byte {
	cypher, _ := aes.NewCipher(key)
	bs := 16 // block size
	output := make([]byte, len(encrypted))
	blockCount := len(encrypted) / bs
	// Iterate over each block starting from the last one
	for i := blockCount - 1; i >= 0; i-- {
		// Get the current block i
		encBlock := encrypted[i*bs:(i+1)*bs]
		// Get the block before i
		var previousBlock []byte
		if i == 0 {
			previousBlock = iv
		} else {
			previousBlock = encrypted[(i-1)*bs:i*bs]
		}
		// Decrypt the current block
		cypher.Decrypt(output[i*bs:(i+1)*bs], encBlock)
		// Xor it against the previous block
		temp := RepeatingKeyXor(output[i*bs:(i+1)*bs], previousBlock)
		for j := 0; j < bs; j++ {
			output[j + i*bs] = temp[j]
		}
	}
	return output
}

func SliceEquals(slice1 []byte, slice2 []byte) bool {
	if len(slice1) != len(slice2) { return false }
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] { return false }
	}
	return true
}

func IsECBEncrypted(data []byte) bool {
	if len(data) == 0 { return false }
	
	// Compare each 16-byte block to the following ones.
	// If two blocks are identical, it means we have some
	// ECB encrypted data.
	
	blockSize := 16
	blockCount := int(float64(len(data)) / float64(blockSize))
	for i := 0; i < blockCount - 1; i++ {
		slice1 := data[i * blockSize : i * blockSize + blockSize]
		for j := i + 1; j < blockCount; j++ {
			slice2 := data[j * blockSize : j * blockSize + blockSize]
			if SliceEquals(slice1, slice2) {
				return true
			}
		}
	}
	
	return false
}

// PKCS#7 padding.
func Pkcs7padding(data []byte, length int) []byte {
	p := length - len(data)
	if p < 0 {
		return data
	}
	output := data
	for i := 0; i < p; i++ {
		output = append(output, byte(p))
	}
	return output
}

// Calculates the number of extra bytes required
// to pad some data with PKCS#7 padding. If
// no padding is required (data length is a multiple
// of 16), we pad with 16 bytes, so that the padding
// can still be identified.
func Pkcs7paddingCount(data []byte) int {
	output := 16 - len(data) % 16
	if output == 0 { output = 16 }
	return output
}

func RemovePkcs7padding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	paddingSize := int(data[len(data) - 1])
	return data[0:len(data) - paddingSize]
}

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

func RandomBytes(count int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var output []byte
	for i := 0; i < count; i++ {
		output = append(output, byte(r.Intn(256)))
	}
	return output
}

func RandomChars(count int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var output []byte
	for i := 0; i < count; i++ {
		output = append(output, byte(48 + r.Intn(126 - 48))) // 48 - 126
	}
	return output
}

func AppendBytes(dest []byte, source []byte) []byte {
	output := dest
	for _, b := range source {
		output = append(output, b)
	}
	return output
}

func PrependBytes(dest []byte, source []byte) []byte {
	var output []byte
	output = AppendBytes(output, source)
	output = AppendBytes(output, dest)
	return output
}

func PrependByte(dest []byte, source byte) []byte {
	var output []byte
	var temp []byte
	temp = append(temp, source)
	output = AppendBytes(output, temp)
	output = AppendBytes(output, dest)
	return output
}

func FillBytes(b byte, length int) []byte {
	var output []byte
	for len(output) < length {
		output = append(output, b)
	}
	return output
}
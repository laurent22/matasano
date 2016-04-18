package main

import (
	"log"
	"./cryptoutil"
	"strings"
	"net/url"
	//"encoding/hex"
)

var randomKey []byte
var randomIv []byte

func removePkcs7padding(data []byte) ([]byte, bool) {
	// Invalid if the data is empty
	if len(data) == 0 {
		return []byte{}, false
	}
	// Get the padding size
	paddingSize := int(data[len(data) - 1])
	// Invalid if the padding is longer than the data
	if len(data) - paddingSize < 0 {
		return []byte{}, false
	}
	// Then check each byte and make sure it's the same value as the padding.
	for i := len(data) - 1; i >= len(data) - paddingSize; i-- {
		if int(data[i]) != paddingSize {
			return []byte{}, false
		}
	}
	// Strip off the padding and return the data
	return data[0:len(data) - paddingSize], true
}

func encrypt(message []byte) []byte {
	var plaintext []byte
	plaintext = cryptoutil.AppendBytes(plaintext, []byte("comment1=cooking%20MCs;userdata="))
	plaintext = cryptoutil.AppendBytes(plaintext, []byte(url.QueryEscape(string(message))))
	plaintext = cryptoutil.AppendBytes(plaintext, []byte(";comment2=%20like%20a%20pound%20of%20bacon"))
	plaintext = cryptoutil.Pkcs7padding(plaintext, len(plaintext) + cryptoutil.Pkcs7paddingCount(plaintext))
	return cryptoutil.AES128CBCEncrypt(plaintext, randomKey, randomIv)
}

func decrypt(ciphertext []byte) []byte {
	output, ok := removePkcs7padding(cryptoutil.AES128CBCDecrypt(ciphertext, randomKey, randomIv))
	if !ok {
		return ciphertext
	}
	return output
}

func isAdmin(ciphertext []byte) bool {
	plaintext := string(decrypt(ciphertext))
	items := strings.Split(plaintext, ";")
	for _, item := range items {
		kv := strings.Split(item, "=")
		k, _ := url.QueryUnescape(kv[0])
		v, _ := url.QueryUnescape(kv[1])
		if k == "admin" && v == "true" {
			return true
		}
	}
	return false
}

func main() {
	bs := 16 // block size
	randomKey = cryptoutil.RandomBytes(bs)
	randomIv = cryptoutil.RandomBytes(bs)
	
	// Why does CBC mode produces the identical 1-bit error in the next ciphertext block?
	//
	// => because each block is XORed against the previous one, so a change of 1 bit in block n is going to
	// flip one bit in block n + 1.
	
	// First, estimate the size of the data by inputting an empty string:
	
	ciphertext := encrypt([]byte(""))
	
	// Then create an input string that is three times bigger than the existing data. That way, no matter what the prefix
	// or suffix is, we know that the middle block will contain a part of our input.
	
	input := cryptoutil.FillBytes('a', len(ciphertext) * 3)
	ciphertext = encrypt(input)
	bi := (len(ciphertext) / bs) / 2 // a block index where we know some of our input is
		
	// Since the block mode is CBC, we know that our input (assumed it's in "block 2" to simplify) is going to be:
	// 1. XORed against block 1
	// 2. Then AES-encrypted
	//
	// When decrypting, block 2 is going to be:
	// 1. AES-Decrypted.
	// 2. XORed against block 1
	//
	// Basically, on block 2, we currently have:
	// ciphertextBlock1 ^ input
	//
	// We need to create newBlock2 so that, when XORed against (ciphertextBlock1 ^ input), it produces adminString:
	//
	// So we want:
	//
	// (ciphertextBlock1 ^ input) ^ newBlock2 = adminString
	//
	// Which is equivalent to:
	//
	// newBlock2 = (ciphertextBlock1 ^ input) ^ adminString	
	
	ciphertextBlock1 := ciphertext[bi * bs:bi * bs + bs]
	newBlock2 := cryptoutil.RepeatingKeyXor([]byte(";admin=true;aaaa"), input)
	newBlock2 = cryptoutil.RepeatingKeyXor(newBlock2, ciphertextBlock1)
	
	// Now that we have the new value for block 1, build the new ciphertext
	
	newciphertext := ciphertext[0:bi * bs]
	newciphertext = cryptoutil.AppendBytes(newciphertext, newBlock2)
	newciphertext = cryptoutil.AppendBytes(newciphertext, ciphertext[len(newciphertext):len(ciphertext)])

	log.Println("Is admin:", isAdmin(newciphertext))
	log.Println(string(decrypt(newciphertext)))	
}
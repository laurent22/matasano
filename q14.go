package main

import (
	"log"
	"./cryptoutil"
	"math/rand"	
	"math"	
	"time"
)

var randomKey []byte
var randomPrefix []byte

func encryption_oracle(plaintext []byte) []byte {
	plaintext = cryptoutil.PrependBytes(plaintext, randomPrefix)
	plaintext = cryptoutil.AppendBytes(plaintext, []byte("target bytes"))
	plaintext = cryptoutil.Pkcs7padding(plaintext, len(plaintext) + cryptoutil.Pkcs7paddingCount(plaintext))
	return cryptoutil.AES128ECBEncrypt(plaintext, randomKey)
}

// Helper function to simply decryption code. Same as the regular encryption_oracle function
// from Q12, except that it does some extra work to remove the first prefix blocks
// (prefix length must be known first)
func encryption_oracle_no_prefix(plaintext []byte, prefixLength int) []byte {
	prefixBlockCount := int(math.Ceil(float64(prefixLength) / 16))
	extraByteCount := 16 - prefixLength
	if extraByteCount == 0 {
		return encryption_oracle(plaintext)
	}
	plaintext = cryptoutil.PrependBytes(plaintext, cryptoutil.FillBytes('a', extraByteCount))
	output := encryption_oracle(plaintext)
	return output[prefixBlockCount * 16 : len(output)]
}

func main() {
	bs := 16 // block size
	randomKey = cryptoutil.RandomBytes(bs)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomPrefix = cryptoutil.RandomChars(r.Intn(8))
	
	// First, we try to find out the length of the prefix.
	
	// To do so, create an input that can generate at least two repeating blocks:
	
	input := cryptoutil.FillBytes('a', bs * 3)
	ciphertext := encryption_oracle(input)
	
	// Now find where the first of these repeating block starts:
	
	prefixLength := 0
	bc := len(ciphertext) / bs // block count
	for bi := 0; bi < bc - 1; bi++ { // block index
		block1 := ciphertext[bi * bs : bi * bs + bs]
		block2 := ciphertext[(bi + 1) * bs : (bi + 1) * bs + bs]
		if cryptoutil.SliceEquals(block1, block2) {
			// Found two repeating blocks.
			// Now prepend different bytes ('b' instead of 'a') to the input until block1 and block2 are different
			// When this happens, it means we filled up the rest of the prefix block with our
			// data and can find from this the length of the prefix.
			byteCount := 0
			for {
				byteCount++
				input = cryptoutil.PrependByte(input, 'b')
				ciphertext := encryption_oracle(input)
				newBlock1 := ciphertext[bi * bs : bi * bs + bs]
				if !cryptoutil.SliceEquals(block1, newBlock1) {
					// Found the length of the prefix
					prefixLength = bs - (byteCount - 1)
					break
				}
			}
			break
		}
	}
	
	// Now that we know the length of the prefix, we can create an input string that moves the target
	// bytes to the beginning of a block. From there, we can use the same method as in Q12 to decrypt
	// the data, one byte at a time.
	
	log.Println(prefixLength)
	log.Println(randomPrefix)
	
	var plaintext []byte
	input = cryptoutil.FillBytes('f', bs - 1)
	bi := 0
	for bi < len(ciphertext) / bs {
		ciphertext = encryption_oracle_no_prefix(input, prefixLength)
		for i := 0; i < 256; i++ {
			var guess []byte
			if len(plaintext) < bs {
				guess = input
				guess = cryptoutil.AppendBytes(guess, plaintext)
			} else {
				guess = plaintext[len(plaintext) - (bs - 1):len(plaintext)]
			}
			guess = append(guess, byte(i))
			
			b1 := encryption_oracle_no_prefix(guess, prefixLength)[0:bs]
			b2 := ciphertext[bs * bi:bs * bi + bs]
			if cryptoutil.SliceEquals(b1, b2) {
				plaintext = append(plaintext, byte(i))
				break
			}
		}
		
		if len(input) == 0 {
			input = cryptoutil.FillBytes('f', bs)
			bi++
		}
		
		input = input[0:len(input)-1]
	}
	
	log.Println("Plain text:", string(plaintext))
}
package main

import (
	"encoding/hex"
	"fmt"
)

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

func main() {
	input := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"
	encrypted := RepeatingKeyXor([]byte(input), []byte(key))
	fmt.Println(hex.EncodeToString(encrypted))
}
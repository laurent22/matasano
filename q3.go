package main

import (
	"encoding/hex"
	"fmt"
	"./charfreq"
)

func main() {
	frequencies := charfreq.NewCharFrequencies()
	source, _ := hex.DecodeString("1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736")
	
	var highestScore float64 = 0
	var highestScoreText []byte
	var highestScoreKey byte
	
	// Loop though each possible key
	var key byte
	for key = 0; key < 255; key++ {
		var plainText []byte
		// Try to build the plain text from the input and key
		for i := 0; i < len(source); i++ {
			plainText = append(plainText, source[i] ^ key)
		}
		// Get the score of the given plain text
		score := frequencies.ScorePlainText(plainText)
		// Save this score if it's the highest
		if score > highestScore {
			highestScore = score
			highestScoreText = plainText
			highestScoreKey = key
		}
	}
	
	fmt.Println("Key:", string(highestScoreKey))
	fmt.Println("Plain:", string(highestScoreText))
}
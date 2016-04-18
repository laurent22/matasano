package main

import (
	"./charfreq"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
)

func XorBytes(bytes []byte, key byte) []byte {
	var output []byte 
	for i := 0; i < len(bytes); i++ {
		output = append(output, bytes[i] ^ key)
	}
	return output
}

func main() {
	frequencies := charfreq.NewCharFrequencies()
	content, _ := ioutil.ReadFile("q1_data.txt")
	lines := strings.Split(string(content), "\n")
	
	var winner struct {
		score float64
		text []byte
		key byte
	}
	winner.score = 0
    
    // Loop through all the lines and apply all possible keys to each of them.
    // Then calculate the score of the plain text and save the line with the
    // highest score.
	var key byte
	for i := 0; i < len(lines); i++ {
		line, _ := hex.DecodeString(strings.TrimSpace(lines[i]))
				
		for key = 0; key < 255; key++ {
			plainText := XorBytes(line, key)
			score := frequencies.ScorePlainText(plainText)
			if score > winner.score {
				winner.score = score
				winner.text = plainText
				winner.key = key
			}
		}
	}
	
	fmt.Println("Key:", string(winner.key))
	fmt.Println("Score:", winner.score)
	fmt.Println("Text:", string(winner.text))
}
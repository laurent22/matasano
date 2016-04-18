package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
)

func sliceEquals(slice1 []byte, slice2 []byte) bool {
	if len(slice1) != len(slice2) { return false }
	for i := 0; i < len(slice1); i++ {
		if slice1[i] != slice2[i] { return false }
	}
	return true
}

func main() {
	content, _ := ioutil.ReadFile("q8_data.txt")
	lines := strings.Split(string(content), "\n")
	
	// Divide the line into blocks of 16 bytes, then compare each block to the
	// following ones. If two blocks are identical, it means we have the
	// the ECB encrypted line.
	
	blockSize := 16
	for lineIndex := 0; lineIndex < len(lines); lineIndex++ {
		line := lines[lineIndex]
		bytes, _ := hex.DecodeString(line)
		if len(bytes) == 0 { continue }
		blockCount := int(float64(len(bytes)) / float64(blockSize))
		for i := 0; i < blockCount - 1; i++ {
			slice1 := bytes[i * blockSize : i * blockSize + blockSize]
			for j := i + 1; j < blockCount; j++ {
				slice2 := bytes[j * blockSize : j * blockSize + blockSize]
				if sliceEquals(slice1, slice2) {
					fmt.Println("Detected:", line)
					return
				}
			}
		}
	}
}
package main

import (
	"log"
)

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

func main() {
	test := []byte("ICE ICE BABY\x04\x04\x04\x04")
	r, ok := removePkcs7padding(test)
	log.Println(string(r), ok)
	
	test = []byte("ICE ICE BABY\x05\x05\x05\x05")
	r, ok = removePkcs7padding(test)
	log.Println(string(r), ok)
	
	test = []byte("ICE ICE BABY\x01\x02\x03\x04")
	r, ok = removePkcs7padding(test)
	log.Println(string(r), ok)
}
package main

import (
	"log"
)

func rightShift() {
	var n int = 3076246805
	log.Printf("%b", n)
	var n1 int = n ^ (n >> 18)
	log.Printf("%b", n1)

	// t := n1 >> 18
	// log.Println(t ^ n1)
	
	var leftPart int = n1 & 0xffffc000
	log.Printf("%b", leftPart)
	rightPart := n1 >> 18
	rightPart = (rightPart ^ n1) & 0x7fff
	log.Printf("%b", rightPart)
	
	leftPart = rightPart | leftPart
	log.Printf("%d", leftPart)	
	

	// 10110111010110111100110100010101
	// 00000000000000000010110111010110  111100110100010101
	// 10110111010110111110000011000011
}


func main() {
	// y = y ^ (y >> 11)
	// y = y ^ ((y << 7) & 0x9d2c5680)
	// y = y ^ ((y << 15) & 0xefc60000)
	// y = y ^ ((y >> 18))
	
	
	// 0x9d2c5680 =
	// 10011101001011000101011010000000

	// 10110111010110111100110100010101
	// 10101101111001101000101010000000    << 7
	// 
	// 10101101111001101000101010000000 &
	// 10011101001011000101011010000000
	//
	//=
}
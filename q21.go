package main

import (
	"log"
	"./mtrand"
	"time"
	"math/rand"	
)

// Guesses the seed for the given random number, starting from the given time.
// The function assumes that the RNG was seeded with the current Unix timestamp,
// and attempts to find which timestamp it was.
func GuessSeed(randomNumber int, fromTime int) int {
	// Try seeds starting from the given time, and stop after an arbitrary
	// number of attempts.
	for seed := fromTime; fromTime - seed < 100000; seed-- {
		g := mtrand.NewGenerator()
		g.Initialize(seed)
		n := g.GetInt()
		if n == randomNumber {
			return seed
		}
	}
	return 0
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	waitingTime := r.Intn(500)
	log.Printf("Waiting for %d seconds...", waitingTime)
		
	time.Sleep(time.Duration(waitingTime) * time.Second)
	
	seed := int(time.Now().Unix())
	log.Printf("Seeding the generator with %d...", seed)
	
	g := mtrand.NewGenerator()
	g.Initialize(seed)
	
	waitingTime = r.Intn(500)
	log.Printf("Waiting for %d seconds...", waitingTime)

	time.Sleep(time.Duration(waitingTime) * time.Second)
	
	randomNumber := g.GetInt()
	log.Printf("Random number: %d", randomNumber)
	
	log.Println("Trying to find the seed...")
	guessedSeed := GuessSeed(randomNumber, int(time.Now().Unix()))
	
	log.Printf("Guessed seed: %d", guessedSeed)
}
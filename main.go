package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var (
	wordLength = 6 // Number of letters in the word.
)

func parseFlags() {
	flag.IntVar(&wordLength, "l", wordLength, "Word length in number of letters")
	flag.Parse()
}

func initialize() {
	parseFlags()

	fmt.Printf("Word length: %d\n", wordLength)

	// Initialize random values.
	rand.Seed(time.Now().UnixNano())
}

func main() {
	initialize()
}

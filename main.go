package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	wildcard        = "-"
	wildcardFlag    = "w"
	excludedFlag    = "x"
	wordPatternFlag = "p"
)

var (
	wordPattern     = "" // Pattern to match, length is inferred from string length.
	wildcardLetters = "" // Letters that can appear in any position where there is a wildecard placeholder.
	excludedLetters = "" // Letters that cannot appear in the position where they are specified.
	tryTheseLetters = map[string]int{}
	matchingWords   []string
)

func parseFlags() {
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + wildcard + "' 1) must include all letters specified by the -" + wildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + excludedFlag + " flag. Example value of 't" + strings.Repeat(wildcard, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&wordPattern, wordPatternFlag, wordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + wildcard + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + wildcard + "' in the -" + wordPatternFlag + " flag."
	flag.StringVar(&wildcardLetters, wildcardFlag, wildcardLetters, wildcardHelp)
	flag.StringVar(&excludedLetters, excludedFlag, excludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	flag.Parse()
}

func initialize() {
	parseFlags()

	fmt.Printf("Word length: %d\n", len(wordPattern))
	fmt.Printf("Word pattern: '%s'\n", wordPattern)
	fmt.Printf("Wild Card letters: '%s'\n", wildcardLetters)
	fmt.Printf("Excluded letters: '%s'\n", excludedLetters)

	if len(wordPattern) == 0 {
		fmt.Println("You must specify a -p <Pattern to Match>")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func filterWord(word string) {
	// filter length.
	if len(word) != len(wordPattern) {
		return
	}

	// filter excluded letters.
	for _, letter := range excludedLetters {
		if strings.Contains(word, string(letter)) {
			return
		}
	}

	// filter included letters.
	for _, letter := range wildcardLetters {
		if !strings.Contains(word, string(letter)) {
			return
		}
	}

	// filter for positional pattern match.
	for i, letter := range word {
		switch string(wordPattern[i : i+1]) {
		case string(letter):
			continue
		case wildcard:
			if !strings.Contains(wildcardLetters, string(letter)) &&
				!strings.Contains(wordPattern, string(letter)) {
				tryTheseLetters[string(letter)]++
			}
			continue
		default:
			return
		}
	}

	matchingWords = append(matchingWords, word)
}

func getMatchingWords() {
	// open file
	f, err := os.Open("english-words/words_alpha.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		filterWord(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func printMatchingWords() {
	if len(matchingWords) == 0 {
		fmt.Println("\nNo matching words!")
		return
	}

	fmt.Printf("\nPossible matching words (%d):\n", len(matchingWords))
	lineLength := 0
	for _, word := range matchingWords {
		fmt.Print(word)
		lineLength += len(word) + 1
		if lineLength+len(word) > 80 {
			fmt.Println()
			lineLength = 0
		} else {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

func printLettersToTry() {
	if len(tryTheseLetters) == 0 {
		fmt.Println("\nNo additional letters to try!")
		return
	}

	keys := make([]string, 0, len(tryTheseLetters))
	for k := range tryTheseLetters {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return tryTheseLetters[keys[i]] > tryTheseLetters[keys[j]]
	})

	fmt.Printf("\nTry these letters (%d):\n", len(keys))
	for _, k := range keys {
		fmt.Printf("%s=%d ", k, tryTheseLetters[k])
	}
	fmt.Println()
}

func main() {
	initialize()
	getMatchingWords()
	printMatchingWords()
	printLettersToTry()
}

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
	maxletters      = 9
)

var (
	wordFile        = "CSW21.txt" // Name/Path of text file containing 1 word per line.
	wordPattern     = ""          // Pattern to match, length is inferred from string length.
	wildcardLetters = ""          // Letters that can appear in any position where there is a wildecard placeholder.
	excludedLetters = ""          // Letters that cannot appear in the position where they are specified.
	matchingWords   []string
	noParkDisSpace  [maxletters]string
	printStatistics = false
)

func parseFlags() {
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + wildcard + "' 1) must include all letters specified by the -" + wildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + excludedFlag + " flag. Example value of 't" + strings.Repeat(wildcard, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&wordPattern, wordPatternFlag, wordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + wildcard + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + wildcard + "' in the -" + wordPatternFlag + " flag."
	flag.StringVar(&wildcardLetters, wildcardFlag, wildcardLetters, wildcardHelp)
	flag.StringVar(&excludedLetters, excludedFlag, excludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	flag.StringVar(&wordFile, "f", wordFile, "Word File: Name/Path of ASCII text file containing one word per line.")
	for disSpace := 0; disSpace < maxletters; disSpace++ {
		noParkDisSpaceHelp := "Letters that don't belong in this position: Letters that appear in the word, but not in postion #" + fmt.Sprintf("%d", disSpace+1) + " Example value of 'ies' means that 'i', 'e', or 's' cannot appear in position #" + fmt.Sprintf("%d", disSpace+1) + "."
		flag.StringVar(&noParkDisSpace[disSpace], fmt.Sprintf("%d", disSpace+1), noParkDisSpace[disSpace], noParkDisSpaceHelp)
	}
	flag.BoolVar(&printStatistics, "s", printStatistics, "Print statistics of letter distribution for each letter position.")
	flag.Parse()
}

func initialize() {
	parseFlags()

	fmt.Printf("Word file: %s\n", wordFile)
	wordPattern = strings.ToLower(wordPattern)
	fmt.Printf("Word length: %d\n", len(wordPattern))
	fmt.Printf("Word pattern: '%s'\n", wordPattern)
	wildcardLetters = strings.ToLower(wildcardLetters)
	fmt.Printf("Wild Card letters: '%s'\n", wildcardLetters)
	excludedLetters = strings.ToLower(excludedLetters)
	fmt.Printf("Excluded letters: '%s'\n", excludedLetters)

	for disSpace := 0; disSpace < maxletters; disSpace++ {
		if len(noParkDisSpace[disSpace]) > 0 {
			noParkDisSpace[disSpace] = strings.ToLower(noParkDisSpace[disSpace])
			fmt.Printf("Can't use letters in postion #%d: '%s'\n", disSpace+1, noParkDisSpace[disSpace])
		}
	}

	if _, err := os.Stat(wordFile); err != nil {
		fmt.Printf("%s does not exist. Please download/specify a valid word file.\n", wordFile)
		os.Exit(1)
	}

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

	// filter words that don't contain wildcard letters.
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
			if len(noParkDisSpace[i]) > 0 {
				// Remove letters that can't be in a specific position.
				if strings.Contains(noParkDisSpace[i], string(letter)) {
					return
				}
			}
			continue
		default:
			return
		}
	}

	matchingWords = append(matchingWords, word)
}

func getMatchingWords() {
	f, err := os.Open(wordFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		filterWord(strings.ToLower(scanner.Text()))
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

	var tryTheseLetters = map[string]int{}

	for _, word := range matchingWords {
		for _, letter := range word {
			if !strings.Contains(wildcardLetters, string(letter)) &&
				!strings.Contains(wordPattern, string(letter)) {
				tryTheseLetters[string(letter)]++
			}
		}
	}

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

func printWordStatistics() {

	if !printStatistics {
		return
	}

	fmt.Println()
	if len(matchingWords) == 0 {
		fmt.Println("\nNo statistics to print!")
		return
	}

	for position := 0; position < len(wordPattern); position++ {
		letterDistribution := map[string]int{}
		for _, word := range matchingWords {
			letterDistribution[string(word[position])]++
		}

		keys := make([]string, 0, len(letterDistribution))
		for k := range letterDistribution {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return letterDistribution[keys[i]] > letterDistribution[keys[j]]
		})

		fmt.Printf("Letter distribution for position #%d:\n", position+1)
		for _, k := range keys {
			if letterDistribution[k] > 0 {
				fmt.Printf("%s=%d ", k, letterDistribution[k])
			}
		}
		fmt.Println()
	}

}

func main() {
	initialize()
	getMatchingWords()
	printMatchingWords()
	printLettersToTry()
	printWordStatistics()
}

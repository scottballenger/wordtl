package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	WildcardChar    = "-"
	WildcardFlag    = "w"
	ExcludedFlag    = "x"
	WordPatternFlag = "p"
	FileFlag        = "f"
	StaisticsFlag   = "s"
	MaxLetters      = 9
)

var (
	WordFile        = "CSW21.txt" // Name/Path of text file containing 1 word per line.
	WordPattern     = ""          // Pattern to match, length is inferred from string length.
	WildcardLetters = ""          // Letters that can appear in any position where there is a wildecard placeholder.
	ExcludedLetters = ""          // Letters that cannot appear in the position where they are specified.
	NoParkDisSpace  [MaxLetters]string
	PrintStatistics = false
)

func parseFlags() {
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + WildcardChar + "' 1) must include all letters specified by the -" + WildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + ExcludedFlag + " flag. Example value of 't" + strings.Repeat(WildcardChar, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&WordPattern, WordPatternFlag, WordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + WildcardChar + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + WildcardChar + "' in the -" + WordPatternFlag + " flag."
	flag.StringVar(&WildcardLetters, WildcardFlag, WildcardLetters, wildcardHelp)
	flag.StringVar(&ExcludedLetters, ExcludedFlag, ExcludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	flag.StringVar(&WordFile, FileFlag, WordFile, "Word File: Name/Path of ASCII text file containing one word per line.")
	for disSpace := 0; disSpace < MaxLetters; disSpace++ {
		noParkDisSpaceHelp := "Letters that don't belong in this position: Letters that appear in the word, but not in postion #" + fmt.Sprintf("%d", disSpace+1) + " Example value of 'ies' means that 'i', 'e', or 's' cannot appear in position #" + fmt.Sprintf("%d", disSpace+1) + "."
		flag.StringVar(&NoParkDisSpace[disSpace], fmt.Sprintf("%d", disSpace+1), NoParkDisSpace[disSpace], noParkDisSpaceHelp)
	}
	flag.BoolVar(&PrintStatistics, StaisticsFlag, PrintStatistics, "Print statistics of letter distribution for each letter position.")
	flag.Parse()
}

func initialize() {
	parseFlags()

	fmt.Printf("Word file: %s\n", WordFile)
	WordPattern = strings.ToLower(WordPattern)
	fmt.Printf("Word length: %d\n", len(WordPattern))
	fmt.Printf("Word pattern: '%s'\n", WordPattern)
	WildcardLetters = strings.ToLower(WildcardLetters)
	fmt.Printf("Wild Card letters: '%s'\n", WildcardLetters)
	ExcludedLetters = strings.ToLower(ExcludedLetters)
	fmt.Printf("Excluded letters: '%s'\n", ExcludedLetters)

	for disSpace := 0; disSpace < MaxLetters; disSpace++ {
		if len(NoParkDisSpace[disSpace]) > 0 {
			NoParkDisSpace[disSpace] = strings.ToLower(NoParkDisSpace[disSpace])
			fmt.Printf("Can't use letters in postion #%d: '%s'\n", disSpace+1, NoParkDisSpace[disSpace])
		}
	}

	if _, err := os.Stat(WordFile); err != nil {
		fmt.Printf("%s does not exist. Please download/specify a valid word file.\n", WordFile)
		os.Exit(1)
	}

	if len(WordPattern) == 0 {
		fmt.Println("You must specify a -p <Pattern to Match>")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func wordMatch(word string,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	noParkDisSpace [MaxLetters]string) bool {

	// filter length.
	if len(word) != len(wordPattern) {
		return false
	}

	// filter excluded letters.
	for _, letter := range excludedLetters {
		if strings.Contains(word, string(letter)) {
			return false
		}
	}

	// filter words that don't contain WildcardChar letters.
	for _, letter := range wildcardLetters {
		if !strings.Contains(word, string(letter)) {
			return false
		}
	}

	// filter for positional pattern match.
	for i, letter := range word {
		switch string(wordPattern[i : i+1]) {
		case string(letter):
			continue
		case WildcardChar:
			if len(noParkDisSpace[i]) > 0 {
				// Remove letters that can't be in a specific position.
				if strings.Contains(noParkDisSpace[i], string(letter)) {
					return false
				}
			}
			continue
		default:
			return false
		}
	}
	return true
}

func getMatchingWords(
	wordFileHandle io.Reader,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	noParkDisSpace [MaxLetters]string) []string {

	var matchingWords []string

	scanner := bufio.NewScanner(wordFileHandle)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if wordMatch(word, wordPattern, excludedLetters, wildcardLetters, noParkDisSpace) {
			matchingWords = append(matchingWords, word)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return matchingWords
}

func printMatchingWords(matchingWords []string) {
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

func printLettersToTry(matchingWords []string, wordPattern string, wildcardLetters string) {
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

func printWordStatistics(matchingWords []string, wordPattern string) {
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

	f, err := os.Open(WordFile)
	if err != nil {
		log.Fatal(err)
	}
	matchingWords := getMatchingWords(f, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
	f.Close()

	printMatchingWords(matchingWords)
	printLettersToTry(matchingWords, WordPattern, WildcardLetters)
	if PrintStatistics {
		printWordStatistics(matchingWords, WordPattern)
	}
}

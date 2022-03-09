package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"wordtl/words"
)

const (
	WordLengthFlag      = "l"
	WordPatternFlag     = "p"
	WildcardFlag        = "w"
	ExcludedFlag        = "x"
	FileFlag            = "f"
	StaisticsFlag       = "s"
	MaxWordsToPrintFlag = "m"
)

var (
	WordLength      = 5
	WordPattern     = "" // Pattern to match, length must be WordLength number of letters.
	WordFile        = "" // Name/Path of text file containing 1 word per line.
	WildcardLetters = "" // Letters that can appear in any position where there is a wildecard placeholder.
	ExcludedLetters = "" // Letters that cannot appear in the position where they are specified.
	NoParkDisSpace  [words.MaxLetters]string
	PrintStatistics = false
	MaxWordsToPrint = 100
)

func parseFlags() {
	flag.IntVar(&WordLength, WordLengthFlag, WordLength, "Word Length: Number of letters to match. Wordle is 5 letters.")
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + words.WildcardChar + "' 1) must include all letters specified by the -" + WildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + ExcludedFlag + " flag. Example value of 't" + strings.Repeat(words.WildcardChar, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&WordPattern, WordPatternFlag, WordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + words.WildcardChar + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + words.WildcardChar + "' in the -" + WordPatternFlag + " flag."
	flag.StringVar(&WildcardLetters, WildcardFlag, WildcardLetters, wildcardHelp)
	flag.StringVar(&ExcludedLetters, ExcludedFlag, ExcludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	flag.StringVar(&WordFile, FileFlag, WordFile, "Word File: Name/Path of ASCII text file containing one word per line.")
	for disSpace := 0; disSpace < words.MaxLetters; disSpace++ {
		noParkDisSpaceHelp := "Letters that don't belong in this position: Letters that appear in the word, but not in postion #" + fmt.Sprintf("%d", disSpace+1) + " Example value of 'ies' means that 'i', 'e', or 's' cannot appear in position #" + fmt.Sprintf("%d", disSpace+1) + "."
		flag.StringVar(&NoParkDisSpace[disSpace], fmt.Sprintf("%d", disSpace+1), NoParkDisSpace[disSpace], noParkDisSpaceHelp)
	}
	flag.BoolVar(&PrintStatistics, StaisticsFlag, PrintStatistics, "Print statistics of letter distribution for each letter position.")
	flag.IntVar(&MaxWordsToPrint, MaxWordsToPrintFlag, MaxWordsToPrint, "Max Words to Print.")
	flag.Parse()
}

func initialize() ([]string, []string) {
	parseFlags()

	fmt.Printf("Word length: %d\n", WordLength)
	if WordPattern == "" {
		WordPattern = strings.Repeat(words.WildcardChar, WordLength)
	}
	WordPattern = strings.ToLower(WordPattern)
	fmt.Printf("Word pattern: '%s'\n", WordPattern)
	WildcardLetters = strings.ToLower(WildcardLetters)
	fmt.Printf("Wild Card letters: '%s'\n", WildcardLetters)
	ExcludedLetters = strings.ToLower(ExcludedLetters)
	fmt.Printf("Excluded letters: '%s'\n", ExcludedLetters)

	for disSpace := 0; disSpace < words.MaxLetters; disSpace++ {
		if len(NoParkDisSpace[disSpace]) > 0 {
			NoParkDisSpace[disSpace] = strings.ToLower(NoParkDisSpace[disSpace])
			fmt.Printf("Can't use letters in postion #%d: '%s'\n", disSpace+1, NoParkDisSpace[disSpace])
		}
	}

	if WordLength <= 1 {
		fmt.Println("WordLength must be greater than 0.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(WordPattern) != WordLength {
		fmt.Printf("WordPattern must be %d letters long.\n", WordLength)
		flag.PrintDefaults()
		os.Exit(1)
	}

	if WordLength == 5 && WordFile == "" {
		fmt.Println("Using built-in wordle words.")
		return words.WordleSolutionWords, append(words.WordleSolutionWords, words.WordleSearchWords...)
	} else if WordFile != "" {
		fmt.Printf("Reading Word file: %s\n", WordFile)
		solutionWords := []string{}

		f, err := os.Open(WordFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			if len(word) == WordLength {
				solutionWords = append(solutionWords, word)
			}
		}
		return solutionWords, solutionWords
	} else {
		fmt.Println("You must specify a -f <Word File>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	return nil, nil
}

func printWords(words []string, description string, maxToPrint int) {
	if len(words) == 0 {
		fmt.Printf("\nNo %s!\n", description)
		return
	}

	if len(words) == 1 {
		fmt.Printf("\n%s - EXACT MATCH! - '%s'\n", description, words[0])
		return
	}

	fmt.Printf("\n%s (%d):\n", description, len(words))
	if len(words) > maxToPrint {
		fmt.Printf("Only printing first %d\n", maxToPrint)
	}
	lineLength := 0
	sort.Strings(words)
	for i, word := range words {
		fmt.Print(word)
		lineLength += len(word) + 1
		if lineLength+len(word) > 80 {
			fmt.Println()
			lineLength = 0
		} else {
			fmt.Print(" ")
		}
		if i == maxToPrint {
			break
		}
	}
	fmt.Println()
}

func printLettersToTry(letters map[string]int) {
	if len(letters) == 0 {
		fmt.Println("\nNo additional letters to try!")
		return
	}

	fmt.Printf("\nTry these letters (%d):\n", len(letters))
	// Sort letters in order of most occurances first.
	keys := make([]string, 0, len(letters))
	for k := range letters {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return letters[keys[i]] > letters[keys[j]]
	})

	for _, k := range keys {
		fmt.Printf("%s=%d ", k, letters[k])
	}
	fmt.Println()
}

func printWordStatistics(letterDistribution []map[string]int, wordLength int) {
	fmt.Println()
	if len(letterDistribution) == 0 {
		fmt.Println("\nNo statistics to print!")
		return
	}

	for position := 0; position < wordLength; position++ {
		keys := make([]string, 0, len(letterDistribution[position]))
		for k := range letterDistribution[position] {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return letterDistribution[position][keys[i]] > letterDistribution[position][keys[j]]
		})

		fmt.Printf("Letter distribution for position #%d:\n", position+1)
		for _, k := range keys {
			if letterDistribution[position][k] > 0 {
				fmt.Printf("%s=%d ", k, letterDistribution[position][k])
			}
		}
		fmt.Println()
	}

}

func main() {
	solutionWords, allWords := initialize()

	matchingWords := words.GetMatchingWords(solutionWords, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
	printWords(matchingWords, "MATCHING WORDS", MaxWordsToPrint)
	remainingLetterCount, remainingLetterOrder := words.GetLetterCount(matchingWords, WordPattern, WildcardLetters)
	remainingLetterDistribution := words.GetLetterDistribution(matchingWords, WordLength)
	printLettersToTry(remainingLetterCount)
	if PrintStatistics {
		printWordStatistics(remainingLetterDistribution, WordLength)
	}
	if len(remainingLetterCount) > 0 {
		eliminationWords := words.GetEliminationWords(remainingLetterOrder, allWords, WordLength, ExcludedLetters, WildcardLetters, NoParkDisSpace)
		printWords(eliminationWords, "ELIMINATION WORDS", MaxWordsToPrint)
		if len(eliminationWords) > 1 {
			bestEliminationWords := words.GetBestEliminationWords(eliminationWords, WordLength, remainingLetterOrder, remainingLetterDistribution)
			printWords(bestEliminationWords, "BEST ELIMINATION WORDS", MaxWordsToPrint)
		}
	}
	fmt.Println()
}

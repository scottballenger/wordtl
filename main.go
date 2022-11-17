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
	WordFileFlag        = "f"
	StatisticsFlag      = "s"
	MaxWordsToPrintFlag = "m"
	GeussFlag           = "g"
	AnswerFlag          = "a"
	BibleWordleFlag     = "b"
)

var (
	WordLength      = words.WordleLength
	MinWordLength   = 3
	WordPattern     = "" // Pattern to match, length must be WordLength number of letters.
	WordFile        = "" // Name/Path of text file containing 1 word per line.
	WildcardLetters = "" // Letters that can appear in any position where there is a wildecard placeholder.
	ExcludedLetters = "" // Letters that cannot appear in the position where they are specified.
	NoParkDisSpace  [words.MaxLetters]string
	PrintStatistics = false
	MaxWordsToPrint = 100
	Guess           = ""
	Answer          = ""
	WordleTitle     = "Wordle"
	DoWordle        = true
	BibleWordle     = false
	DoBibleWordle   = false
)

func parseFlags() {
	flag.IntVar(&WordLength, WordLengthFlag, WordLength, "Word Length: Number of letters to match. Wordle is 5 letters.")
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + words.WildcardChar + "' 1) must include all letters specified by the -" + WildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + ExcludedFlag + " flag. Example value of 't" + strings.Repeat(words.WildcardChar, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&WordPattern, WordPatternFlag, WordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + words.WildcardChar + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + words.WildcardChar + "' in the -" + WordPatternFlag + " flag."
	flag.StringVar(&WildcardLetters, WildcardFlag, WildcardLetters, wildcardHelp)
	flag.StringVar(&ExcludedLetters, ExcludedFlag, ExcludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	wordFileHelp := "OPTIONAL Word File: Name/Path of ASCII text file containing one word per line. Will use the Wordle list from https://www.nytimes.com/games/wordle/index.html (or https://www.thelivingwordle.com if -" + BibleWordleFlag + " is specified) if this flag is not specified."
	flag.StringVar(&WordFile, WordFileFlag, WordFile, wordFileHelp)
	for disSpace := 0; disSpace < words.MaxLetters; disSpace++ {
		noParkDisSpaceHelp := "Letters that don't belong in this position: Letters that appear in the word, but not in postion #" + fmt.Sprintf("%d", disSpace+1) + " Example value of '-" + fmt.Sprintf("%d", disSpace+1) + " ies' means that 'i', 'e', or 's' cannot appear in position #" + fmt.Sprintf("%d", disSpace+1) + "."
		flag.StringVar(&NoParkDisSpace[disSpace], fmt.Sprintf("%d", disSpace+1), NoParkDisSpace[disSpace], noParkDisSpaceHelp)
	}
	flag.BoolVar(&BibleWordle, BibleWordleFlag, BibleWordle, "Use Bible Wordle words from https://www.thelivingwordle.com.")
	flag.BoolVar(&PrintStatistics, StatisticsFlag, PrintStatistics, "Print statistics of letter distribution for each letter position.")
	guessHelp := "Guess: This is your guess. Please include an Answer (-" + AnswerFlag + ") to filter the next guess. REQUIRED if -" + AnswerFlag + " is included."
	flag.StringVar(&Guess, GeussFlag, Guess, guessHelp)
	answerHelp := "Answer: Enter the following characters for each letter in your guess - '" + words.MatchedChar + "' for matching characters, '" + words.WildcardChar + "' for matching characters that are in the wwrong location, '" + words.MissedChar + "' for non-matching characters. Example value of 'X-X=X' would be match for 4th character; non-match for 1st, 3rd, and 5th character; and 2nd character is in word, but not in the 2nd position."
	flag.StringVar(&Answer, AnswerFlag, Answer, answerHelp)
	flag.IntVar(&MaxWordsToPrint, MaxWordsToPrintFlag, MaxWordsToPrint, "Max Words to Print.")
	flag.Parse()
}

func initialize() ([]string, []string, string, string) {

	parseFlags()

	fmt.Printf("Word length: %d\n", WordLength)
	Guess = strings.ToLower(Guess)
	fmt.Printf("Guess:  '%s'\n", Guess)
	fmt.Printf("Answer: '%s'\n", Answer)
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

	if WordLength < MinWordLength {
		fmt.Printf("\nERROR: WordLength must be greater than %d. Entered word length is %d.\n\n", MinWordLength-1, WordLength)
		os.Exit(1)
	}

	if len(WordPattern) != WordLength {
		fmt.Printf("\nERROR: WordPattern must be %d letters long. '%s' is %d lettters.\n\n", WordLength, WordPattern, len(WordPattern))
		os.Exit(1)
	}

	if Guess != "" && len(Guess) != WordLength {
		fmt.Printf("\nERROR: Guess must be %d letters long. '%s' is %d lettters.\n\n", WordLength, Guess, len(Guess))
		os.Exit(1)
	}

	if Guess != "" && len(Guess) == WordLength && len(Answer) != WordLength {
		fmt.Printf("\nERROR: Answer must be %d letters long. '%s' is %d lettters.\n\n", WordLength, Answer, len(Answer))
		os.Exit(1)
	}

	if Answer != "" && len(Guess) != WordLength {
		fmt.Printf("\nERROR: Guess must be provided with Answer '%s'.\n\n", Answer)
		os.Exit(1)
	}

	DoWordle = (WordLength == words.WordleLength) && (WordFile == "")
	DoBibleWordle = BibleWordle && DoWordle

	if DoWordle {
		if DoBibleWordle {
			WordleTitle = "Bible " + WordleTitle
		}
		fmt.Printf("Using built-in %s words.\n", WordleTitle)

		solutionWords := []string{}
		allWords := []string{}
		if DoBibleWordle {
			for _, solutionWord := range words.BibleWordleSolutionWords {
				solutionWords = append(solutionWords, strings.ToLower(solutionWord))
			}
			allWords = append(allWords, words.BibleWordleSolutionWords...)
			allWords = append(allWords, words.BibleWordleSearchWords...)
		} else {
			for _, solutionWord := range words.WordleSolutionWords {
				solutionWords = append(solutionWords, strings.ToLower(solutionWord))
			}
			allWords = append(allWords, words.WordleSolutionWords...)
			allWords = append(allWords, words.WordleSearchWords...)
		}

		// Remove duplicates in all words.
		updatedWords := []string{}
		visited := make(map[string]bool)
		for _, word := range allWords {
			word = strings.ToLower(word)
			if visited[word] {
				continue
			} else {
				updatedWords = append(updatedWords, word)
				visited[word] = true
			}
		}

		return solutionWords, updatedWords, Guess, Answer

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
		if len(solutionWords) == 0 {
			fmt.Printf("\nERROR: You must specify a -f <Word File> that includes %d letter words.\n\n", WordLength)
			os.Exit(1)
		}
		return solutionWords, solutionWords, Guess, Answer
	} else {
		fmt.Printf("\nERROR: You must specify a -f <Word File> for %d letter words.\n\n", WordLength)
		flag.PrintDefaults()
		os.Exit(1)
	}
	return nil, nil, "", ""
}

func printWords(words []string, description string, exclamation string, maxToPrint int) {
	if len(words) == 0 {
		fmt.Printf("\nNo %s!\n", description)
		return
	}

	if len(words) == 1 {
		fmt.Printf("\n%s - %s! - '%s'\n", description, exclamation, words[0])
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

func getWordSolutions(solutionWords []string, allWords []string) ([]string, []string, []string) {
	eliminationWords := []string{}
	bestEliminationWords := []string{}

	matchingWords := words.GetMatchingWords(solutionWords, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
	printWords(matchingWords, "MATCHING WORDS", "EXACT MATCH", MaxWordsToPrint)
	if len(matchingWords) > 1 {
		remainingLetterCount, remainingLetterOrder := words.GetLetterCount(matchingWords, WordPattern, WildcardLetters)
		remainingLetterDistribution := words.GetLetterDistribution(matchingWords, WordLength)
		printLettersToTry(remainingLetterCount)
		if PrintStatistics {
			printWordStatistics(remainingLetterDistribution, WordLength)
		}
		if len(remainingLetterOrder) > 0 {
			eliminationWords = words.GetEliminationWords(remainingLetterOrder, allWords, WordLength, ExcludedLetters, WildcardLetters, NoParkDisSpace)
			printWords(eliminationWords, "ELIMINATION WORDS", "BEST CHOICE", MaxWordsToPrint)
			if len(eliminationWords) > 1 {
				bestEliminationWords = words.GetBestEliminationWords(matchingWords, eliminationWords, WordLength, remainingLetterOrder, remainingLetterCount)
				printWords(bestEliminationWords, "BEST ELIMINATION WORDS", "BEST CHOICE", MaxWordsToPrint)
			}
		}
	}
	return matchingWords, eliminationWords, bestEliminationWords
}

func getBestGuess(matchingWords []string, eliminationWords []string, bestEliminationWords []string) string {
	if len(matchingWords) == 1 || len(matchingWords) == 2 {
		return matchingWords[0]
	} else if len(bestEliminationWords) > 0 {
		return bestEliminationWords[0]
	} else if len(eliminationWords) > 0 {
		return eliminationWords[0]
	} else if len(matchingWords) > 1 {
		return matchingWords[0]
	}
	return ""
}

func printNextGuess(
	guess string,
	wordPattern string,
	wildcardLetters string,
	excludedLetters string,
	noParkDisSpace [words.MaxLetters]string) {

	noParkDeesSpaces := ""
	for i, disSpace := range noParkDisSpace {
		if len(disSpace) > 0 {
			noParkDeesSpaces += fmt.Sprintf("-%d %s ", i+1, disSpace)
		}
	}

	bibleWordleArgs := ""
	if DoWordle && DoBibleWordle {
		bibleWordleArgs = "-" + BibleWordleFlag + " "
	}

	wordPatternArgs := ""
	if len(wordPattern) > 0 {
		wordPatternArgs = "-" + WordPatternFlag + " " + wordPattern + " "
	}

	wildcardLettersArgs := ""
	if len(wildcardLetters) > 0 {
		wildcardLettersArgs = "-" + WildcardFlag + " " + wildcardLetters + " "
	}

	excludedLettersArgs := ""
	if len(excludedLetters) > 0 {
		excludedLettersArgs = "-" + ExcludedFlag + " " + excludedLetters + " "
	}

	guessArgs := ""
	if len(guess) > 0 {
		guessArgs = "-" + GeussFlag + " " + guess + " "
	}

	wordLengthArg := ""
	if WordLength != words.WordleLength {
		wordLengthArg = "-" + WordLengthFlag + " " + fmt.Sprintf("%d", WordLength) + " "
	}

	wordFileArg := ""
	if WordFile != "" {
		wordFileArg = "-" + WordFileFlag + " " + WordFile + " "
	}

	fmt.Printf("\nTry:\n%s %s%s%s%s%s%s%s%s\n", os.Args[0], bibleWordleArgs, wordFileArg, wordLengthArg, wordPatternArgs, wildcardLettersArgs, noParkDeesSpaces, excludedLettersArgs, guessArgs)
}

func main() {
	solutionWords, allWords, guess, answer := initialize()

	WordPattern, WildcardLetters, ExcludedLetters, NoParkDisSpace = words.TranslateGuessResults(guess, answer, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
	matchingWords, eliminationWords, bestEliminationWords := getWordSolutions(solutionWords, allWords)
	guess = getBestGuess(matchingWords, eliminationWords, bestEliminationWords)
	printNextGuess(guess, WordPattern, WildcardLetters, ExcludedLetters, NoParkDisSpace)
	fmt.Println()
}

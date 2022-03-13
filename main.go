package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
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
	AutoGuessFlag       = "a"
	GeussFlag           = "g"
	WordleDayFlag       = "d"
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
	AutoGuess       = false
	Guess           = ""
	WordleDay       = 1
	TodaysDay       = 1
	DoWordle        = true
)

func parseFlags() {
	flag.IntVar(&WordLength, WordLengthFlag, WordLength, "Word Length: Number of letters to match. wordle is 5 letters.")
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + words.WildcardChar + "' 1) must include all letters specified by the -" + WildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + ExcludedFlag + " flag. Example value of 't" + strings.Repeat(words.WildcardChar, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	flag.StringVar(&WordPattern, WordPatternFlag, WordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + words.WildcardChar + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + words.WildcardChar + "' in the -" + WordPatternFlag + " flag."
	flag.StringVar(&WildcardLetters, WildcardFlag, WildcardLetters, wildcardHelp)
	flag.StringVar(&ExcludedLetters, ExcludedFlag, ExcludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	flag.StringVar(&WordFile, WordFileFlag, WordFile, "OPTIONAL Word File: Name/Path of ASCII text file containing one word per line. Will use the wordle list from https://www.nytimes.com/games/wordle/index.html if this flag is not specified.")
	for disSpace := 0; disSpace < words.MaxLetters; disSpace++ {
		noParkDisSpaceHelp := "Letters that don't belong in this position: Letters that appear in the word, but not in postion #" + fmt.Sprintf("%d", disSpace+1) + " Example value of '-" + fmt.Sprintf("%d", disSpace+1) + " ies' means that 'i', 'e', or 's' cannot appear in position #" + fmt.Sprintf("%d", disSpace+1) + "."
		flag.StringVar(&NoParkDisSpace[disSpace], fmt.Sprintf("%d", disSpace+1), NoParkDisSpace[disSpace], noParkDisSpaceHelp)
	}
	flag.BoolVar(&PrintStatistics, StatisticsFlag, PrintStatistics, "Print statistics of letter distribution for each letter position.")
	flag.IntVar(&WordleDay, WordleDayFlag, WordleDay, "Wordle Day: The # of the Wordle solution to use.")
	flag.BoolVar(&AutoGuess, AutoGuessFlag, AutoGuess, "Try to guess the word by iterating through guesses.")
	flag.StringVar(&Guess, GeussFlag, Guess, "Guess: This is your guess.")
	flag.IntVar(&MaxWordsToPrint, MaxWordsToPrintFlag, MaxWordsToPrint, "Max Words to Print.")
	flag.Parse()
}

func initialize() ([]string, []string, string, string) {

	rand.Seed(time.Now().UnixNano())

	timeFormat := "2006-01-02"
	t, _ := time.Parse(timeFormat, words.StartDate)
	year, month, day := time.Now().Date()
	todaysDate := fmt.Sprintf("%d-%02d-%02d", year, int(month), day)
	now, _ := time.Parse(timeFormat, todaysDate)
	WordleDay = int(now.Sub(t).Hours()/24) + 1
	TodaysDay = WordleDay
	fmt.Printf("Todays Wordle Day: %d\n", WordleDay)

	parseFlags()

	fmt.Printf("Word length: %d\n", WordLength)
	fmt.Printf("Guess: '%s'\n", Guess)
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

	if WordleDay < 1 {
		fmt.Printf("Wordle Day must be greater than 0. You entered %d, using 1.\n", WordleDay)
		WordleDay = 1
	}

	DoWordle := (WordLength == words.WordleLength) && (WordFile == "")

	if DoWordle {
		fmt.Println("Using built-in wordle words.")
		maxWordleDay := len(words.WordleSolutionWords)
		if WordleDay > maxWordleDay {
			fmt.Printf("Wordle Day must be less than %d. You entered %d, using %d.\n", maxWordleDay+1, WordleDay, maxWordleDay)
			WordleDay = maxWordleDay
		}
		fmt.Printf("Yesterday's wordle: '%s'\n", words.GetWordle(WordleDay-1))
		fmt.Printf("Solving wordle for Day: %d\n", WordleDay)
		return words.WordleSolutionWords, append(words.WordleSolutionWords, words.WordleSearchWords...), words.GetWordle(WordleDay), Guess
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
		if WordleDay == TodaysDay {
			WordleDay = rand.Intn(len(solutionWords) - 1)
		}
		fmt.Printf("Solving %d letter word from %s for Day: %d\n", WordLength, WordFile, WordleDay)
		return solutionWords, solutionWords, solutionWords[WordleDay], Guess
	} else {
		fmt.Printf("\nERROR: You must specify a -f <Word File> for %d letter words.\n\n", WordLength)
		flag.PrintDefaults()
		os.Exit(1)
	}
	return nil, nil, "", ""
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

func getWordSolutions(solutionWords []string, allWords []string) ([]string, []string, []string) {
	eliminationWords := []string{}
	bestEliminationWords := []string{}

	matchingWords := words.GetMatchingWords(solutionWords, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
	printWords(matchingWords, "MATCHING WORDS", MaxWordsToPrint)
	if len(matchingWords) > 1 {
		remainingLetterCount, remainingLetterOrder := words.GetLetterCount(matchingWords, WordPattern, WildcardLetters)
		remainingLetterDistribution := words.GetLetterDistribution(matchingWords, WordLength)
		printLettersToTry(remainingLetterCount)
		if PrintStatistics {
			printWordStatistics(remainingLetterDistribution, WordLength)
		}
		if len(remainingLetterOrder) > 0 {
			eliminationWords = words.GetEliminationWords(remainingLetterOrder, allWords, WordLength, ExcludedLetters, WildcardLetters, NoParkDisSpace)
			printWords(eliminationWords, "ELIMINATION WORDS", MaxWordsToPrint)
			if len(eliminationWords) > 1 {
				bestEliminationWords = words.GetBestEliminationWords(eliminationWords, WordLength, remainingLetterOrder, remainingLetterDistribution)
				printWords(bestEliminationWords, "BEST ELIMINATION WORDS", MaxWordsToPrint)
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

	dayArg := ""
	if WordleDay != TodaysDay {
		dayArg = "-" + WordleDayFlag + " " + fmt.Sprintf("%d", WordleDay) + " "
	}

	wordLengthArg := ""
	if WordLength != words.WordleLength {
		wordLengthArg = "-" + WordLengthFlag + " " + fmt.Sprintf("%d", WordLength) + " "
	}

	wordFileArg := ""
	if WordFile != "" {
		wordFileArg = "-" + WordFileFlag + " " + WordFile + " "
	}

	fmt.Printf("\nTry:\n%s %s%s%s%s%s%s%s%s\n", os.Args[0], wordFileArg, wordLengthArg, dayArg, wordPatternArgs, wildcardLettersArgs, noParkDeesSpaces, excludedLettersArgs, guessArgs)
}

func main() {
	solutionWords, allWords, word, guess := initialize()

	if len(guess) > 0 {
		matchingGuessWords := words.GetMatchingWords(allWords, guess, "", "", [words.MaxLetters]string{})
		if len(matchingGuessWords) != 1 {
			fmt.Printf("\nERROR: Guess '%s' is not in the dictionary.\n\n", guess)
			os.Exit(1)
		}
	}

	numTries := 1
	if AutoGuess {
		numTries = 6
	}
	for try := 1; try <= numTries; try++ {
		if AutoGuess {
			fmt.Printf("\nAttempt #%d:\n", try)
		}
		if guess == "" {
			matchingWords, eliminationWords, bestEliminationWords := getWordSolutions(solutionWords, allWords)
			guess = getBestGuess(matchingWords, eliminationWords, bestEliminationWords)
			if !AutoGuess {
				printNextGuess(guess, WordPattern, WildcardLetters, ExcludedLetters, NoParkDisSpace)
				guess = ""
			}
		}

		if guess != "" {
			guessed, answer := words.GuessWord(word, guess)
			fmt.Printf("\nGuessed '%s', match=%v\n", guess, guessed)
			fmt.Printf("         %s\n", answer)
			if guessed {
				guess = ""
				break
			}
			WordPattern, WildcardLetters, ExcludedLetters, NoParkDisSpace = words.TranslateGuessResults(guess, answer, WordPattern, ExcludedLetters, WildcardLetters, NoParkDisSpace)
			matchingWords, eliminationWords, bestEliminationWords := getWordSolutions(solutionWords, allWords)
			guess = getBestGuess(matchingWords, eliminationWords, bestEliminationWords)
			printNextGuess(guess, WordPattern, WildcardLetters, ExcludedLetters, NoParkDisSpace)
		} else {
			break
		}
	}
	if AutoGuess && guess != "" {
		fmt.Printf("\nDid not guess word after %d tries, next guess is '%s'\n", numTries, guess)
	}
	fmt.Println()
}

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
	"wordtl/words"

	"github.com/gookit/color"
)

const (
	ExcludeAllFlag                = "exclude-all"
	ExcludeByPosFlag              = "exclude-pos"
	WordFileFlag                  = "file"
	GuessFlag                     = "guess"
	ResultFlag                    = "guess-result"
	IgnoreWordleSolutionWordsFlag = "ignore-wordle-solution-words"
	IgnoreWordleUsedWordsFlag     = "ignore-wordle-used-words"
	WordLengthFlag                = "length"
	MaxWordsToPrintFlag           = "max-print"
	NotInPosFlag                  = "not"
	WordPatternFlag               = "pattern"
	DiagnosticsFlag               = "stats"
	UseWordleSolutionWordsFlag    = "use-wordle-solution-words"
	UseWordleUsedWordsFlag        = "use-wordle-used-words"
	WildcardFlag                  = "wildcards"

	MaxTries = 6

	ModeAutoPlay    = "auto"
	ModeManualGuess = "manual"
	ModeWordSearch  = "search"
	ModeHelp        = "help"
)

var (
	WordLength       = words.WordleLength
	MinWordLength    = 3
	WordPattern      = "" // Pattern to match, length must be WordLength number of letters.
	WordFile         = "" // Name/Path of text file containing 1 word per line.
	WildcardLetters  = "" // Letters that can appear in any position where there is a wildecard placeholder.
	ExcludedByPosStr = "" // Letters that cannot appear in the position where they are specified.
	ExcludedByPosMap map[int]string
	ExcludedLetters  = "" // Letters that cannot appear anywhere in the word.
	PrintDiagnostics = false
	MaxWordsToPrint  = 100
	Guess            = ""
	Result           = ""
	WordleTitle      = "Wordle"
	DoWordle         = true
	Mode             = ModeAutoPlay
	TodaysDay        = 1
	UsedWordsFile    = "words/wordle_words_used.txt"

	IgnoreWordleSolutionWords = false
	IgnoreWordleUsedWords     = false
)

func parseFlags() {
	wordleCmd := flag.NewFlagSet(ModeAutoPlay, flag.ExitOnError)
	guessCmd := flag.NewFlagSet(ModeManualGuess, flag.ExitOnError)
	searchCmd := flag.NewFlagSet(ModeWordSearch, flag.ExitOnError)
	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	subcommands := map[string]*flag.FlagSet{
		wordleCmd.Name(): wordleCmd,
		guessCmd.Name():  guessCmd,
		searchCmd.Name(): searchCmd,
	}

	// Manual Guess Flags
	addSearchFlags(guessCmd)
	guessHelp := "Guess: This is your guess. Please include a Result (-" + ResultFlag + ") to filter the next guess. REQUIRED if -" + ResultFlag + " is included."
	guessCmd.StringVar(&Guess, GuessFlag, Guess, guessHelp)
	resultHelp := "Result: Enter the following characters for each letter in your guess - '" + words.MatchedChar + "' for matching characters, '" + words.WildcardChar + "' for matching characters that are in the wrong location, '" + words.MissedChar + "' for non-matching characters. Example value of '" + words.MissedChar + words.WildcardChar + words.MissedChar + words.MatchedChar + words.MissedChar + "' would be match for 4th character; non-match for 1st, 3rd, and 5th character; and 2nd character is in word, but not in the 2nd position."
	guessCmd.StringVar(&Result, ResultFlag, Result, resultHelp)

	// Search Flags
	addSearchFlags(searchCmd)
	useWordleSolutionWords := false
	useWordleUsedWords := false

	// Global Flags
	for _, fs := range subcommands {
		fs.IntVar(&WordLength, WordLengthFlag, WordLength, "Word Length: Number of letters in each word. Wordle is 5 letters.")
		wordFileHelp := "OPTIONAL Word File: Name/Path of ASCII text file containing one word per line. Will use the Wordle list from https://www.nytimes.com/games/wordle/index.html if this flag is not specified."
		fs.StringVar(&WordFile, WordFileFlag, WordFile, wordFileHelp)
		if fs == searchCmd {
			fs.BoolVar(&useWordleSolutionWords, UseWordleSolutionWordsFlag, useWordleSolutionWords, "Consider Wordle solution words for results. Can ony be used when -"+WordFileFlag+" is NOT specified.")
			fs.BoolVar(&useWordleUsedWords, UseWordleUsedWordsFlag, useWordleUsedWords, "Consider previously used Wordle solution words for results. Can ony be used when -"+WordFileFlag+" is NOT specified.")
		} else {
			fs.BoolVar(&IgnoreWordleSolutionWords, IgnoreWordleSolutionWordsFlag, IgnoreWordleSolutionWords, "Do not consider Wordle solution words for results. Can ony be used when -"+WordFileFlag+" is NOT specified.")
			fs.BoolVar(&IgnoreWordleUsedWords, IgnoreWordleUsedWordsFlag, IgnoreWordleUsedWords, "Do not consider previously used Wordle solution words for results. Can ony be used when -"+WordFileFlag+" is NOT specified.")
		}
		fs.BoolVar(&PrintDiagnostics, DiagnosticsFlag, PrintDiagnostics, "Print statistics of letter distribution for each letter position.")
		fs.IntVar(&MaxWordsToPrint, MaxWordsToPrintFlag, MaxWordsToPrint, "Max Words to Print.")
	}

	if len(os.Args) < 2 {
		printUsage(subcommands, "expected subcommand")
		os.Exit(1)
	}

	if os.Args[1] == helpCmd.Name() {
		printUsage(subcommands, "")
		os.Exit(1)
	}

	cmd := subcommands[os.Args[1]]
	if cmd == nil {
		printUsage(subcommands, fmt.Sprintf("unknown subcommand '%s', see usage for known subcommands.\n", os.Args[1]))
		os.Exit(1)
	}

	cmd.Parse(os.Args[2:])
	Mode = cmd.Name()
	fmt.Println()
	fmt.Println(getModeDescription(Mode))
	fmt.Println()

	if Mode == ModeWordSearch {
		IgnoreWordleSolutionWords = !useWordleSolutionWords
		IgnoreWordleUsedWords = !useWordleUsedWords
	}
}

func addSearchFlags(fs *flag.FlagSet) {
	wordPatternHelp := "Pattern to Match: Known letters will be in the position that they appear. Wildecard placeholders '" + words.WildcardChar + "' 1) must include all letters specified by the -" + WildcardFlag + " flag and 2) can be any other letter that is not excluded by the -" + ExcludeAllFlag + " flag. Example value of 't" + strings.Repeat(words.WildcardChar, 4) + "' would lookup words with a 't' in the beginning of a 5 letter word."
	fs.StringVar(&WordPattern, WordPatternFlag, WordPattern, wordPatternHelp)
	wildcardHelp := "Wildcard Letters: Letters that must appear in any position where there is a wildecard placeholder '" + words.WildcardChar + "'. Example value of 'r' means that there must be at least 1 'r' in any place where there is a '" + words.WildcardChar + "' in the -" + WordPatternFlag + " flag."
	fs.StringVar(&WildcardLetters, WildcardFlag, WildcardLetters, wildcardHelp)
	fs.StringVar(&ExcludedLetters, ExcludeAllFlag, ExcludedLetters, "Excluded Letters: Letters that cannot appear in the word. Example value of 'ies' means that 'i', 'e', or 's' cannot appear anywhere in the word.")
	fs.StringVar(&ExcludedByPosStr, ExcludeByPosFlag, ExcludedByPosStr, "Excluded Letters by Position: Letters that cannot appear in a specific position of the word. JSON - Example value of '{\"1\":\"ab\",\"4\":\"cd\"}' means that 'a' or 'b' cannot appear in position #1 of the word and 'c' or 'd' cannot appear in position #4 of the word. Position must be an integer greater than or equal to 1 and should be less than or equal to the word length.")
}

func printUsage(subcommandFlagset map[string]*flag.FlagSet, errString string) {
	fmt.Println()
	if len(errString) > 0 {
		fmt.Println("[ERROR] " + errString)
		fmt.Println()
	}
	fmt.Println("usage: " + os.Args[0] + " subcommmand [flags]")
	fmt.Println()
	fmt.Println("Available subcommands:")
	subcommandSpacing := len(ModeHelp)
	var subcommands []string
	for _, fs := range subcommandFlagset {
		subcommand := fs.Name()
		subcommands = append(subcommands, subcommand)
		if len(subcommand) > subcommandSpacing {
			subcommandSpacing = len(subcommand)
		}
	}
	subcommands = append(subcommands, ModeHelp)
	subcommandSpacing += 3

	for _, subcommand := range subcommands {
		fmt.Println("   " + subcommand + strings.Repeat(" ", subcommandSpacing-len(subcommand)) + getModeDescription(subcommand))
	}
	fmt.Println()
	fmt.Println("For specific subcommand flags, enter '" + os.Args[0] + " subcommmand -h'")
	fmt.Println()
	fmt.Println("Also see https://github.com/scottballenger/wordtl/blob/main/README.md for a detailed description.")
	fmt.Println()
}

func getModeDescription(mode string) string {
	switch mode {
	case ModeAutoPlay:
		return "Auto Play: Try to guess the word in " + fmt.Sprintf("%d", MaxTries) + " tries"
	case ModeManualGuess:
		return "Manual Guess: Get help with a single guess"
	case ModeWordSearch:
		return "Search All Words: dictionary lookup"
	case ModeHelp:
		return "Print subcommand help message"
	default:
		return ""
	}
}

func initialize() ([]string, []string, map[string]bool, string, string) {

	parseFlags()

	fmt.Printf("Word length: %d\n", WordLength)
	if len(WordPattern) > 0 {
		WordPattern = strings.ToLower(WordPattern)
		fmt.Printf("Word pattern: '%s'\n", WordPattern)
	} else {
		WordPattern = strings.Repeat(words.WildcardChar, WordLength)
	}
	if len(WildcardLetters) > 0 {
		WildcardLetters = strings.ToLower(WildcardLetters)
		fmt.Printf("Wild Card letters: '%s'\n", WildcardLetters)
	}
	if len(ExcludedLetters) > 0 {
		ExcludedLetters = strings.ToLower(ExcludedLetters)
		fmt.Printf("Excluded letters: '%s'\n", ExcludedLetters)
	}
	ExcludedByPosMap = make(map[int]string)
	if len(ExcludedByPosStr) > 0 {
		if err := json.Unmarshal([]byte(ExcludedByPosStr), &ExcludedByPosMap); err != nil {
			fmt.Println("Invalid JSON for -" + ExcludeByPosFlag + " '" + ExcludedByPosStr + "', see usage with '" + os.Args[0] + " " + Mode + " -h'")
			os.Exit(1)
		}
		ints := make([]int, 0, len(ExcludedByPosMap))
		for pos, letters := range ExcludedByPosMap {
			letters = strings.ToLower(letters)
			ExcludedByPosMap[pos] = letters
			ints = append(ints, pos)
		}
		sort.Ints(ints[:])
		for _, pos := range ints {
			if pos > 0 {
				cantUseError := ""
				if pos > WordLength {
					cantUseError = " [Invalid due to word length of " + fmt.Sprintf("%d", WordLength) + "]"
				}
				fmt.Printf("Can't use letters in postion #%d: '%s'%s\n", pos, ExcludedByPosMap[pos], cantUseError)
			} else {
				fmt.Println("Position " + fmt.Sprintf("%d", pos) + " out of range for -" + ExcludeByPosFlag + " '" + ExcludedByPosStr + "', see usage with '" + os.Args[0] + " " + Mode + " -h'")
				os.Exit(1)
			}
		}
	}

	Guess = strings.ToLower(Guess)
	if Mode == ModeManualGuess {
		fmt.Printf("Guess:  '%s'\n", Guess)
		fmt.Printf("Result: '%s'\n", Result)
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

	if Guess != "" && len(Guess) == WordLength && len(Result) != WordLength {
		fmt.Printf("\nERROR: Result must be %d letters long. '%s' is %d lettters.\n\n", WordLength, Result, len(Result))
		os.Exit(1)
	}

	if Result != "" && len(Guess) != WordLength {
		fmt.Printf("\nERROR: Guess must be provided with Result '%s'.\n\n", Result)
		os.Exit(1)
	}

	DoWordle = (WordLength == words.WordleLength) && (WordFile == "")

	if DoWordle {
		fmt.Printf("Using built-in %s words.\n", WordleTitle)

		if (Mode == ModeAutoPlay || Mode == ModeManualGuess) && !IgnoreWordleUsedWords {
			startDate := words.WordleStartDate
			timeFormat := "2006-01-02"
			t, _ := time.Parse(timeFormat, startDate)
			year, month, day := time.Now().Date()
			todaysDate := fmt.Sprintf("%d-%02d-%02d", year, int(month), day)
			now, _ := time.Parse(timeFormat, todaysDate)
			TodaysDay = int(now.Sub(t).Hours() / 24)
			fmt.Printf("Todays %s Day: %d\n", WordleTitle, TodaysDay)
		}

		usedWords := make(map[string]bool)
		solutionWords := []string{}
		if IgnoreWordleSolutionWords {
			fmt.Printf("Ignoring built-in %s solution words.\n", WordleTitle)
		} else {
			fmt.Printf("Using built-in %s solution words.\n", WordleTitle)
			if IgnoreWordleUsedWords {
				solutionWords = words.WordleSolutionWords
			} else {
				// Wordle already used words are contained in a separate file.
				fmt.Printf("Removing previously used %s solution words.\n", WordleTitle)
				f, err := os.Open(UsedWordsFile)
				if err != nil {
					log.Println(err)
				} else {
					defer f.Close()
					scanner := bufio.NewScanner(f)
					for scanner.Scan() {
						word := strings.ToLower(scanner.Text())
						if len(word) == WordLength {
							usedWords[word] = true
						}
					}
				}

				// Remove the used words from the solution words.
				for _, word := range words.WordleSolutionWords {
					word = strings.ToLower(word)
					if !usedWords[word] {
						solutionWords = append(solutionWords, word)
					}
				}
			}
		}

		allWords := append(words.WordleSolutionWords, words.WordleSearchWords...)
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
		allWords = updatedWords

		if len(solutionWords) == 0 {
			solutionWords = allWords
		}

		return solutionWords, allWords, usedWords, Guess, Result

	} else if WordFile != "" {
		fmt.Printf("Reading Word file: %s\n", WordFile)
		allWords := []string{}

		f, err := os.Open(WordFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			word := strings.ToLower(scanner.Text())
			if len(word) == WordLength {
				allWords = append(allWords, word)
			}
		}
		if len(allWords) == 0 {
			fmt.Printf("\nERROR: '%s' does NOT include any %d letter words.\n\n", WordFile, WordLength)
			os.Exit(1)
		}
		return allWords, allWords, nil, Guess, Result
	} else {
		fmt.Printf("\nERROR: You must specify a -f <Word File> for %d letter words.\n\n", WordLength)
		flag.PrintDefaults()
		os.Exit(1)
	}
	return nil, nil, nil, "", ""
}

func printWords(words []string, description string, exclamation string, maxToPrint int) {
	if len(words) == 0 {
		fmt.Printf("\nNo %s!\n", description)
		return
	}

	if len(words) == 1 && len(exclamation) > 0 {
		fmt.Printf("\n%s - %s! - '%s'\n", description, exclamation, words[0])
		return
	}

	fmt.Printf("\n%s (%d):\n", description, len(words))
	if len(words) > maxToPrint {
		fmt.Printf("Only printing first %d\n", maxToPrint)
	}
	lineLength := 0
	sortedWords := []string{}
	sortedWords = append(sortedWords, words...) // Create a copy so sort does not disturb the original array.
	sort.Strings(sortedWords)
	for i, word := range sortedWords {
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

func printWordDiagnostics(letterDistribution []map[string]int, wordLength int) {
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

	matchingWords := words.GetMatchingWords(solutionWords, WordPattern, ExcludedLetters, WildcardLetters, true, ExcludedByPosMap)
	printWords(matchingWords, "MATCHING WORDS", "EXACT MATCH", MaxWordsToPrint)
	if len(matchingWords) > 1 {
		remainingLetterCount, remainingLetterOrder := words.GetLetterCount(matchingWords, WordPattern, WildcardLetters)
		remainingLetterDistribution := words.GetLetterDistribution(matchingWords, WordLength)
		printLettersToTry(remainingLetterCount)
		if PrintDiagnostics {
			printWordDiagnostics(remainingLetterDistribution, WordLength)
		}
		if len(remainingLetterOrder) > 0 {
			eliminationWords = words.GetEliminationWords(remainingLetterOrder, allWords, WordLength, ExcludedLetters, WildcardLetters, ExcludedByPosMap)
			if len(eliminationWords) < 2*MaxWordsToPrint {
				printWords(eliminationWords, "ELIMINATION WORDS", "BEST CHOICE", MaxWordsToPrint)
			}
			if len(eliminationWords) > 1 {
				bestEliminationWords = words.GetBestEliminationWords(matchingWords, eliminationWords, WordLength, remainingLetterOrder, remainingLetterCount, remainingLetterDistribution)
				printWords(bestEliminationWords, "BEST ELIMINATION WORDS", "BEST CHOICE", MaxWordsToPrint)
				if len(bestEliminationWords) > 1 {
					bestEliminationWord := []string{}
					bestEliminationWord = append(bestEliminationWord, bestEliminationWords[0])
					printWords(bestEliminationWord, "BEST ELIMINATION WORD", "BEST CHOICE", MaxWordsToPrint)
				}
			}
		}
	}
	return matchingWords, eliminationWords, bestEliminationWords
}

func getBestGuess(matchingWords []string, eliminationWords []string, bestEliminationWords []string) string {
	if len(matchingWords) == 1 || len(matchingWords) == 2 {
		fmt.Println()
		fmt.Println("Using MATCHING WORD - '" + matchingWords[0] + "'")
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
	excludedByPosMap map[int]string) {

	excludedByPosStr := ""
	if len(excludedByPosMap) > 0 {
		json, _ := json.Marshal(excludedByPosMap)
		excludedByPosStr = fmt.Sprintf("-%s '%s' ", ExcludeByPosFlag, string(json))
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
		excludedLettersArgs = "-" + ExcludeAllFlag + " " + excludedLetters + " "
	}

	guessArgs := ""
	if len(guess) > 0 {
		guessArgs = "-" + GuessFlag + " " + guess + " -" + ResultFlag
	} else {
		if Mode == ModeManualGuess {
			if !IgnoreWordleSolutionWords {
				if !IgnoreWordleUsedWords {
					guessArgs = "-" + IgnoreWordleUsedWordsFlag
				} else {
					guessArgs = "-" + IgnoreWordleSolutionWordsFlag
				}
			}
		}
	}

	wordLengthArg := ""
	if WordLength != words.WordleLength {
		wordLengthArg = "-" + WordLengthFlag + " " + fmt.Sprintf("%d", WordLength) + " "
	}

	wordFileArg := ""
	if WordFile != "" {
		wordFileArg = "-" + WordFileFlag + " " + WordFile + " "
	}

	ignoreWordleSolutionWordsFlag := ""
	if IgnoreWordleSolutionWords {
		ignoreWordleSolutionWordsFlag = "-" + IgnoreWordleSolutionWordsFlag + " "
	}

	ignoreWordleUsedWordsFlag := ""
	if IgnoreWordleUsedWords {
		ignoreWordleUsedWordsFlag = "-" + IgnoreWordleUsedWordsFlag + " "
	}

	fmt.Printf("\nTry:\n%s %s %s%s%s%s%s%s%s%s%s\n", os.Args[0], Mode, wordFileArg, wordLengthArg, ignoreWordleSolutionWordsFlag, ignoreWordleUsedWordsFlag, wordPatternArgs, wildcardLettersArgs, excludedByPosStr, excludedLettersArgs, guessArgs)
}

func getUserInputRange(defaultVal string, valName string, startChar string, endChar string, validCharsMsg string, validCharsHelp string, validLength int) string {
	validChars := ""
	if startChar[0] < endChar[0] && len(startChar) == 1 && len(endChar) == 1 {
		for char := startChar[0]; char <= endChar[0]; char++ {
			validChars += string(char)
		}
		return getUserInput(defaultVal, valName, validChars, validCharsMsg, validCharsHelp, validLength)
	}
	fmt.Println()
	fmt.Println("Invalid range starting with '" + startChar + "' and ending with '" + endChar + "'.")
	fmt.Println()
	return ""
}

func getUserInput(defaultVal string, valName string, validChars string, validCharsMsg string, validCharsHelp string, validLength int) string {
	// TODO: Add unit test
	exitStr := "0"
	defaultStr := ""
	userInput := ""
	if len(defaultVal) > 0 {
		defaultStr = "default = '" + defaultVal + "', "
	}
	for {
		fmt.Print(valName, " (", defaultStr, "exit = '"+exitStr+"'): ")
		reader := bufio.NewReader(os.Stdin)
		userInput, _ = reader.ReadString('\n')
		userInput = strings.TrimSuffix(userInput, "\n")
		userInput = strings.ToLower(userInput)
		if userInput == "" {
			userInput = defaultVal
		}
		if len(userInput) == len(exitStr) && userInput == exitStr {
			os.Exit(0)
		}
		if len(userInput) == validLength {
			validInput := true
			invalidChars := ""
			for _, thisChar := range userInput {
				letter := string(thisChar)
				if !strings.Contains(strings.ToLower(validChars), letter) {
					validInput = false
					invalidChars += letter
				}
			}
			if validInput {
				break
			} else {
				fmt.Println()
				fmt.Println("Value must contain only the following characters: '" + validCharsMsg + "'" + validCharsHelp + ". Your input: '" + userInput + "' includes the following invalid characters: '" + invalidChars + "'.")
				fmt.Println()
			}
		} else {
			fmt.Println()
			fmt.Println("Value must be " + fmt.Sprintf("%d", validLength) + " characters, '" + userInput + "' is " + fmt.Sprintf("%d", len(userInput)) + " characters.")
			fmt.Println()
		}
	}

	return userInput
}

func printWordleResult(guess string, result string) bool {
	match := color.New(color.BgGreen, color.Bold)
	almost := color.New(color.BgLightYellow, color.Bold)
	miss := color.New(color.BgDarkGray, color.Bold)
	incorrect := color.New(color.BgHiRed, color.Bold)
	correctForm := true

	if len(guess) == len(result) {
		for ndx, letter := range guess {
			char := strings.ToUpper(string(letter))
			switch string(result[ndx]) {
			case words.MatchedChar:
				match.Print(" " + string(char) + " ")
			case words.MissedChar:
				miss.Print(" " + string(char) + " ")
			case words.WildcardChar:
				almost.Print(" " + string(char) + " ")
			default:
				incorrect.Print(" " + string(char) + " ")
				correctForm = false
			}
			if ndx < len(guess)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()

		if !correctForm {
			fmt.Println("Result: '" + result + "' must be in the proper form.")
		}
	} else {
		fmt.Println("Guess: '" + guess + "' and Result: '" + result + "' must be the same length.")
		correctForm = false
	}

	return correctForm
}

func printWordleSolution(guesses [MaxTries]string, results [MaxTries]string, numTries int, foundSolution bool) {
	if foundSolution {
		fmt.Println()
		fmt.Println("Congratulations, you have found the solution word in " + fmt.Sprintf("%d", numTries+1) + " turns!")
		fmt.Println()
	} else {
		if numTries+1 > 1 {
			fmt.Println()
			fmt.Println("Result after " + fmt.Sprintf("%d", numTries+1) + " guesses:")
		} else {
			return
		}
	}
	for i := 0; i <= numTries; i++ {
		printWordleResult(guesses[i], results[i])
		fmt.Println()
	}
}

func isResultCorrect(result string, validlength int) bool {
	if len(result) != validlength {
		return false
	}
	correct := true
	for i := 0; i < len(result); i++ {
		correct = correct && (string(result[i]) == words.MatchedChar)
	}
	return correct
}

func WordSearch(solutionWords []string) {
	if len(WildcardLetters) > 0 {
		letterCount := map[string]int{}
		for _, letter := range WildcardLetters {
			letterCount[string(letter)]++
		}
		letterDistribution := []map[string]int{}
		for position := 0; position < WordLength; position++ {
			letterDistribution = append(letterDistribution, map[string]int{})
			for _, letter := range WildcardLetters {
				letterDistribution[position][string(letter)]++
			}
		}

		matchingWords := words.GetMatchingWords(solutionWords, WordPattern, ExcludedLetters, "", false, ExcludedByPosMap)
		if len(matchingWords) > 1 {
			matchingWords = words.GetBestEliminationWords([]string{}, matchingWords, WordLength, WildcardLetters, letterCount, letterDistribution)
		}
		if len(matchingWords) > 0 {
			wordSearchTitle := "ALL"
			if DoWordle {
				wordSearchTitle += " " + strings.ToUpper(WordleTitle)
				if !IgnoreWordleSolutionWords {
					wordSearchTitle += " SOLUTION"
					if !IgnoreWordleUsedWords {
						wordSearchTitle += " (minus USED)"
					}
				}
			}
			printWords(matchingWords, "SEARCH "+wordSearchTitle+" WORDS", "EXACT MATCH", MaxWordsToPrint)
		} else {
			fmt.Println()
			fmt.Println("NO MATCHING WORDS. Please change args to get matching results.")
			fmt.Println()
		}
		fmt.Println()
	} else {
		wordSearchHelp := "Nothing to SEARCH! Please use the -" + WildcardFlag + " flag to specify letters to search for."
		fmt.Println()
		fmt.Println(wordSearchHelp)
		fmt.Println()
	}
}

func ManualGuess(guess string, result string, solutionWords []string, allWords []string) {
	if isResultCorrect(result, WordLength) {
		fmt.Println()
		fmt.Println("Congratulations, '" + guess + "' is the solution word!")
		fmt.Println()
	} else {
		WordPattern, WildcardLetters, ExcludedLetters, ExcludedByPosMap = words.TranslateGuessResults(guess, result, WordPattern, ExcludedLetters, WildcardLetters, ExcludedByPosMap)
		matchingWords, eliminationWords, bestEliminationWords := getWordSolutions(solutionWords, allWords)
		guess = getBestGuess(matchingWords, eliminationWords, bestEliminationWords)
		printNextGuess(guess, WordPattern, WildcardLetters, ExcludedLetters, ExcludedByPosMap)
		fmt.Println()
	}
}

func AutoPlay(guess string, result string, solutionWords []string, allWords []string, usedWords map[string]bool) {
	const (
		yes = "y"
		no  = "n"
	)

	var guesses [MaxTries]string
	var results [MaxTries]string
	for try := 0; try < MaxTries; try++ {
		WordPattern, WildcardLetters, ExcludedLetters, ExcludedByPosMap = words.TranslateGuessResults(guess, result, WordPattern, ExcludedLetters, WildcardLetters, ExcludedByPosMap)
		matchingWords, eliminationWords, bestEliminationWords := getWordSolutions(solutionWords, allWords)
		if (len(matchingWords) == 0) && (len(solutionWords) != len(allWords)) {
			fmt.Println()
			useAllWords := getUserInput(yes, "No matching words found in Solution Words.\n\nDo you want to search All Words?", yes+no, yes+" or "+no, "", 1)
			if useAllWords == yes {
				IgnoreWordleSolutionWords = true
				solutionWords = allWords
				matchingWords, eliminationWords, bestEliminationWords = getWordSolutions(solutionWords, allWords)
				var previouslyUsedWords []string
				for _, word := range matchingWords {
					if usedWords[word] {
						previouslyUsedWords = append(previouslyUsedWords, word)
					}
				}
				if len(previouslyUsedWords) > 0 {
					printWords(previouslyUsedWords, "FYI - The following words were previously used as Wordle Solutions", "", MaxWordsToPrint)
				}
			}
		}
		guess = getBestGuess(matchingWords, eliminationWords, bestEliminationWords)

		fmt.Println()
		fmt.Println("TRY #" + fmt.Sprintf("%d", try+1))
		fmt.Println("------")
		fmt.Println()
		userGuess := guess
		userResult := result
		for {
			userGuess = getUserInputRange(userGuess, "Enter your Guess", "a", "z", "a-z", "", WordLength)
			userResult = getUserInput("", "Enter your Result", words.MatchedChar+words.WildcardChar+words.MissedChar, words.MatchedChar+words.WildcardChar+words.MissedChar, " (where '"+words.MatchedChar+"' is a matching character in position, '"+words.WildcardChar+"' is a matching character out of position, and '"+words.MissedChar+"' is a non-matching character)", WordLength)
			fmt.Println()
			correctForm := printWordleResult(userGuess, userResult)
			fmt.Println()

			if correctForm {
				correct := getUserInput(yes, "Is this correct?", yes+no, yes+" or "+no, "", 1)
				if correct == yes {
					guess = userGuess
					result = userResult
					break
				}
			}
		}

		guesses[try] = guess
		results[try] = result

		foundSolution := isResultCorrect(result, WordLength)
		printWordleSolution(guesses, results, try, foundSolution)
		if foundSolution {
			if DoWordle {
				if !usedWords[guess] && !IgnoreWordleUsedWords {
					addUsedWord := getUserInput(yes, "Would you like to add '"+guess+"' to the list of already used words?", yes+no, yes+" or "+no, "", 1)
					if addUsedWord == yes {
						f, err := os.OpenFile(UsedWordsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							log.Println(err)
						} else {
							defer f.Close()
							newUsedWord := "\n" + guess
							_, err = f.WriteString(newUsedWord)
							if err != nil {
								log.Println(err)
							}
						}
					}
				}
			}
			break
		}
	}
}

func main() {
	solutionWords, allWords, usedWords, guess, result := initialize()

	switch Mode {
	case ModeWordSearch:
		WordSearch(solutionWords)
	case ModeManualGuess:
		ManualGuess(guess, result, solutionWords, allWords)
	default:
		AutoPlay(guess, result, solutionWords, allWords, usedWords)
	}
}

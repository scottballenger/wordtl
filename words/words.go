package words

import (
	"fmt"
	"sort"
	"strings"
)

const (
	WildcardChar = "-"
	MatchedChar  = "="
	MissedChar   = "x"
)

func WordMatch(
	word string,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	matchAllWildcardLetters bool,
	excludedByPosMap map[int]string) bool {

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
	if len(wildcardLetters) > 0 {
		if matchAllWildcardLetters {
			// All wildcard letters must be included.
			for _, letter := range wildcardLetters {
				if !strings.Contains(word, string(letter)) {
					return false
				}
			}
		} else {
			// Any wildcard letters can be matched.
			match := false
			for _, letter := range wildcardLetters {
				if strings.Contains(word, string(letter)) {
					match = true
				}
			}
			if !match {
				return false
			}
		}
	}

	// filter for positional pattern match.
	for i, letter := range word {
		switch string(wordPattern[i : i+1]) {
		case string(letter):
			continue
		case WildcardChar:
			if len(excludedByPosMap[i+1]) > 0 {
				// Remove letters that can't be in a specific position.
				if strings.Contains(excludedByPosMap[i+1], string(letter)) {
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

func GetMatchingWords(
	words []string,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	matchAllWildcardLetters bool,
	excludedByPosMap map[int]string) []string {

	var matchingWords []string

	for _, word := range words {
		if WordMatch(word, wordPattern, excludedLetters, wildcardLetters, matchAllWildcardLetters, excludedByPosMap) {
			matchingWords = append(matchingWords, word)
		}
	}

	return matchingWords
}

func GetLetterCount(words []string, wordPattern string, wildcardLetters string) (map[string]int, string) {
	var letterCount = map[string]int{}
	letterOrdering := ""

	for _, word := range words {
		for _, letter := range word {
			if !strings.Contains(wildcardLetters, string(letter)) &&
				!strings.Contains(wordPattern, string(letter)) {
				letterCount[string(letter)]++
			}
		}
	}

	// Sort letters in order of most occurances first.
	keys := make([]string, 0, len(letterCount))
	for k := range letterCount {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return letterCount[keys[i]] > letterCount[keys[j]]
	})

	// Create string with most important letters first.
	for _, k := range keys {
		letterOrdering += k
	}

	return letterCount, letterOrdering
}

func GetLetterDistribution(words []string, wordLength int) []map[string]int {
	letterDistribution := []map[string]int{}

	if len(words) == 0 {
		return letterDistribution
	}

	for position := 0; position < wordLength; position++ {
		letterDistribution = append(letterDistribution, map[string]int{})
		for _, word := range words {
			letterDistribution[position][string(word[position])]++
		}
	}

	return letterDistribution
}

func GetEliminationWords(
	eliminationLetters string,
	words []string,
	wordLength int,
	excludedLetters string,
	wildcardLetters string,
	excludedByPosMap map[int]string) []string {

	fmt.Printf("\nTrying elimination letters: '%s'\n", eliminationLetters)

	return GetMatchingWords(words, strings.Repeat(WildcardChar, wordLength), "", eliminationLetters, false, map[int]string{})
}

func GetBestEliminationWords(words []string, eliminationWords []string, wordLength int, eliminationLetters string, letterCounts map[string]int, letterDistribution []map[string]int, debug bool) []string {
	bestEliminationWords := eliminationWords
	if len(eliminationWords) == 0 {
		return bestEliminationWords
	}

	lastLetters := ""
	remainingWords := []string{}
	included := make(map[string]bool)

	for index, letter := range eliminationLetters {
		if debug {
			fmt.Print("letter:", string(letter), "\n")
		}

		// Skip letter if already searched.
		if strings.Contains(lastLetters, string(letter)) {
			if debug {
				fmt.Print("\tskipping:", string(letter), "\n")
			}
			continue
		}

		if index >= wordLength {
			if len(bestEliminationWords) > 0 {
				// Found enough letters.
				if debug {
					fmt.Println("Searched ", index, "letters. Stopping.")
				}
				break
			}
		}

		// Check for other letters with the same letter count.
		count := letterCounts[string(letter)]
		letters := string(letter)
		for letterKey, letterCount := range letterCounts {
			if letterKey != string(letter) && letterCount == count {
				letters = letters + letterKey
			}
		}
		lastLetters = letters
		if debug {
			if len(letters) > 1 {
				fmt.Print("\tall letters:", letters, "\n")
			}
		}

		// Search all the letters
		for _, word := range bestEliminationWords {
			for _, thisLetter := range letters {
				if strings.Contains(word, string(thisLetter)) {
					if !included[word] {
						remainingWords = append(remainingWords, word)
						included[word] = true
						if debug {
							fmt.Print(word, " ")
						}
					}
					break
				}
			}
		}
		if debug {
			fmt.Println("\nlen(remainingWords): ", len(remainingWords))
		}
		if len(remainingWords) > 0 {
			bestEliminationWords = remainingWords
			if debug {
				fmt.Println("len(bestEliminationWords): ", len(bestEliminationWords))
			}
		}
		if len(bestEliminationWords) == 1 {
			break
		}
	}
	if debug {
		fmt.Println("1: len(bestEliminationWords): ", len(bestEliminationWords))
	}

	// Get the words with the most letters in order of highest number first
	if len(bestEliminationWords) > 1 {
		eliminationLettersCount := map[string]int{}
		for _, word := range bestEliminationWords {
			letterCount := 0
			for _, letter := range eliminationLetters {
				if strings.Contains(word, string(letter)) {
					letterCount++
				}
			}
			eliminationLettersCount[word] = letterCount
		}
		if debug {
			fmt.Println("eliminationLettersCount: ", eliminationLettersCount)
		}

		// Sort words in order of most elimination letters first.
		keys := make([]string, 0, len(eliminationLettersCount))
		for k := range eliminationLettersCount {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return eliminationLettersCount[keys[i]] > eliminationLettersCount[keys[j]]
		})

		// Save the words with the most elimination letters.
		bestEliminationWords = []string{}
		numEliminationLetters := eliminationLettersCount[keys[0]]
		for _, k := range keys {
			if eliminationLettersCount[k] == numEliminationLetters {
				bestEliminationWords = append(bestEliminationWords, k)
			} else {
				break
			}
		}

		// Score the remaining words with total elimination letter occurances.
		eliminationLettersScore := map[string]int{}
		for _, word := range bestEliminationWords {
			letterScore := 0
			for _, letter := range eliminationLetters {
				if strings.Contains(word, string(letter)) {
					letterScore += letterCounts[string(letter)]
				}
			}
			eliminationLettersScore[word] = letterScore
		}
		if debug {
			fmt.Println("eliminationLettersScore: ", eliminationLettersScore)
		}

		// Sort words in order of highest elimination score first.
		keys = make([]string, 0, len(eliminationLettersScore))
		for k := range eliminationLettersScore {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return eliminationLettersScore[keys[i]] > eliminationLettersScore[keys[j]]
		})

		// Save the words with the highest elimination score.
		bestEliminationWords = []string{}
		numEliminationOccurances := eliminationLettersScore[keys[0]]
		for _, k := range keys {
			if eliminationLettersScore[k] == numEliminationOccurances {
				bestEliminationWords = append(bestEliminationWords, k)
			} else {
				break
			}
		}
	}
	if debug {
		fmt.Println("2: len(bestEliminationWords): ", len(bestEliminationWords))
		fmt.Println("len(words): ", len(words))
	}

	// Prefer any solution words that are in remaining list.
	if len(bestEliminationWords) > len(words) {
		matchingWords := []string{}
		for _, word := range words {
			for _, eliminationWord := range bestEliminationWords {
				if eliminationWord == word {
					matchingWords = append(matchingWords, word)
				}
			}
		}
		if debug {
			fmt.Println("matchingWords: ", matchingWords)
		}
		if len(matchingWords) > 0 {
			bestEliminationWords = matchingWords
		}
	}
	if debug {
		fmt.Println("3: len(bestEliminationWords): ", len(bestEliminationWords))
	}

	// Sort the remaining words with respect to individual elimination letter distribution.
	if len(bestEliminationWords) > 1 {
		eliminationWordScore := map[string]int{}
		for _, word := range bestEliminationWords {
			for position, letter := range word {
				eliminationWordScore[word] += letterDistribution[position][string(letter)]
			}
		}
		if debug {
			fmt.Println("eliminationWordScore: ", eliminationWordScore)
		}

		// Sort words in order of highest elimination score first.
		keys := make([]string, 0, len(eliminationWordScore))
		for k := range eliminationWordScore {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return eliminationWordScore[keys[i]] > eliminationWordScore[keys[j]]
		})

		// Save the words with the highest elimination score first.
		bestEliminationWords = keys
	}
	if debug {
		fmt.Println("4: len(bestEliminationWords): ", len(bestEliminationWords))
		fmt.Println("bestEliminationWords: ", bestEliminationWords)
	}
	return bestEliminationWords
}

func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

func GuessWord(word string, guess string) (bool, string) {
	match := false
	word = strings.ToLower(word)
	guess = strings.ToLower(guess)
	result := ""
	if len(guess) == len(word) {
		match = true
		for i := range guess {
			guessLetter := string(guess[i])
			wordleLetter := string(word[i])
			if guessLetter == wordleLetter {
				result += MatchedChar
			} else {
				if strings.Contains(word, guessLetter) {
					result += WildcardChar
				} else {
					result += MissedChar
				}
				match = false
			}
		}
	}
	return match, result
}

func TranslateGuessResults(
	guess string,
	results string,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	excludedByPosMap map[int]string) (string, string, string, map[int]string) {

	if len(guess) == len(results) {
		for i := range results {
			guessLetter := string(guess[i])
			switch string(results[i]) {
			case MatchedChar:
				wordPattern = replaceAtIndex(wordPattern, rune(guess[i]), i)
			case WildcardChar:
				if !strings.Contains(wildcardLetters, guessLetter) {
					wildcardLetters += guessLetter
				}
				if !strings.Contains(excludedByPosMap[i+1], guessLetter) {
					excludedByPosMap[i+1] += guessLetter
				}
			case MissedChar:
				// Can be another instance of an exsiting letter.
				matchLater := false
				for j := i + 1; j < len(results); j++ {
					if string(guess[j]) == guessLetter && (string(results[j]) == MatchedChar || string(results[j]) == WildcardChar) {
						matchLater = true
						break
					}
				}
				if !matchLater && !strings.Contains(wordPattern, guessLetter) && !strings.Contains(wildcardLetters, guessLetter) {
					if !strings.Contains(excludedLetters, guessLetter) {
						excludedLetters += guessLetter
					}
				} else {
					if !strings.Contains(excludedByPosMap[i+1], guessLetter) {
						excludedByPosMap[i+1] += guessLetter
					}
				}
			}
		}
	}
	return wordPattern, wildcardLetters, excludedLetters, excludedByPosMap
}

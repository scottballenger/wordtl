package words

import (
	"fmt"
	"sort"
	"strings"
)

const (
	WildcardChar = "-"
	MatchedChar  = "="
	MissedChar   = "X"
	MaxLetters   = 9
)

func WordMatch(
	word string,
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
	if len(wildcardLetters) > 0 {
		if len(wildcardLetters) <= len(word) {
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
			if i < MaxLetters {
				if len(noParkDisSpace[i]) > 0 {
					// Remove letters that can't be in a specific position.
					if strings.Contains(noParkDisSpace[i], string(letter)) {
						return false
					}
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
	noParkDisSpace [MaxLetters]string) []string {

	var matchingWords []string

	for _, word := range words {
		if WordMatch(word, wordPattern, excludedLetters, wildcardLetters, noParkDisSpace) {
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
	noParkDisSpace [MaxLetters]string) []string {

	var eliminationWords []string

	fmt.Printf("\nTrying elimination letters: '%s'\n", eliminationLetters)
	numLetters := wordLength
	if len(eliminationLetters) < wordLength {
		numLetters = len(eliminationLetters)
	}
	for ; numLetters > 0; numLetters-- {
		eliminationWords = GetMatchingWords(words, strings.Repeat(WildcardChar, wordLength), "", eliminationLetters[:numLetters], [MaxLetters]string{})
		if len(eliminationWords) > 0 {
			break
		}
	}

	return eliminationWords
}

func GetBestEliminationWords(words []string, eliminationWords []string, wordLength int, eliminationLetters string, letterCounts map[string]int) []string {
	bestEliminationWords := eliminationWords
	if len(eliminationWords) == 0 {
		return bestEliminationWords
	}

	lastLetters := ""
	for _, letter := range eliminationLetters {
		remainingWords := []string{}
		// fmt.Print("letter:", string(letter), "\n")

		// Skip letter if already searched.
		if strings.Contains(lastLetters, string(letter)) {
			// fmt.Print("\tskipping:", string(letter), "\n")
			continue
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
		// if len(letters) > 1 {
		// 	fmt.Print("\tall letters:", letters, "\n")
		// }

		// Search all the letters
		for _, word := range bestEliminationWords {
			for _, thisLetter := range letters {
				if strings.Contains(word, string(thisLetter)) {
					remainingWords = append(remainingWords, word)
					break
				}
			}
		}
		// fmt.Printf("remainingWords:%v\n", remainingWords)
		if len(remainingWords) > 0 {
			bestEliminationWords = remainingWords
			// fmt.Printf("bestEliminationWords:%v\n", bestEliminationWords)
		} else {
			break
		}
		if len(bestEliminationWords) == 1 {
			break
		}
	}

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
		// fmt.Println("eliminationLettersCount: ", eliminationLettersCount)

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
		// fmt.Println("eliminationLettersScore: ", eliminationLettersScore)

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

		matchingWords := []string{}
		// Prefer any solution words that are in remining list
		for _, word := range words {
			for _, eliminationWord := range bestEliminationWords {
				if eliminationWord == word {
					matchingWords = append(matchingWords, word)
				}
			}
		}
		// fmt.Println("matchingWords: ", matchingWords)
		if len(matchingWords) > 0 {
			bestEliminationWords = matchingWords
		}

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
	noParkDisSpace [MaxLetters]string) (string, string, string, [MaxLetters]string) {

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
				if i < MaxLetters {
					if !strings.Contains(noParkDisSpace[i], guessLetter) {
						noParkDisSpace[i] += guessLetter
					}
				}
			case MissedChar:
				if !strings.Contains(excludedLetters, guessLetter) {
					excludedLetters += guessLetter
				}
			}
		}
	}
	return wordPattern, wildcardLetters, excludedLetters, noParkDisSpace
}

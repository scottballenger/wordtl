package words

import (
	"bufio"
	"io"
	"log"
	"strings"
)

const (
	WildcardChar = "-"
	MaxLetters   = 9
)

func WordMatch(word string,
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

func GetMatchingWords(
	wordFileHandle io.Reader,
	wordPattern string,
	excludedLetters string,
	wildcardLetters string,
	noParkDisSpace [MaxLetters]string) []string {

	var matchingWords []string

	scanner := bufio.NewScanner(wordFileHandle)
	for scanner.Scan() {
		word := strings.ToLower(scanner.Text())
		if WordMatch(word, wordPattern, excludedLetters, wildcardLetters, noParkDisSpace) {
			matchingWords = append(matchingWords, word)

		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return matchingWords
}

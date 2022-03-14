package words

import (
	"strings"
	"testing"
)

func Test_wordMatch(t *testing.T) {
	type args struct {
		word            string
		wordPattern     string
		excludedLetters string
		wildcardLetters string
		noParkDisSpace  [MaxLetters]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Empty Word",
			args: args{word: "", wordPattern: "t----"},
			want: false,
		},
		{
			name: "Excluded Letters",
			args: args{word: "abcde", wordPattern: "t----", excludedLetters: "a"},
			want: false,
		},
		{
			name: "Excluded Letters",
			args: args{word: "abcde", wordPattern: "t----", wildcardLetters: "f"},
			want: false,
		},
		{
			name: "Pattern Match all Letters",
			args: args{word: "abcde", wordPattern: "abcde"},
			want: true,
		},
		{
			name: "Pattern Match no Letters",
			args: args{word: "abcde", wordPattern: "fghij"},
			want: false,
		},
		{
			name: "Wildcard Match all Letters",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcde"},
			want: true,
		},
		{
			name: "Wildcard Match, but can't be in current position",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcde",
				noParkDisSpace: [MaxLetters]string{"", "b"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordMatch(tt.args.word, tt.args.wordPattern, tt.args.excludedLetters, tt.args.wildcardLetters, tt.args.noParkDisSpace); got != tt.want {
				t.Errorf("wordMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMatchingWords(t *testing.T) {
	type args struct {
		words           []string
		wordPattern     string
		excludedLetters string
		wildcardLetters string
		noParkDisSpace  [MaxLetters]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Empty File",
			args: args{words: []string{}},
			want: []string{},
		},
		{
			name: "Blank Word",
			args: args{words: []string{"     "}, wordPattern: "t----"},
			want: []string{},
		},
		{
			name: "Pattern Match all Letters",
			args: args{words: []string{"abcde"}, wordPattern: "abcde"},
			want: []string{"abcde"},
		},
		{
			name: "Pattern Match multiple Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "targa"}, wordPattern: "t----"},
			want: []string{"tabor", "talar", "tardo", "tardy", "targa"},
		},
		{
			name: "Pattern Match not all Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "barga"}, wordPattern: "t----"},
			want: []string{"tabor", "talar", "tardo", "tardy"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMatchingWords(tt.args.words, tt.args.wordPattern, tt.args.excludedLetters, tt.args.wildcardLetters, tt.args.noParkDisSpace); strings.TrimSpace(strings.Join(got, "")) != strings.TrimSpace(strings.Join(tt.want, "")) {
				t.Errorf("getMatchingWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
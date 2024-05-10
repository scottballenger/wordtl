package words

import (
	"reflect"
	"strings"
	"testing"
)

func Test_wordMatch(t *testing.T) {
	type args struct {
		word                    string
		wordPattern             string
		excludedLetters         string
		wildcardLetters         string
		matchAllWildcardLetters bool
		excludedByPosMap        map[int]string
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
			name: "Wildcard Match all Letters, must match all",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcde", matchAllWildcardLetters: true},
			want: true,
		},
		{
			name: "Wildcard does not Match all Letters, must match all",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcdef", matchAllWildcardLetters: true},
			want: false,
		},
		{
			name: "Wildcard Match any Letters",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcdef", matchAllWildcardLetters: false},
			want: true,
		},
		{
			name: "Wildcard Match, but can't be in current position",
			args: args{word: "abcde", wordPattern: "-----", wildcardLetters: "abcde", matchAllWildcardLetters: true,
				excludedByPosMap: map[int]string{2: "b"}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WordMatch(tt.args.word, tt.args.wordPattern, tt.args.excludedLetters, tt.args.wildcardLetters, tt.args.matchAllWildcardLetters, tt.args.excludedByPosMap); got != tt.want {
				t.Errorf("wordMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getMatchingWords(t *testing.T) {
	type args struct {
		words                   []string
		wordPattern             string
		excludedLetters         string
		wildcardLetters         string
		matchAllWildcardLetters bool
		excludedByPosMap        map[int]string
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
			name: "Pattern Match all Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "targa"}, wordPattern: "t----"},
			want: []string{"tabor", "talar", "tardo", "tardy", "targa"},
		},
		{
			name: "Pattern Match not all Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "barga"}, wordPattern: "t----"},
			want: []string{"tabor", "talar", "tardo", "tardy"},
		},
		{
			name: "Wildcard Match all Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "barga"}, wordPattern: "-----", wildcardLetters: "tar", matchAllWildcardLetters: false},
			want: []string{"tabor", "talar", "tardo", "tardy", "barga"},
		},
		{
			name: "Wildcard Match not all Words",
			args: args{words: []string{"tabor", "talar", "tardo", "tardy", "barga"}, wordPattern: "-----", wildcardLetters: "tar", matchAllWildcardLetters: true},
			want: []string{"tabor", "talar", "tardo", "tardy"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMatchingWords(tt.args.words, tt.args.wordPattern, tt.args.excludedLetters, tt.args.wildcardLetters, tt.args.matchAllWildcardLetters, tt.args.excludedByPosMap); strings.TrimSpace(strings.Join(got, "")) != strings.TrimSpace(strings.Join(tt.want, "")) {
				t.Errorf("getMatchingWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranslateGuessResults(t *testing.T) {
	type args struct {
		guess            string
		results          string
		wordPattern      string
		excludedLetters  string
		wildcardLetters  string
		excludedByPosMap map[int]string
	}
	tests := []struct {
		name                 string
		args                 args
		wantWordPattern      string
		wantWildcardLetters  string
		wantExcludedLetters  string
		wantExcludedByPosMap map[int]string
	}{
		{
			name:                 "Match All",
			args:                 args{guess: "abcde", results: "=====", wordPattern: "-----", excludedLetters: "", wildcardLetters: "", excludedByPosMap: map[int]string{}},
			wantWordPattern:      "abcde",
			wantWildcardLetters:  "",
			wantExcludedLetters:  "",
			wantExcludedByPosMap: map[int]string{},
		},
		{
			name:                 "Match Some, Others In Wrong Position",
			args:                 args{guess: "abcde", results: "-===-", wordPattern: "-----", excludedLetters: "", wildcardLetters: "", excludedByPosMap: map[int]string{}},
			wantWordPattern:      "-bcd-",
			wantWildcardLetters:  "ae",
			wantExcludedLetters:  "",
			wantExcludedByPosMap: map[int]string{1: "a", 5: "e"},
		},
		{
			name:                 "Repeating Letter Matching Earlier In Word",
			args:                 args{guess: "chick", results: "===xx", wordPattern: "--i--", excludedLetters: "", wildcardLetters: "", excludedByPosMap: map[int]string{}},
			wantWordPattern:      "chi--",
			wantWildcardLetters:  "",
			wantExcludedLetters:  "k",
			wantExcludedByPosMap: map[int]string{4: "c"},
		},
		{
			name:                 "Repeating Letter Matching Later In Word",
			args:                 args{guess: "chick", results: "x===x", wordPattern: "--i--", excludedLetters: "", wildcardLetters: "", excludedByPosMap: map[int]string{}},
			wantWordPattern:      "-hic-",
			wantWildcardLetters:  "",
			wantExcludedLetters:  "k",
			wantExcludedByPosMap: map[int]string{1: "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWordPattern, gotWildcardLetters, gotExcludedLetters, gotExcludedByPosMap := TranslateGuessResults(tt.args.guess, tt.args.results, tt.args.wordPattern, tt.args.excludedLetters, tt.args.wildcardLetters, tt.args.excludedByPosMap)
			if gotWordPattern != tt.wantWordPattern {
				t.Errorf("TranslateGuessResults() gotWordPattern = %v, want %v", gotWordPattern, tt.wantWordPattern)
			}
			if gotWildcardLetters != tt.wantWildcardLetters {
				t.Errorf("TranslateGuessResults() gotWildcardLetters = %v, want %v", gotWildcardLetters, tt.wantWildcardLetters)
			}
			if gotExcludedLetters != tt.wantExcludedLetters {
				t.Errorf("TranslateGuessResults() gotExcludedLetters = %v, want %v", gotExcludedLetters, tt.wantExcludedLetters)
			}
			if !reflect.DeepEqual(gotExcludedByPosMap, tt.wantExcludedByPosMap) {
				t.Errorf("TranslateGuessResults() gotExcludedByPosMap = %v, want %v", gotExcludedByPosMap, tt.wantExcludedByPosMap)
			}
		})
	}
}

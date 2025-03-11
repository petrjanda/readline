package readline

import (
	"sort"
	"strings"
	"unicode"
)

// FuzzyCompleter implements AutoCompleter with fuzzy matching capabilities
type FuzzyCompleter struct {
	items []string
}

// NewFuzzyCompleter creates a new FuzzyCompleter with the provided items
func NewFuzzyCompleter(items []string) *FuzzyCompleter {
	return &FuzzyCompleter{
		items: items,
	}
}

// Do implements the AutoCompleter interface
func (f *FuzzyCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	// First, find the start of the current word
	wordStart := pos
	for wordStart > 0 && !unicode.IsSpace(line[wordStart-1]) {
		wordStart--
	}

	// Extract the current word
	currentWord := string(line[wordStart:pos])

	// Find matching items using fuzzy matching
	var candidates [][]rune

	// If the current word is empty, return all items
	if len(currentWord) == 0 {
		for _, item := range f.items {
			candidates = append(candidates, []rune(item))
		}
		// Even though there's nothing to replace, we need to return a length
		// that will allow the completion mechanism to insert text
		return candidates, 0
	}

	// For non-empty input, perform fuzzy matching
	for _, item := range f.items {
		if fuzzyMatch(currentWord, item) {
			candidates = append(candidates, []rune(item))
		}
	}

	// Sort candidates by relevance if needed
	// (could implement a more sophisticated sorting algorithm here)
	sort.Slice(candidates, func(i, j int) bool {
		return string(candidates[i]) < string(candidates[j])
	})

	return candidates, pos - wordStart
}

// fuzzyMatch checks if the pattern appears in the text in order
func fuzzyMatch(pattern, text string) bool {
	patternLower := strings.ToLower(pattern)
	textLower := strings.ToLower(text)

	// If the pattern is a substring, it's definitely a match
	if strings.Contains(textLower, patternLower) {
		return true
	}

	// Otherwise, check if the characters appear in order
	patIndex := 0
	for _, char := range textLower {
		if patIndex < len(patternLower) && char == rune(patternLower[patIndex]) {
			patIndex++
		}
	}

	return patIndex == len(patternLower)
}

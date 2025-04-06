package readline

import (
	"sort"
	"unicode"
)

// FuzzyCompleter implements AutoCompleter with fuzzy matching capabilities
type FuzzyCompleter struct {
	lazyMatch func(term string) [][]rune

	maxItems int
}

type FuzzyCompleterOption func(*FuzzyCompleter)

func WithMaxItems(maxItems int) FuzzyCompleterOption {
	return func(f *FuzzyCompleter) {
		f.maxItems = maxItems
	}
}

// NewFuzzyCompleter creates a new FuzzyCompleter with the provided items
func NewFuzzyCompleter(lazyMatch func(term string) [][]rune, opts ...FuzzyCompleterOption) *FuzzyCompleter {
	f := &FuzzyCompleter{
		lazyMatch: lazyMatch,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
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

	candidates := f.lazyMatch(currentWord)

	if f.maxItems > 0 && len(candidates) > f.maxItems {
		candidates = candidates[:f.maxItems]
	}

	// Sort candidates by relevance if needed
	// (could implement a more sophisticated sorting algorithm here)
	sort.Slice(candidates, func(i, j int) bool {
		return string(candidates[i]) < string(candidates[j])
	})

	return candidates, pos - wordStart
}

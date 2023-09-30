package prunningRadixTrie

import (
	"math/rand"
	"fmt"
	"testing"
	"strconv"
)

func BenchmarkAddTopKSuggestion(b *testing.B) {
	trie := NewPruningRadixTrie()
	results := make([]struct {
		term              string
		termFrequencyCount int64
	}, 0)

	// Initialize the results slice with some data
	for i := 0; i < 100; i++ {
		results = append(results, struct {
			term              string
			termFrequencyCount int64
		}{term: "term" + strconv.Itoa(i), termFrequencyCount: int64(rand.Intn(100))})
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate adding a new suggestion
		term := "newTerm"
		termFrequencyCount := int64(rand.Intn(100))
		trie.addTopKSuggestion(term, termFrequencyCount, 10, &results)
	}
}

func ExampleAddTopKSuggestion() {
	trie := NewPruningRadixTrie()
	trie.AddTerm("apple", 5)
	trie.AddTerm("appetizer", 3)
	trie.AddTerm("appetite", 2)
	trie.AddTerm("banana", 4)

	results, termFrequencyCountPrefix := trie.GetTopkTermsForPrefix("app", 3, true)
	fmt.Println("Top K Terms:")
	for _, result := range results {
		fmt.Printf("%s: %d\n", result.term, result.termFrequencyCount)
	}
	fmt.Println("Term Frequency Count Prefix:", termFrequencyCountPrefix)

	trie.WriteTermsToFile("terms.txt")
}

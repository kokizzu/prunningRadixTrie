package prunningRadixTrie

import (
	"fmt"
	"math"
	"sort"
	"io"
	"bufio"
	"strconv"
	"strings"
	"os"
)

type Node struct {
	Children                []struct {
		key  string
		node *Node
	}
	termFrequencyCount      int64
	termFrequencyCountChildMax int64
}

func NewNode(termFrequencyCount int64) *Node {
	return &Node{
		termFrequencyCount: termFrequencyCount,
	}
}

type PruningRadixTrie struct {
	termCount      int64
	termCountLoaded int64
	trie           *Node
}

func NewPruningRadixTrie() *PruningRadixTrie {
	return &PruningRadixTrie{
		trie: &Node{},
	}
}

func (t *PruningRadixTrie) AddTerm(term string, termFrequencyCount int64) {
	nodeList := make([]*Node, 0)
	t.addTerm(t.trie, term, termFrequencyCount, 0, 0, &nodeList)
}

func (t *PruningRadixTrie) UpdateMaxCounts(nodeList []*Node, termFrequencyCount int64) {
 	for _, node := range nodeList {
        if termFrequencyCount > node.termFrequencyCountChildMax {
            node.termFrequencyCountChildMax = termFrequencyCount
        }
    }
}

func (t *PruningRadixTrie) addTerm(curr *Node, term string, termFrequencyCount int64, id int, level int, nodeList *[]*Node) {
	*nodeList = append(*nodeList, curr)

	common := 0
	if curr.Children != nil {
		for j := 0; j < len(curr.Children); j++ {
			child := curr.Children[j].node
			key := curr.Children[j].key

			for i := 0; i < int(math.Min(float64(len(term)), float64(len(key)))); i++ {
				if term[i] == key[i] {
					common = i + 1
				} else {
					break
				}
			}

			if common > 0 {
				if common == len(term) && common == len(key) {
					if child.termFrequencyCount == 0 {
						t.termCount++
					}
					child.termFrequencyCount += termFrequencyCount
					t.UpdateMaxCounts(*nodeList, child.termFrequencyCount)
				} else if common == len(term) {
					newChild := &Node{
						termFrequencyCount: termFrequencyCount,
						Children: []struct {
							key  string
							node *Node
						}{
							{key[len(term):], child},
						},
						termFrequencyCountChildMax: int64(math.Max(float64(child.termFrequencyCountChildMax), float64(child.termFrequencyCount))),
					}
					t.UpdateMaxCounts(*nodeList, termFrequencyCount)
					curr.Children[j] = struct {
						key  string
						node *Node
					}{term[:common], newChild}
					sort.Slice(curr.Children, func(i, j int) bool {
						return curr.Children[i].node.termFrequencyCountChildMax > curr.Children[j].node.termFrequencyCountChildMax
					})
					t.termCount++
				} else if common == len(key) {
					t.addTerm(child, term[common:], termFrequencyCount, id, level+1, nodeList)
				} else {
					newChild := &Node{
						termFrequencyCount: termFrequencyCount,
						Children: []struct {
							key  string
							node *Node
						}{
							{key[common:], child},
							{term[common:], &Node{termFrequencyCount: termFrequencyCount}},
						},
						termFrequencyCountChildMax: int64(math.Max(float64(child.termFrequencyCountChildMax), math.Max(float64(termFrequencyCount), float64(child.termFrequencyCount)))),
					}
					t.UpdateMaxCounts(*nodeList, termFrequencyCount)
					curr.Children[j] = struct {
						key  string
						node *Node
					}{term[:common], newChild}
					sort.Slice(curr.Children, func(i, j int) bool {
						return curr.Children[i].node.termFrequencyCountChildMax > curr.Children[j].node.termFrequencyCountChildMax
					})
					t.termCount++
				}
				return
			}
		}
	}

	if curr.Children == nil {
		curr.Children = []struct {
			key  string
			node *Node
		}{
			{term, &Node{termFrequencyCount: termFrequencyCount}},
		}
	} else {
		curr.Children = append(curr.Children, struct {
			key  string
			node *Node
		}{term, &Node{termFrequencyCount: termFrequencyCount}})
		sort.Slice(curr.Children, func(i, j int) bool {
			return curr.Children[i].node.termFrequencyCountChildMax > curr.Children[j].node.termFrequencyCountChildMax
		})
	}
	t.termCount++
	t.UpdateMaxCounts(*nodeList, termFrequencyCount)
}

func (t *PruningRadixTrie) FindAllChildTerms(prefix string, topK int, termFrequencyCountPrefix *int64, prefixString string, results *[]struct {
	term              string
	termFrequencyCount int64
}, pruning bool) {
	t.findAllChildTerms(prefix, t.trie, topK, termFrequencyCountPrefix, prefixString, results, nil, pruning)
}

func (t *PruningRadixTrie) findAllChildTerms(prefix string, curr *Node, topK int, termFrequencyCountPrefix *int64, prefixString string, results *[]struct {
	term              string
	termFrequencyCount int64
}, file io.Writer, pruning bool) {
	if pruning && topK > 0 && len(*results) == topK && curr.termFrequencyCountChildMax <= (*results)[topK-1].termFrequencyCount {
		return
	}

	noPrefix := prefix == ""

	if curr.Children != nil {
		for _, child := range curr.Children {
			key := child.key
			node := child.node

			if pruning && topK > 0 && len(*results) == topK && node.termFrequencyCount <= (*results)[topK-1].termFrequencyCount && node.termFrequencyCountChildMax <= (*results)[topK-1].termFrequencyCount {
				if !noPrefix {
					break
				} else {
					continue
				}
			}

			if noPrefix || (len(key) >= len(prefix) && key[:len(prefix)] == prefix) {
				if node.termFrequencyCount > 0 {
					if prefix == key {
						*termFrequencyCountPrefix = node.termFrequencyCount
					}

					if file != nil {
						fmt.Fprintf(file, "%s%s\t%d\n", prefixString, key, node.termFrequencyCount)
					} else if topK > 0 {
						t.addTopKSuggestion(prefixString+key, node.termFrequencyCount, topK, results)
					} else {
						*results = append(*results, struct {
							term              string
							termFrequencyCount int64
						}{prefixString + key, node.termFrequencyCount})
					}
				}

				if node.Children != nil && len(node.Children) > 0 {
					t.findAllChildTerms("", node, topK, termFrequencyCountPrefix, prefixString+key, results, file, pruning)
				}
				if !noPrefix {
					break
				}
			} else if len(prefix) >= len(key) && prefix[:len(key)] == key {
				if node.Children != nil && len(node.Children) > 0 {
					t.findAllChildTerms(prefix[len(key):], node, topK, termFrequencyCountPrefix, prefixString+key, results, file, pruning)
				}
				break
			}
		}
	}
}

func (t *PruningRadixTrie) GetTopkTermsForPrefix(prefix string, topK int, pruning bool) ([]struct {
	term              string
	termFrequencyCount int64
}, int64) {
	results := make([]struct {
		term              string
		termFrequencyCount int64
	}, 0)

	termFrequencyCountPrefix := int64(0)

	t.FindAllChildTerms(prefix, topK, &termFrequencyCountPrefix, "", &results, pruning)

	return results, termFrequencyCountPrefix
}

func (t *PruningRadixTrie) WriteTermsToFile(path string) {
	if t.termCountLoaded == t.termCount {
		return
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	prefixCount := int64(0)
	t.findAllChildTerms("", t.trie, 0, &prefixCount, "", nil, file, true)
	fmt.Printf("%d terms written.\n", t.termCount)
}

func (t *PruningRadixTrie) ReadTermsFromFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("Could not find file", path)
		return false
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lineParts := strings.Split(line, "\t")

		if len(lineParts) == 2 {
			if count, err := strconv.ParseInt(lineParts[1], 10, 64); err == nil {
				t.AddTerm(lineParts[0], count)
			}
		}
	}

	t.termCountLoaded = t.termCount
	fmt.Printf("%d terms loaded.\n", t.termCount)
	return true
}

func (t *PruningRadixTrie) addTopKSuggestion(term string, termFrequencyCount int64, topK int, results *[]struct {
	term              string
	termFrequencyCount int64
}) {
	if len(*results) < topK || termFrequencyCount >= (*results)[topK-1].termFrequencyCount {
		index := sort.Search(len(*results), func(i int) bool {
			return termFrequencyCount > (*results)[i].termFrequencyCount
		})
		newResult := struct {
			term              string
			termFrequencyCount int64
		}{term, termFrequencyCount}
		*results = append(*results, struct {
			term              string
			termFrequencyCount int64
		}{})

		copy((*results)[index+1:], (*results)[index:])
		(*results)[index] = newResult

		if len(*results) > topK {
			*results = (*results)[:topK]
		}
	}
}

package bk_tree

import (
	"bytes"
	"encoding/gob"
	"os"
	"sort"

	"github.com/avalanche-pwn/cdrepo/internal/searchif"
)

const maxEditDistance = 100

type edge struct {
	value int
	elem  *node
}

func (e *edge) init(v int) {
	e.value = v
	e.elem = nil
}

type node struct {
	value searchif.SearchNode
	edges []*edge
}

func (n *node) init(s searchif.SearchNode) {
	n.value = s
	n.edges = make([]*edge, 0)
}

type BKTree struct {
	root *node
}

func levenshtein(a, b string, maxDistance int) int {
	// Normalize lengths so that len(a) >= len(b) to minimize workspace.
	na := len(a)
	nb := len(b)
	if na < nb {
		// swap so a is the longer
		a, b = b, a
		na, nb = nb, na
	}

	// If the shorter string is empty, distance is length of the longer.
	if nb == 0 {
		if maxDistance >= 0 && na > maxDistance {
			return maxDistance + 1
		}
		return na
	}

	// If bounded and length difference already exceeds bound, early exit.
	lenDiff := na - nb
	if maxDistance >= 0 && lenDiff > maxDistance {
		return maxDistance + 1
	}

	// Convert to rune slices if you want to handle unicode properly.
	// For byte-wise (UTF-8 bytes) comparison use []byte(a) / []byte(b).
	ar := []rune(a)
	br := []rune(b)
	na = len(ar)
	nb = len(br)
	// Ensure ar is the longer (swap if necessary after rune conversion)
	if na < nb {
		ar, br = br, ar
		na, nb = nb, na
	}

	// Two-row DP: previous and current
	prev := make([]int, nb+1)
	curr := make([]int, nb+1)

	// Initialize previous row: distance from empty string to prefix of br
	for j := 0; j <= nb; j++ {
		prev[j] = j
	}

	// Main loop
	for i := 1; i <= na; i++ {
		curr[0] = i

		// Track minimal value in current row for early cutoff
		rowMin := curr[0]

		ai := ar[i-1]
		for j := 1; j <= nb; j++ {
			cost := 0
			if ai != br[j-1] {
				cost = 1
			}
			// substitution, insertion, deletion
			s := prev[j-1] + cost
			ins := curr[j-1] + 1
			del := prev[j] + 1

			// min(s, ins, del)
			if ins < s {
				s = ins
			}
			if del < s {
				s = del
			}
			curr[j] = s
			if s < rowMin {
				rowMin = s
			}
		}

		// If bounded and the smallest value in this row exceeds maxDistance, early exit
		if maxDistance >= 0 && rowMin > maxDistance {
			return maxDistance + 1
		}

		// swap prev and curr
		prev, curr = curr, prev
	}

	result := prev[nb]
	if maxDistance >= 0 && result > maxDistance {
		return maxDistance + 1
	}
	return result
}

func (bktree *BKTree) Add(s searchif.SearchNode) searchif.SearchNode {
	if bktree.root == nil {
		var r node
		bktree.root = &r
		bktree.root.init(s)
		return s
	}
	current := bktree.root
	for {
		distance := levenshtein(current.value.Key(), s.Key(), -1)
		for _, e := range current.edges {
			if e.value != distance {
				continue
			}
			current = e.elem
			continue
		}
		if current.value.Key() == s.Key() {
			return current.value
		}
		var n node
		var e edge
		n.init(s)
		e.init(distance)
		e.elem = &n
		current.edges = append(current.edges, &e)
		break
	}
	return s
}

func (bktree *BKTree) Search(s string) []*searchif.SearchResult {
	if bktree.root == nil {
		return nil
	}
	candidates := make([]*node, 0)
	result := make([]*searchif.SearchResult, 0)
	candidates = append(candidates, bktree.root)
	for len(candidates) != 0 {
		var current *node
		current, candidates = candidates[0], candidates[1:]
		distance := levenshtein(current.value.Key(), s, -1)
		if distance < maxEditDistance {
			result = append(result,
				&searchif.SearchResult{Score: distance, Value: current.value})
		}
		low := distance - maxEditDistance
		high := distance + maxEditDistance
		for _, e := range current.edges {
			if e.value <= high && e.value >= low {
				candidates = append(candidates, e.elem)
			}
		}
	}
	sort.Sort(searchif.ByScore(result))
	return result
}

func (bktree *BKTree) Save(s string) {
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.

	// Encode (send) some values.
	enc.Encode(bktree.Encode())
	os.WriteFile(s, network.Bytes(), 0644)
}

func (bktree *BKTree) Read(s string) {
	f, _ := os.Open(s)
	defer f.Close()
	dec := gob.NewDecoder(f) // Will write to network.
	var test modelOfBKTree

	// Encode (send) some values.
	dec.Decode(&test)
	bktree.Decode(test)
}

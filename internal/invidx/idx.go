package invidx

import (
	"strings"
	"sort"

	"github.com/avalanche-pwn/cdrepo/internal/bk_tree"
	"github.com/avalanche-pwn/cdrepo/internal/searchif"
)

const searchSpaceLimit int = 3

type internalSearchNode struct {
	word    string
	entries []int
}

func (i *internalSearchNode) Key() string {
	return i.word
}

type InvIdx struct {
	mappingSearcher searchif.FuzzySearcher
	entries         []searchif.SearchNode
}

func searchFactory() searchif.FuzzySearcher {
	return &bk_tree.BKTree{}
}

func (i *InvIdx) Init(s searchif.FuzzySearcher) {
	if s != nil {
		i.mappingSearcher = s
		return
	}
	i.mappingSearcher = searchFactory()
}

func (i *InvIdx) Add(s searchif.SearchNode) searchif.SearchNode {
	for _, entry := range i.entries {
		if entry.Key() == s.Key() {
			return entry
		}
	}
	idx := len(i.entries)
	i.entries = append(i.entries, s)

	for word := range strings.SplitSeq(s.Key(), "/") {
		isn := internalSearchNode{word: word}
		res := i.mappingSearcher.Add(&isn).(*internalSearchNode)
		res.entries = append(res.entries, idx)
	}
	return s
}

func (i *InvIdx) Search(s string) []*searchif.SearchResult {
	m := make(map[int]int)
	for word := range strings.SplitSeq(s, "/") {
		res := i.mappingSearcher.Search(word)
		if len(res) == 0 {
			continue
		}
		for match_idx, match := range res[:searchSpaceLimit] {
			current := match.Value.(*internalSearchNode)
			for _, idx := range current.entries {
				count, _ := m[idx]
				m[idx] = count + searchSpaceLimit - match_idx
			}
		}
	}

	var ss []*searchif.SearchResult
	for k, v := range m {
		res := searchif.SearchResult{Score: v, Value: i.entries[k]}
		ss = append(ss, &res)
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Score > ss[j].Score
	})
	return ss
}

func (i *InvIdx) Read(s string) {}
func (i *InvIdx) Save(s string) {}

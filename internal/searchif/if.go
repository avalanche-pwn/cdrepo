package searchif

type SearchResult struct {
	Score int
	Value SearchNode
}

type ViewSearchResult struct {
	Score int
	Value string
}

type ByScore []*SearchResult

func (r ByScore) Len() int           { return len(r) }
func (r ByScore) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool { return r[i].Score < r[j].Score }

type SearchNode interface {
	Key() string
}

type FuzzySearcher interface {
	Add(s SearchNode) SearchNode
	Read(s string)
	Save(s string)
	Search(s string) []*SearchResult
}

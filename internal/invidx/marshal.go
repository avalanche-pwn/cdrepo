package invidx

import (
	"encoding/gob"

	"github.com/avalanche-pwn/cdrepo/internal/searchif"
)

type modelOfInvIdx struct {
	MappingSearcher searchif.Decodable
	Entries         []searchif.SearchNode
}

func init() {
	gob.Register(modelOfInvIdx{})
	gob.Register(&internalSearchNode{})
}

func (m modelOfInvIdx) Decode() searchif.FuzzySearcher {
	var i InvIdx
	i.Init(m.MappingSearcher.Decode())
	i.entries = m.Entries
	return &i
}

func (i *InvIdx) Encode() searchif.Decodable {
	return modelOfInvIdx{MappingSearcher: i.mappingSearcher.Encode(),
		Entries: i.entries}
}

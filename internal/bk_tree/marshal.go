package bk_tree

import (
	"encoding/gob"

	"github.com/avalanche-pwn/cdrepo/internal/searchif"
)

type modelOfEdge struct {
	Value int
	Elem  *modelOfNode
}

type modelOfNode struct {
	Value searchif.SearchNode
	Edges []*modelOfEdge
}

type modelOfBKTree struct {
	Root *modelOfNode
}

func (e *edge) Encode() modelOfEdge {
	model := e.elem.Encode()
	return modelOfEdge{Value: e.value, Elem: &model}
}

func (m modelOfEdge) Decode() edge {
	n := m.Elem.Decode()
	return edge{value: m.Value, elem: &n}
}

func (n *node) Encode() modelOfNode {
	model := modelOfNode{Value: n.value}
	for _, e := range n.edges {
		edge_m := e.Encode()
		model.Edges = append(model.Edges, &edge_m)
	}
	return model
}

func (m *modelOfNode) Decode() node {
	n := node{value: m.Value}
	for _, model_edge := range m.Edges {
		e := model_edge.Decode()
		n.edges = append(n.edges, &e)
	}
	return n
}

func (m *modelOfBKTree) Decode() searchif.FuzzySearcher {
	root := m.Root.Decode()
	return &BKTree{root: &root}
}

func (bktree *BKTree) Encode() searchif.Decodable {
	var root_model modelOfNode
	if bktree.root != nil {
		root_model = bktree.root.Encode()
	}
	model := modelOfBKTree{Root: &root_model}
	return &model
}

func init() {
	gob.Register(&modelOfBKTree{})
}

package bk_tree

type modelOfEdge struct {
	Value int
	Elem  *modelOfNode
}

type modelOfNode struct {
	Value string
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

func (bktree *BKTree) Encode() any {
	var root_model modelOfNode
	if bktree.root != nil {
		root_model = bktree.root.Encode()
	}
	model := modelOfBKTree{Root: &root_model}
	return model
}

func (bktree *BKTree) Decode(model modelOfBKTree) {
	root := model.Root.Decode()
	bktree.root = &root
}

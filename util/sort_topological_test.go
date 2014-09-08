package util_test

import (
	"code.google.com/p/refola/util"
	"testing"
)

type graph []*node

type node struct {
	next, prev graph
	val        int
}

func (g *graph) contains(n *node) bool {
	for _, v := range *g {
		if v == n {
			return true
		}
	}
	return false
}

// implement util.TopologicalSortable
func (g *graph) Len() int {
	return len(*g)
}
func (g *graph) Compare(i, j int) util.Comparation {
	switch {
	case (*g)[i].next.contains((*g)[j]):
		return util.Less
	case (*g)[i].prev.contains((*g)[j]):
		return util.Greater
	default:
		return util.Other
	}
}
func (g *graph) Swap(i, j int) {
	(*g)[i], (*g)[j] = (*g)[j], (*g)[i]
}

type namedGraph struct {
	name string
	g    graph
}

// convert a lists of arrows into a graph
// WARNING: these are zero-indexed and going beyond len(arrows) is a panic
func makeGraph(arrows [][]int) *graph {
	g := make([]*node, len(arrows))
	for i := 0; i < len(arrows); i++ {
		g[i] = new(node)
	}
	for i, orig := range arrows {
		for _, dest := range orig {
			(*g)[i].to = append((*g)[i].to, (*g)[dest])
			(*g)[dest].from = append((*g)[dest], orig)
		}
	}
	return g
}

func TestTopological(t *testing.T) {
	// graphs made by drawing human-level pseudorandom arrows on the vertices of a polygon in KlourPaint with the randInt function of a TI-83 Plus used to shuffle the order
	// a connected graph with no loops
	connected := makeGraph()
	cycle := makeGraph()
	disconnected := makeGraph()

	graphs := []namedGraph{{"connected", connected}, {"cycle", cycle}, {"disconnected", disconnected}}
	var failed bool

	for _, v := range graphs {
		if err := testGraph(v); err != nil {
			failed = true
			t.Fail()
			t.Log("Failed sorting %s graph. Error: %s.", v.name, err)
		}
	}
}

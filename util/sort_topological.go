// Utility functions used by other (refola) packages
// Topological sort - given a collection with some set of things with partial ordering (with many orderings determined implicitly by intermediaries), arrange the collection s such that s[i]<s[j] implies i<j.

package util

import (
	"errors"
	"fmt"
	"sort"
)

type Comparation int

const (
	Less = Comparation(iota)
	Greater
	Other
)

// a listing of indices of nodes in a found cycle
type Cycle []int

type CycleError struct {
	Cycles []Cycle
}

func (ce CycleError) String() string {
	s := "Cycles found in given TopologicalSortable, with indices as follows:"
	for i, v := range ce.Cycles {
		s += fmt.Sprintf("\n%d: %s", i, v)
	}
	return s
}

// TODO: implement topological sort
// Resources:
//	http://en.wikipedia.org/wiki/Topological_sorting
//	http://en.wikipedia.org/wiki/Partial_sorting
//	http://math.stackexchange.com/questions/55891/algorithm-to-sort-based-on-a-partial-order
// Things needed for topological sorting of partially-ordered lists with some orders only implicit....
// TODO: shorter name...
type TopologicalSortable interface {
	Len() int                     // How many things to sort.
	Compare(i, j int) Comparation // i is [Less/Greater] than j or Other
	Swap(i, j int)                // Switch locations of elements i and j.
}

// Same as in the standard library's "sort" so we can pass it to them.
type Sortable interface {
	Len() int
	Less(i, j int) bool
	Swap(i, j int)
}

// just passes it through to the standard library's sort
func Sort(s Sortable) {
	sort.Sort(s)
}

type node struct {
	from []*node // nodes that this one comes from
	to   []*node // nodes that this one goes to (comes before)
	val  int     // array index that this node originally corresponded to
}

// convert the array into nodes so it can be topologically sorted
func buildGraph(t TopologicalSortable) []*node {
	l := t.Len()
	fmt.Println(l)
	nodes := make([]*node, l)
	// prebuild nodes list before referring to future nodes
	for i, _ := range nodes {
		nodes[i] = new(node)
		nodes[i].val = i
	}
	// loop over everything in t, building nodes describing the relationships between its elements
	// TODO: see if there's a way to do this in less than O(n²) time
	for i, _ := range nodes {
		// check each node after i for relation with i
		for j := i + 1; j < l; j++ {
			switch t.Compare(i, j) {
			case Less:
				nodes[i].to = append(nodes[i].to, nodes[j])
				nodes[j].from = append(nodes[j].from, nodes[i])
			case Greater:
				nodes[i].from = append(nodes[i].from, nodes[j])
				nodes[j].to = append(nodes[j].to, nodes[i])
			}
		}
	}
	return nodes
}

// Arrange the nodes so that their new indices have all their dependencies at lower (earlier) indices
// "Take first thing that doesn't need any currently untaken thing before it" rough algorithm taken from https://en.wikipedia.org/wiki/Topological_sorting and details filled in as encountered.
func placeNodes(n []*node) []*node {
	newN := make([]*node, len(n))
	for len(newN) < len(n) {
		for i, v := range n {
			if v.from == nil {
				// remove v from each v.to[foo].from
				for _, to := range v.to {
					for j, from := range to.from {
						if from == v {
							to.from[j] = nil
						}
					}
				}
				// move v to newN
				newN = append(newN, v)
				n[i] = nil
			}
		}
	}
	return newN
}

// Make the TopologicalSortable match the rearrangement given in the nodes.
func arrangeOriginal(graph []*node, t TopologicalSortable) {
	// Algorithm:
	// Iterate through nodes
	// 	At each node i:
	// 		Swap t's index i with whatever the node's original index got moved to
	// 		Record node movement

	// tracks where indices have moved to
	moved := make([]int, len(graph))
	for i := 0; i < len(graph); i++ {
		moved[i] = i
	}

	for i, v := range graph {
		t.Swap(i, v.val)
		moved[i], moved[v.val] = v.val, moved[i]
	}
}

// Sort according to a directed acyclic graph built from comparing elements
// TODO: make efficient for large t.Len() -- use different structure than array?
func Topological(t TopologicalSortable) error {
	graph := buildGraph(t)
	graph = placeNodes(graph)
	if t.Len() != len(graph) {
		panic("graph and TopologicalSortabel sizes differ")
	}
	arrangeOriginal(graph, t)

	// TODO: make sure this is stable or sort it separately afterwards

	// TODO: The only error this emits should be a CycleError
	return errors.New("Topological sort not yet implemented.")
}

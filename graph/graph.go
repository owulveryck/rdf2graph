package graph

import (
	"log"

	rdf "github.com/deiu/gon3"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
)

type dict map[string]rdf.Term

func (d dict) getOrInsert(t rdf.Term) rdf.Term {
	if t, ok := d[t.RawValue()]; ok {
		return t
	}
	d[t.RawValue()] = t
	return t
}

// NewGraph ...
func NewGraph(rdfGraph *rdf.Graph) Graph {

	// tree is a map of subjects, containing predicates referencing objects
	tree := make(map[rdf.Term]map[rdf.Term][]rdf.Term)
	termDict := dict(make(map[string]rdf.Term))
	for s := range rdfGraph.IterTriples() {
		subject := termDict.getOrInsert(s.Subject)
		predicate := termDict.getOrInsert(s.Predicate)
		object := termDict.getOrInsert(s.Object)
		if _, ok := tree[subject]; !ok {
			tree[subject] = make(map[rdf.Term][]rdf.Term)
		}
		sub := tree[subject]
		if sub[predicate] == nil {
			sub[predicate] = make([]rdf.Term, 0)
		}
		sub[predicate] = append(sub[predicate], object)
	}
	g := simple.NewDirectedGraph()
	reference := make(map[rdf.Term]*Node, len(tree))
	// create the nodes
	for s, po := range tree {
		n := &Node{
			id:             g.NewNode().ID(),
			Subject:        s,
			PredicatObject: po,
		}
		g.AddNode(n)
		reference[s] = n
	}
	// create the edges
	for s, po := range tree {

		me, ok := reference[s]
		if !ok {
			log.Fatal("wot? node is not found")
		}
		for predicate, objects := range po {
			for _, object := range objects {
				if to, ok := reference[object]; ok {
					if me == to {
						continue
					}
					e := Edge{
						F:    me,
						T:    to,
						Term: predicate,
					}
					g.SetEdge(e)
				}
			}

		}
	}
	return Graph{
		DirectedGraph: g,
		Reference:     reference,
	}
}

// Graph is carrying the information
type Graph struct {
	*simple.DirectedGraph
	// Reference of the term and their associated nodes
	Reference map[rdf.Term]*Node
}

// Node is a node of the graph, it carries n tuple associated with one subject
type Node struct {
	id             int64
	Subject        rdf.Term
	PredicatObject map[rdf.Term][]rdf.Term
}

// ID of the node
func (n *Node) ID() int64 {
	return n.id
}

// Edge ...
type Edge struct {
	F, T graph.Node
	Term rdf.Term
}

// From ...
func (e Edge) From() graph.Node {
	return e.F
}

// To ...
func (e Edge) To() graph.Node {
	return e.T
}

// ReversedEdge returns a new Edge with the F and T fields
// swapped.
func (e Edge) ReversedEdge() graph.Edge { return Edge{F: e.T, T: e.F, Term: e.Term} }

// FindNode returns a node whose Term match t's rawstring
// it returns nil if no matching node is found
func (g *Graph) FindNode(t rdf.Term) *Node {
	for term, n := range g.Reference {
		if term.Equals(t) {
			return n
		}
	}
	return nil
}

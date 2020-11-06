package main

import (
	"fmt"
	"log"
	"os"

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

func main() {
	// Set a base URI
	baseURI := "https://example.org/foo"
	//rdfType := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	// Create a new graph

	parser, err := rdf.NewParser(baseURI).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := NewGraph(parser)

	n, ok := findID(g.Reference, rdf.NewIRI("http://schema.org/PostalAddress"))
	if !ok {
		log.Fatal("not found")
	}
	fmt.Println(n.Subject)
	it := g.DirectedGraph.To(n.ID())
	for it.Next() {
		to := it.Node().(*node)
		e := g.Edge(to.ID(), n.ID()).(edge)
		fmt.Printf("%v -> %v -> %v\n", n.Subject, e.term, to.Subject)
	}
}

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
	reference := make(map[rdf.Term]*node, len(tree))
	// create the nodes
	for s, po := range tree {
		n := &node{
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
					e := edge{
						F:    me,
						T:    to,
						term: predicate,
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

type Graph struct {
	*simple.DirectedGraph
	Reference map[rdf.Term]*node
}

func findID(dict map[rdf.Term]*node, t rdf.Term) (*node, bool) {
	for term, n := range dict {
		if term.Equals(t) {
			return n, true
		}
	}
	return nil, false
}

type node struct {
	id             int64
	Subject        rdf.Term
	PredicatObject map[rdf.Term][]rdf.Term
}

func (n *node) ID() int64 {
	return n.id
}

type edge struct {
	F, T graph.Node
	term rdf.Term
}

func (e edge) From() graph.Node {
	return e.F
}

func (e edge) To() graph.Node {
	return e.T
}

// ReversedEdge returns a new Edge with the F and T fields
// swapped.
func (e edge) ReversedEdge() graph.Edge { return edge{F: e.T, T: e.F, term: e.term} }

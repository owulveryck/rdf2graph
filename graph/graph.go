package graph

import (
	"log"

	rdf "github.com/owulveryck/gon3"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

// Dict reference a term by its rawvalue
type Dict map[string]rdf.Term

func (d Dict) getOrInsert(t rdf.Term) rdf.Term {
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
	termDict := Dict(make(map[string]rdf.Term))
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
			id:              g.NewNode().ID(),
			Subject:         s,
			PredicateObject: po,
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
		Dict:          termDict,
		Reference:     reference,
	}
}

// Graph is carrying the information
type Graph struct {
	*simple.DirectedGraph
	// Reference of the term and their associated nodes
	Reference map[rdf.Term]*Node
	// Dict ...
	Dict map[string]rdf.Term
}

// ToWithEdges return al the nodes reaching n with an edge whosh subject is one of ts
func (g *Graph) ToWithEdges(n *Node, ts ...rdf.Term) graph.Nodes {
	nodes := make(map[int64]graph.Node)
	it := g.DirectedGraph.To(n.ID())
	for it.Next() {
		to := it.Node().(*Node)
		e := g.Edge(to.ID(), n.ID()).(Edge)
		for i := range ts {
			if e.Term.Equals(ts[i]) {
				nodes[to.ID()] = to
			}
		}
	}
	return iterator.NewNodes(nodes)
}

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

// GetTerm returns the term matching s (short forms are expanded according to the namespaces of the graph)
func (g *Graph) GetTerm(s string) rdf.Term {
	for k, v := range g.Dict {
		if k == s {
			return v
		}
	}
	return nil
}

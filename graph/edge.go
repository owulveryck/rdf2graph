package graph

import (
	rdf "github.com/owulveryck/gon3"
	"gonum.org/v1/gonum/graph"
)

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

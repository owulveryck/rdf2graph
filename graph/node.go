package graph

import rdf "github.com/owulveryck/gon3"

// Node is a node of the graph, it carries n tuple associated with one subject
type Node struct {
	id              int64
	Subject         rdf.Term
	PredicateObject map[rdf.Term][]rdf.Term
}

// GetObjectsFromPredicateIRI ...
//func (n *Node) GetObjectsFromPredicateIRI(s string) []rdf.Term {
//	return nil
//}

// ID of the node
func (n *Node) ID() int64 {
	return n.id
}

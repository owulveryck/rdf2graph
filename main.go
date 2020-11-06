package main

import (
	"fmt"
	"log"
	"os"

	rdf "github.com/deiu/gon3"
	"github.com/owulveryck/rdf2raph/graph"
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
	g := graph.NewGraph(parser)

	n, ok := g.FindNode(rdf.NewIRI("http://schema.org/PostalAddress"))
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

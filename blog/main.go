package main

import (
	"fmt"
	"log"
	"os"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
)

func main() {
	baseURI := "https://example.org/foo"
	parser := rdf.NewParser(baseURI)
	gr, err := parser.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := graph.NewGraph(gr)
	postalAddress := g.Dict["http://schema.org/PostalAddress"]
	node := g.Reference[postalAddress]
	it := g.To(node.ID())
	for it.Next() {
		n := it.Node().(*graph.Node) // need inference here because gonum's simple graph returns a type graph.Node which is an interface
		e := g.Edge(n.ID(), node.ID()).(graph.Edge)
		fmt.Printf("%v -%v-> %v\n", node.Subject, e.Term, n.Subject)
	}
}

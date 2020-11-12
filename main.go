package main

import (
	"fmt"
	"log"
	"os"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
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
	// Create a new graph

	parser := rdf.NewParser(baseURI)
	gr, err := parser.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := graph.NewGraph(gr)
	/*
		fmt.Println(parser.GetNamespaces())
		rdfType, ok := g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#type"]
		if !ok {
			log.Fatal("Term not present in the graph: ", "rdfType")
		}
		schemaOrgRangeInclude, ok := g.Dict["http://schema.org/rangeIncludes"]
		if !ok {
			log.Fatal("Term not present in the graph: ", "rangeInclude")
		}
		schemaOrgDomainIncludes, ok := g.Dict["http://schema.org/domainIncludes"]
		if !ok {
			log.Fatal("Term not present in the graph: ", "domainIncludes")
		}
		rdfsSubClassOf, ok := g.Dict["http://www.w3.org/2000/01/rdf-schema#subClassOf"]
		if !ok {
			log.Fatal("Term not present in the graph: ", "rdfsSubClassOf")
		}
		rdfProperty, ok := g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"]
		if !ok {
			log.Fatal("Term not present in the graph: ", "rdfProperty")
		}

		n := g.FindNode(rdf.NewIRI(os.Args[1]))
		if n == nil {
			log.Fatal("not found")
		}
		fmt.Println(n.Subject)
		for k, v := range n.PredicateObject {
			n := g.FindNode(k)
			if n == nil {
				fmt.Printf("\t%v: %v\n", k, v)
			}
		}
		it := g.DirectedGraph.To(n.ID())
		for it.Next() {
			to := it.Node().(*graph.Node)
			e := g.Edge(to.ID(), n.ID()).(graph.Edge)
			fmt.Printf("\t\t-> %v -> %v\n", e.Term, to.Subject)
		}
	*/
	b, err := dot.Marshal(g.DirectedGraph, "test", " ", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	/*
		_ = rdfProperty
		_ = rdfType
		_ = schemaOrgDomainIncludes
		_ = schemaOrgRangeInclude
		_ = rdfsSubClassOf
	*/
}

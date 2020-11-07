package main

import (
	"log"
	"net/http"
	"os"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
)

func main() {
	// Set a base URI
	baseURI := "https://example.org/foo"
	// Create a new graph

	parser, err := rdf.NewParser(baseURI).Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	g := graph.NewGraph(parser)
	h := &handler{
		g: g,
	}
	err = http.ListenAndServe(":8080", h)
	if err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	g graph.Graph
}

func (*handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

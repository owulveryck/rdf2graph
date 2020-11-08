package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
)

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
	h := &handler{
		namespaces: parser.GetNamespaces(),
		g:          g,
	}
	http.Handle("/", h)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	namespaces map[string]*rdf.IRI
	g          graph.Graph
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)
	n, err := h.getNode(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// get the type of the node
	rdfsType := h.g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#type"]
	rdfProperty := h.g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"]
	rdfClass := h.g.Dict["http://www.w3.org/2000/01/rdf-schema#Class"]
	typ, ok := n.PredicateObject[rdfsType]
	if !ok {
		http.Error(w, "Node has no type", http.StatusInternalServerError)
		return
	}
	switch typ[0] {
	case rdfProperty:
		fmt.Fprintf(w, "This is a property:\n %v", n)
	case rdfClass:
		fmt.Fprintf(w, "This is a class:\n %v", n)
	default:
		fmt.Fprintf(w, "This is something else:\n %v", n)
	}
}

func (h *handler) getNode(s string) (*graph.Node, error) {
	var n *graph.Node
	var term rdf.Term
	for k, v := range h.namespaces {
		if strings.Contains(s, k+":") {
			term = v
		}
	}
	if term == nil {
		return nil, fmt.Errorf("No matching namespace found for name %v", s)
	}
	colon := strings.Index(s, ":")
	s = term.RawValue() + s[colon+1:]
	term, ok := h.g.Dict[s]
	if !ok {
		return nil, fmt.Errorf("No term found for namespace %v", s)

	}
	n = h.g.FindNode(term)
	if n == nil {
		return nil, fmt.Errorf("No node found for term %v", term)
	}
	return n, nil

}

type propertyDisplay struct {
	URL         string
	Label       string
	Description string
	Types       []string
}

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
	w.Header().Add("content-type", "text/html")
	fmt.Fprint(w, "<!DOCTYPE html>")
	fmt.Fprint(w, "<html>")
	fmt.Fprint(w, "<head>")
	fmt.Fprint(w, `<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>`)
	fmt.Fprint(w, "</head>")
	fmt.Fprint(w, "<body>")
	fmt.Fprintf(w, "<h1>%v</h1>", n.Subject.RawValue())
	fmt.Fprint(w, `<table border="1" style="border: 1px #ccc;">`)
	fmt.Fprint(w, "<thead>")
	fmt.Fprint(w, `<tr><th>Property</th><th>Type</th><th><Description></th></tr>`)
	fmt.Fprint(w, "</thead>")
	h.processNode(n, w, nil)
	fmt.Fprint(w, "</table>")
	fmt.Fprint(w, "</body></html>")
}

func (h *handler) processNode(n *graph.Node, w http.ResponseWriter, visited []*graph.Node) {
	// get the type of the node
	rdfsType := h.g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#type"]
	rdfProperty := h.g.Dict["http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"]
	rdfClass := h.g.Dict["http://www.w3.org/2000/01/rdf-schema#Class"]
	rdfLabel := h.g.Dict["http://www.w3.org/2000/01/rdf-schema#label"]
	rdfComment := h.g.Dict["http://www.w3.org/2000/01/rdf-schema#comment"]
	rdfRangeIncludes := h.g.Dict["http://schema.org/rangeIncludes"]
	rdfDomainIncludes := h.g.Dict["http://schema.org/domainIncludes"]
	typ, ok := n.PredicateObject[rdfsType]
	if !ok {
		http.Error(w, "Node has no type", http.StatusInternalServerError)
		return
	}
	switch typ[0] {
	case rdfProperty:
		fmt.Fprintf(w, `<tr><td><a href="%v">%v</a>`,
			h.minifyHREF(n.Subject),
			n.PredicateObject[rdfLabel][0].RawValue())
		fmt.Fprint(w, "<td>")
		for _, v := range n.PredicateObject[rdfRangeIncludes] {
			fmt.Fprintf(w, `<a href="%v">%v</a><br>`, h.minifyHREF(v), h.minifyHREF(v))
		}
		fmt.Fprint(w, "</td>")
		fmt.Fprintf(w, "<td>%v</td>", n.PredicateObject[rdfComment][0].RawValue())
		fmt.Fprint(w, "</tr>")

	case rdfClass:
		fmt.Fprint(w, `<tbody style="border: 1px solid #ccc;">`)
		fmt.Fprintf(w, `<tr><td colspan="3"><b>Properties from: <a href="%v">%v</a></b> <pre>%v</pre></td></tr>`,
			h.minifyHREF(n.Subject),
			h.minifyHREF(n.Subject),
			n.PredicateObject[rdfComment][0].RawValue(),
		)
		// GetProperties
		it := h.g.DirectedGraph.To(n.ID())
		for it.Next() {
			nn := it.Node().(*graph.Node)
			e := h.g.DirectedGraph.Edge(nn.ID(), n.ID()).(graph.Edge)
			if e.Term.Equals(rdfDomainIncludes) {
				if nn.PredicateObject[rdfsType][0].Equals(rdfProperty) {
					h.processNode(nn, w, nil)
				}
			}
		}
		// GetClasses
		it = h.g.DirectedGraph.From(n.ID())
		for it.Next() {
			n := it.Node().(*graph.Node)
			if n.PredicateObject[rdfsType][0].Equals(rdfClass) {
				h.processNode(n, w, nil)
			}
		}
		//var class *graph.Node
		//h.processNode(prop,w)

		fmt.Fprint(w, "</tbody>")
	default:
		fmt.Fprintf(w, "This is something else:\n %v", n)
	}

}

// minifyHREF returns a string with the nameprefix
func (h *handler) minifyHREF(t rdf.Term) string {
	found := false
	for k := range h.g.Reference {
		if k.Equals(t) {
			found = true
		}
	}
	if !found {
		return t.RawValue()
	}
	for k, v := range h.namespaces {
		if strings.Contains(t.RawValue(), v.RawValue()) {
			return "/" + strings.Replace(t.RawValue(), v.RawValue(), k+":", -1)
		}
	}
	return t.RawValue()
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

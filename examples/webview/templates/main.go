package main

import (
	"html/template"
	"log"
	"os"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
)

// Current object to apply to the template
type Current struct {
	Graph *graph.Graph
	Node  *graph.Node
}

// HasPredicate returns true if the couple predicate/object exists in the tuple
//
// ex:
//   {{- if .HasPredicate "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" "http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"}}{{end}}
func (g Current) HasPredicate(predicate, object string) bool {
	predicateT := g.linkToTerm(predicate)
	objectT := g.linkToTerm(object)
	if _, ok := g.Node.PredicateObject[predicateT]; !ok {
		return false
	}
	for _, o := range g.Node.PredicateObject[predicateT] {
		if o.Equals(objectT) {
			return true
		}
	}
	return false
}

// Objects return all the objects for the predicate in the current node
//
// to return the label of an object
//    {{ index (.Objects "http://www.w3.org/2000/01/rdf-schema#label") 0 }}
func (g Current) Objects(predicate string) []string {
	predicateT := g.linkToTerm(predicate)
	output := make([]string, 0)
	if _, ok := g.Node.PredicateObject[predicateT]; !ok {
		return output
	}
	for _, o := range g.Node.PredicateObject[predicateT] {
		output = append(output, o.RawValue())
	}
	return output
}

// To the node with edge holding a value fro m links
func (g Current) To(links ...string) []Current {
	n := g.Node
	linksT := g.linksToTerms(links...)
	ns := make([]Current, 0)
	it := g.Graph.DirectedGraph.To(n.ID())
	for it.Next() {
		from := it.Node().(*graph.Node)
		e := g.Graph.Edge(from.ID(), n.ID()).(graph.Edge)
		if len(links) == 0 {
			ns = append(ns, Current{
				Graph: g.Graph,
				Node:  from,
			})

		} else {
			for _, link := range linksT {
				if e.Term.Equals(link) {
					ns = append(ns, Current{
						Graph: g.Graph,
						Node:  from,
					})
				}
			}
		}
	}
	return ns
}

// From the node with edge holding a value fro m links
func (g Current) From(links ...string) []Current {
	n := g.Node
	linksT := g.linksToTerms(links...)
	ns := make([]Current, 0)
	it := g.Graph.DirectedGraph.From(n.ID())
	for it.Next() {
		to := it.Node().(*graph.Node)
		e := g.Graph.Edge(n.ID(), to.ID()).(graph.Edge)
		if len(links) == 0 {
			ns = append(ns, Current{
				Graph: g.Graph,
				Node:  to,
			})

		} else {
			for _, link := range linksT {
				if e.Term.Equals(link) {
					ns = append(ns, Current{
						Graph: g.Graph,
						Node:  to,
					})
				}
			}
		}
	}
	return ns
}
func (g *Current) linkToTerm(s string) rdf.Term {
	if t, ok := g.Graph.Dict[s]; ok {
		return t
	}
	log.Fatal("No term found for ", s)
	return nil
}

func (g *Current) linksToTerms(links ...string) []rdf.Term {
	terms := make([]rdf.Term, len(links))
	for i := range links {
		terms[i] = g.linkToTerm(links[i])
	}
	return terms
}

// Has s as a predicate
func main() {
	tmpl, err := template.New("render").ParseFiles("test.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	baseURI := "https://example.org/foo"
	f, err := os.Open("../../../client/schemaorg-current-http.ttl")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Create a new graph

	parser := rdf.NewParser(baseURI)
	gr, err := parser.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	g := graph.NewGraph(gr)
	it := g.Nodes()
	for it.Next() {
		n := it.Node().(*graph.Node)
		if n.Subject.RawValue() == "http://schema.org/PostalAddress" {
			tmpl.ExecuteTemplate(os.Stdout, "main", Current{
				Graph: &g,
				Node:  n,
			})
		}
	}
}

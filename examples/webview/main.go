package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
	rdfTmpl "github.com/owulveryck/rdf2graph/template"
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
	b, err := ioutil.ReadFile("index.tmpl")
	if err != nil {
		log.Fatal(err)

	}
	funcMap := template.FuncMap{
		"minifyhref": h.MinifyHREF,
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(string(b))
	h.tmpl = tmpl
	http.Handle("/", h)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type handler struct {
	namespaces map[string]*rdf.IRI
	g          graph.Graph
	tmpl       *template.Template
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)
	n, err := h.getNode(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("content-type", "text/html")
	err = rdfTmpl.Apply(w, h.tmpl, "main", n, &h.g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (h *handler) MinifyHREF(s string) string {
	found := false
	for k := range h.g.Reference {
		if k.RawValue() == s {
			found = true
		}
	}
	if !found {
		return s
	}
	for k, v := range h.namespaces {
		if strings.Contains(s, v.RawValue()) {
			return "/" + strings.Replace(s, v.RawValue(), k+":", -1)
		}
	}
	return s
}

func (h *handler) getNode(s string) (string, error) {
	var n *graph.Node
	var term rdf.Term
	for k, v := range h.namespaces {
		if strings.Contains(s, k+":") {
			term = v
		}
	}
	if term == nil {
		return "", fmt.Errorf("No matching namespace found for name %v", s)
	}
	colon := strings.Index(s, ":")
	s = term.RawValue() + s[colon+1:]
	term, ok := h.g.Dict[s]
	if !ok {
		return "", fmt.Errorf("No term found for namespace %v", s)

	}
	n = h.g.FindNode(term)
	if n == nil {
		return "", fmt.Errorf("No node found for term %v", term)
	}
	return term.RawValue(), nil

}

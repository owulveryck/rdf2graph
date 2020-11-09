package template

import (
	"html/template"
	"os"
	"strings"

	rdf "github.com/owulveryck/gon3"
	"github.com/owulveryck/rdf2graph/graph"
)

func ExampleApply() {

	const templ = `
{{- define "http://www.w3.org/1999/02/22-rdf-syntax-ns#type http://www.w3.org/1999/02/22-rdf-syntax-ns#Property" -}}
property {{ .Node.Subject.RawValue }} {{ index (.Objects "http://www.w3.org/2000/01/rdf-schema#label") 0 }}
{{- end }}

{{- define "http://www.w3.org/1999/02/22-rdf-syntax-ns#type http://www.w3.org/2000/01/rdf-schema#Class" -}}
Class: {{ .Node.Subject.RawValue }}
{{- range .To }}
{{- if .HasPredicate "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" "http://www.w3.org/1999/02/22-rdf-syntax-ns#Property"}}
{{ template "http://www.w3.org/1999/02/22-rdf-syntax-ns#type http://www.w3.org/1999/02/22-rdf-syntax-ns#Property" . -}}
{{ end -}}
{{ end -}}
{{ end -}}

{{ define "main" }}
start
{{ template "http://www.w3.org/1999/02/22-rdf-syntax-ns#type http://www.w3.org/2000/01/rdf-schema#Class" . }}
end
{{ end }}
	`
	const ontology = `
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix example: <http://example.com/> .

# a is a shortcut for http://www.w3.org/1999/02/22-rdf-syntax-ns#type
example:PostalAddress a rdfs:Class ;
    rdfs:label "PostalAddress" ;
    rdfs:comment "The mailing address." .
	
example:addressCountry a rdf:Property ;
    rdfs:label "addressCountry" ;
    rdfs:domain example:PostalAddress ;
	rdfs:comment "A comment" .
	
example:address a rdf:Property ;
    rdfs:label "address" ;
	rdfs:domain example:PostalAddress ;
	rdfs:comment "Physical address of the item." .
	`

	baseURI := "https://example.org/foo"
	parser := rdf.NewParser(baseURI)
	gr, _ := parser.Parse(strings.NewReader(ontology))
	g := graph.NewGraph(gr)

	tmpl, _ := template.New("global").Parse(templ)
	Apply(os.Stdout, tmpl, "main", "http://example.com/PostalAddress", &g)
	// output:
	// start
	// Class: http://example.com/PostalAddress
	// property http://example.com/addressCountry addressCountry
	// property http://example.com/address address
	// end
}

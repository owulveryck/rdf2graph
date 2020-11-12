package graph

import (
	"bytes"
	"html"
	"log"
	"strings"
	"text/template"

	"gonum.org/v1/gonum/graph/encoding"
)

// Attributes fulfills the gonum encoding Attributer
func (n *Node) Attributes() []encoding.Attribute {
	replaceSlash := func(s string) string {
		return strings.Replace(s, `/`, `\/`, -1)
	}
	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"wrap":       wrap,
		"htmlescape": html.EscapeString,
		"slash":      replaceSlash,
	}
	t := template.Must(template.New("node").Funcs(funcMap).Parse(tmpl))
	var b bytes.Buffer
	err := t.Execute(&b, n)
	if err != nil {
		log.Fatal(err)
	}

	shape := encoding.Attribute{
		Key:   "shape",
		Value: "plain",
	}
	label := encoding.Attribute{
		Key: "label",
		//Value: n.Subject.RawValue(),
		Value: b.String(),
	}
	/*	type temp struct {
			Subject string
			Val     map[string]string
		}
	*/
	return []encoding.Attribute{shape, label}
}

// Attributes fulfills the gonum encoding Attributer
func (e Edge) Attributes() []encoding.Attribute {
	shape := encoding.Attribute{
		Key:   "shape",
		Value: "plain",
	}
	label := encoding.Attribute{
		Key:   "label",
		Value: e.Term.String(),
	}
	fontsize := encoding.Attribute{
		Key:   "fontsize",
		Value: "10",
	}
	return []encoding.Attribute{label, fontsize, shape}
}

func wrap(text string) string {
	lineWidth := 40
	words := strings.Fields(strings.TrimSpace(text))
	if len(words) == 0 {
		return text
	}
	wrapped := words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if len(word)+1 > spaceLeft {
			wrapped += `<br/>` + word
			spaceLeft = lineWidth - len(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + len(word)
		}
	}
	return wrapped
}

const tmpl = `<<table border='1' cellborder='1'>
	   <tr><td colspan="2" bgcolor="lightblue">{{ .Subject.RawValue | htmlescape}}</td></tr>
	   <tr><td bgcolor="lightgrey">Predicate</td><td bgcolor="lightgrey">object</td></tr>
	   {{range $key, $value := .PredicateObject -}}
	   {{range $value -}}
	   <tr><td>{{ $key.RawValue | htmlescape }}</td><td>{{.RawValue | htmlescape | wrap }}</td></tr>
	   {{end -}}
	   {{end -}}
	 </table>>`

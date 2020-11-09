package graph

import "gonum.org/v1/gonum/graph/encoding"

// Attributes fulfills the gonum encoding Attributer
func (n *Node) Attributes() []encoding.Attribute {
	shape := encoding.Attribute{
		Key:   "shape",
		Value: "record",
	}
	label := encoding.Attribute{
		Key:   "label",
		Value: n.Subject.RawValue(),
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
	label := encoding.Attribute{
		Key:   "label",
		Value: e.Term.String(),
	}
	fontsize := encoding.Attribute{
		Key:   "fontsize",
		Value: "10",
	}
	return []encoding.Attribute{label, fontsize}
}

const tmpl = `<
	 <table border='1' cellborder='0'>
	   {{range $key, $value := .Val}}
	   <tr><td>{{$key}}</td><td>{{$value}}</td></tr>
	   {{end}}
     </table>

>`

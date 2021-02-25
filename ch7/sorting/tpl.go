package main

import "html/template"

var (
	rowTpl,
	headerTpl,
	tableTpl *template.Template
)

func init() {

	headerTpl = template.Must(template.New("theader").Parse(`
	<tr>
	{{range .Headers}}
		<td>
			<a href="/?order={{.}} title="Sort by {{.}}">{{.}}</a>
		</td>
	{{end}}
	</tr>`))
	rowTpl = template.Must(template.New("theader").Parse(`
	<tr>
		<td>
			{{.Title}}
		</td>
		<td>
			{{.Artist}}
		</td>
		<td>
			{{.Album}}
		</td>
		<td>
			{{.Year}}
		</td>
		<td>
			{{.Length}}
		</td>
	</tr>`))
}

package main

import "html/template"

var (
	htmlTpl *template.Template
)

func init() {
	htmlTpl = template.Must(template.New("page").Parse(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Tracks</title>
		<style>
			th,td {
				text-align: left;
				padding-left: 32px;
			}
			tr:hover{
				background-color: #ffdfb3
			}
		</style>
	</head>
	<body>{{template "table" .}}</body>
	</html>`))

	htmlTpl.New("table").Parse(`<table>
	{{template "header" .}}
	{{range .Tracks}}{{template "row" .}}{{end}}
	</table>`)

	htmlTpl.New("header").Parse(`
	{{$headers := .Headers}}
	{{$ordered := .OrderedField}}
	{{$order := .Order}}
	<tr>
	{{range $headers}}
		<th>
		{{if eq $ordered . -}}
			<a href="/?field={{.}}&order={{$order}}" 
				title="Sort by {{.}}{{$order}}ending.">{{.}}
				{{if eq $order "asc"}}&#x1F815;{{else}}&#x1F807;{{end}}</a>

		{{else}}
			<a href="/?field={{.}}" title="Sort by {{.}}">{{.}}</a>
		{{- end}}
		</th>
	{{end}}
	</tr>`)

	htmlTpl.New("row").Parse(`
	<tr>
		<td>{{.Title}}</td>
		<td>{{.Artist}}</td>
		<td>{{.Album}}</td>
		<td>{{.Year}}</td>
		<td>{{.Length}}</td>
	</tr>`)
}

package main

var cmdTemplate = `
The supported commands are:
{{range .}}{{if .CanRun}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "ep help [action]" for detailed information on the usage for any given command.
`

var podcastTemplate = `
Here are the currently added podcasts by tag:

{{range .}}
	{{.Tag | printf "%-11s"}} {{.Name}}{{end}}

`

var epissodesTemplate = `
Here are the currently stored podcast episodes:

{{range $index, $element := .}}
	{{box $index | printf "%-5s"}} {{$element.Title}}{{end}}

`
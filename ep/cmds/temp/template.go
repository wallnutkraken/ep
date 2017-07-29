package temp

import (
	"os"
	"strconv"
	"text/template"
)

var CmdTemplate = `
The supported commands are:
{{range .}}{{if .CanRun}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "ep help [action]" for detailed information on the usage for any given command.
`

var PodcastTemplate = `
Here are the currently added podcasts by tag:

{{range .}}
	{{.Tag | printf "%-11s"}} {{.Name}}{{end}}

`

var EpissodesTemplate = `
Here are the currently stored podcast episodes:

{{range $index, $element := .}}
	{{box $index | printf "%-5s"}} {{$element.Title}}{{end}}

`


func WriteTemplate(text string, data interface{}) error {
	temp := template.New("top")
	temp.Funcs(template.FuncMap{"box": boxAndIncrement})
	template.Must(temp.Parse(text))

	return temp.Execute(os.Stdout, data)
}

func boxAndIncrement(num int) string {
	return "[" + strconv.Itoa(num+1) + "]"
}
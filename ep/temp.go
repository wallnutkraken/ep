package main

import (
	"os"
	"strconv"
	"text/template"
)

func writeTemplate(text string, data interface{}) error {
	temp := template.New("top")
	temp.Funcs(template.FuncMap{"box": boxAndIncrement})
	template.Must(temp.Parse(text))

	return temp.Execute(os.Stdout, data)
}

func boxAndIncrement(num int) string {
	return "[" + strconv.Itoa(num+1) + "]"
}

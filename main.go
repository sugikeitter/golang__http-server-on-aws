package main

import (
	"html/template"
	"net/http"
	"time"
)

var tmpl = template.Must(template.New("index.html").
	//Parse("<html><body><div>{{.Time}}</div>{{.Name}}さん</body></html>"))
	ParseFiles("template/index.html"))

type Result struct {
	Time string
	Name string
}

func handler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		name = "Nanashi"
	}
	result := Result{
		Time: currentTime(),
		Name: name,
	}
	tmpl.Execute(w, result)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func currentTime() string {
	t := time.Now()
	return t.Format("2006/01/02 15:04:05.000")
}

package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("index.html")
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.Execute(w, nil)
	})
	counter := 0
	http.HandleFunc("/increment", func(w http.ResponseWriter, r *http.Request) {
		tmplStr := `<div id="counter">{{ . }}</div>`
		tmpl := template.Must(template.New("counter").Parse(tmplStr))
		tmpl.ExecuteTemplate(w, "counter", counter)
		counter++
	})

	// Start the HTTP server on port 8080
	return http.ListenAndServe(":8080", nil)
}

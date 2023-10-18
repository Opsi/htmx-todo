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

type Todo struct {
	ID    int
	Title string
}

func run() error {

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("parse templates: %w", err)
	}

	todos := []Todo{
		{ID: 1, Title: "Buy groceries"},
		{ID: 2, Title: "Finish homework"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = tmpl.ExecuteTemplate(w, "index.html", todos)
		if err != nil {
			fmt.Println(err)
			return
		}
	})

	// Start the HTTP server on port 8080
	return http.ListenAndServe(":8080", nil)
}

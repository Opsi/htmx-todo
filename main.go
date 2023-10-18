package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

type Todo struct {
	ID      int
	Title   string
	Checked bool
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

	http.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		todoIDRegex := regexp.MustCompile(`/todo/(\d+)$`)
		matches := todoIDRegex.FindStringSubmatch(r.URL.Path)
		if len(matches) < 2 {
			http.Error(w, "invalid todo id", http.StatusBadRequest)
			return
		}
		todoID, err := strconv.Atoi(matches[1])
		if err != nil {
			http.Error(w, "invalid todo id", http.StatusBadRequest)
			return
		}
		todoIndex := slices.IndexFunc(todos, func(t Todo) bool {
			return t.ID == todoID
		})
		if todoIndex == -1 {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}
		todo := todos[todoIndex]
		todo.Checked = !todo.Checked
		todos[todoIndex] = todo
		tmpl.ExecuteTemplate(w, "todo.html", todo)
	})

	// Start the HTTP server on port 8080
	return http.ListenAndServe(":8080", nil)
}

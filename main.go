package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func readTemplates() (*template.Template, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	return tmpl, nil
}

func extractTodoID(r *http.Request) (int, error) {
	todoIDRegex := regexp.MustCompile(`/todo/(\d+)$`)
	matches := todoIDRegex.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid todo id")
	}
	todoID, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid todo id")
	}
	return todoID, nil
}

func run() error {
	todoRepo := NewTodoRepo()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := readTemplates()
		if err != nil {
			fmt.Println(err)
			return
		}
		tmpl.ExecuteTemplate(w, "index.html", todoRepo.GetAll())
	})

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := readTemplates()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		switch {
		case r.Method == http.MethodPost:
			var todo Todo
			if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			todo = todoRepo.Create(todo.Title)
			tmpl.ExecuteTemplate(w, "todo.html", todo)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		todoID, err := extractTodoID(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tmpl, err := readTemplates()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		switch {
		case r.Method == http.MethodGet:
			todo, ok := todoRepo.Get(todoID)
			if !ok {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			tmpl.ExecuteTemplate(w, "todo.html", todo)
			return
		case r.Method == http.MethodPut:
			var update TodoUpdate
			if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			todo, ok := todoRepo.Update(todoID, update)
			if !ok {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			tmpl.ExecuteTemplate(w, "todo.html", todo)
			return
		case r.Method == http.MethodDelete:
			if ok := todoRepo.Delete(todoID); !ok {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server on port 8080
	return http.ListenAndServe(":8080", nil)
}

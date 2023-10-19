package main

import (
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

func extractInt(r *http.Request, regex string) (int, error) {
	todoIDRegex := regexp.MustCompile(regex)
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

	http.HandleFunc("/edittodo/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		todoID, err := extractInt(r, `/edittodo/(\d+)$`)
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
		todo, ok := todoRepo.Get(todoID)
		if !ok {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}
		tmpl.ExecuteTemplate(w, "edittodo.html", todo)
	})

	http.HandleFunc("/toggletodo/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		todoID, err := extractInt(r, `/toggletodo/(\d+)$`)
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
		todo, ok := todoRepo.Toggle(todoID)
		if !ok {
			http.Error(w, "todo not found", http.StatusNotFound)
			return
		}
		tmpl.ExecuteTemplate(w, "todo.html", todo)
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
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			title := r.Form.Get("title")
			if title == "" {
				http.Error(w, "title cannot be empty", http.StatusBadRequest)
				return
			}
			todo := todoRepo.Create(title)
			tmpl.ExecuteTemplate(w, "todo.html", todo)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		todoID, err := extractInt(r, `/todo/(\d+)$`)
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
			// encode form data
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			title := r.Form.Get("title")
			if title == "" {
				http.Error(w, "title cannot be empty", http.StatusBadRequest)
				return
			}
			todo, ok := todoRepo.Update(todoID, title)
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

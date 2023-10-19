package main

import "sync"

type Todo struct {
	ID      int
	Title   string
	Checked bool
}

type TodoUpdate struct {
	Title   *string
	Checked *bool
}

type todoRepo struct {
	mu    sync.Mutex
	todos []Todo
	id    int
}

var _ TodoRepo = (*todoRepo)(nil)

type TodoRepo interface {
	GetAll() []Todo
	Get(id int) (Todo, bool)
	Create(title string) Todo
	Update(id int, update TodoUpdate) (Todo, bool)
	Delete(id int) bool
}

func NewTodoRepo() TodoRepo {
	return &todoRepo{
		todos: []Todo{
			{ID: 1, Title: "Buy groceries"},
			{ID: 2, Title: "Finish homework"},
		},
		id: 3,
	}
}

func (r *todoRepo) GetAll() []Todo {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Todo{}, r.todos...)
}

func (r *todoRepo) Get(id int) (Todo, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, true
		}
	}
	return Todo{}, false
}

func (r *todoRepo) Create(title string) Todo {
	r.mu.Lock()
	defer r.mu.Unlock()
	newTodo := Todo{
		ID:      r.id,
		Title:   title,
		Checked: false,
	}
	r.todos = append(r.todos, newTodo)
	r.id++
	return newTodo
}

func (r *todoRepo) Update(id int, update TodoUpdate) (Todo, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, todo := range r.todos {
		if todo.ID != id {
			continue
		}
		if update.Title != nil {
			r.todos[i].Title = *update.Title
		}
		if update.Checked != nil {
			r.todos[i].Checked = *update.Checked
		}
		return r.todos[i], true
	}
	return Todo{}, false
}

func (r *todoRepo) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, todo := range r.todos {
		if todo.ID != id {
			continue
		}
		r.todos = append(r.todos[:i], r.todos[i+1:]...)
		return true
	}
	return false
}

package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tearingItUp786/go-lang-todo/models"
	"github.com/tearingItUp786/go-lang-todo/templates"
	"github.com/tearingItUp786/go-lang-todo/views"
)

// BaseHandler
type BaseHandler struct {
	todoService  *models.BaseModel
	homeTemplate Template
	todoTemplate Template
}

type BaseHandlerInput struct {
	TodoService *models.BaseModel
}

// NewBaseHandler returns a new BaseHandler
func NewTodoContaroller(baseInput BaseHandlerInput) *BaseHandler {
	homeTemplate := views.Must(views.ParseFS(
		templates.FS,
		"index.gohtml", "template.gohtml", "todo-templates.gohtml",
	))

	todoTemplate := views.Must(views.ParseFS(
		templates.FS,
		"todo-templates.gohtml",
	))

	return &BaseHandler{
		todoService:  baseInput.TodoService,
		homeTemplate: homeTemplate,
		todoTemplate: todoTemplate,
	}
}

type Data struct {
	ToDos []models.ToDo
}

func (h *BaseHandler) GetToDos(w http.ResponseWriter, r *http.Request) {
	h.todoService.InsertToDo(
		fmt.Sprintf("Hello World %v", time.Now().Format("2006-01-02, 15:04:05")),
	)
	todos, err := h.todoService.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := Data{
		ToDos: todos,
	}
	h.homeTemplate.Execute(w, r, data)
}

func (h *BaseHandler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, err := h.todoService.ToggleTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.todoTemplate.ExecuteTemplate(w, r, "swap-todo", row)
}

func (h *BaseHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.todoService.DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write([]byte(""))
}

func (h *BaseHandler) GetEditToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	todo, err := h.todoService.GetSingleToDo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.todoTemplate.ExecuteTemplate(w, r, "edit-todo", todo)
}

func (h *BaseHandler) PatchEditToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	oldTodo, err := h.todoService.GetSingleToDo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newToDoDone := r.FormValue("todo-done")
	toDoDoneAsBool := newToDoDone == "true"
	newToDoText := r.FormValue("todo-text")

	todo, err := h.todoService.UpdateSingleToDo(id, newToDoText, toDoDoneAsBool)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if oldTodo.Done != todo.Done {

		h.todoTemplate.ExecuteTemplate(w, r, "swap-todo", todo)
		return
	}

	h.todoTemplate.ExecuteTemplate(w, r, "swap-single", todo)
}

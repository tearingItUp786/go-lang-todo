package controllers

import (
	"net/http"

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
		"index.gohtml", "template.gohtml", "todo/normal.gohtml",
	))

	todoTemplate := views.Must(views.ParseFS(
		templates.FS,
		"render.gohtml", "todo/normal.gohtml",
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

	h.todoTemplate.Execute(w, r, row)
}

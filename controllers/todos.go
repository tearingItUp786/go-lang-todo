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
func NewTodoController(baseInput BaseHandlerInput) *BaseHandler {
	homeTemplate := views.Must(views.ParseFS(
		templates.FS,
		"index.gohtml", "template.gohtml", "base.gohtml", "todo-templates.gohtml",
	))

	todoTemplate := views.Must(views.ParseFS(
		templates.FS,
		"base.gohtml", "todo-templates.gohtml",
	))

	return &BaseHandler{
		todoService:  baseInput.TodoService,
		homeTemplate: homeTemplate,
		todoTemplate: todoTemplate,
	}
}

type EnhancedToDo struct {
	models.ToDo
	Error bool
}

func NewEnhancedToDo(todo models.ToDo) EnhancedToDo {
	return EnhancedToDo{
		ToDo: models.ToDo{
			Id:   todo.Id,
			Text: todo.Text,
			Done: todo.Done,
		},
		Error: false,
	}
}

type Data struct {
	ToDos []EnhancedToDo
}

func (h *BaseHandler) GetToDos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoService.GetTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	enhancedToDos := []EnhancedToDo{}
	for _, todo := range todos {
		todo := NewEnhancedToDo(todo)
		enhancedToDos = append(enhancedToDos, todo)
	}

	data := Data{
		ToDos: enhancedToDos,
	}
	h.homeTemplate.Execute(w, r, data)
}

func (h *BaseHandler) NewTodo(w http.ResponseWriter, r *http.Request) {
	todoText := r.FormValue("todo-text")
	row, count, err := h.todoService.InsertToDo(todoText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	enhancedRow := NewEnhancedToDo(row)
	// let's get the base html if we just added "one" todo
	if count <= 1 {
		data := Data{
			ToDos: []EnhancedToDo{enhancedRow},
		}
		h.todoTemplate.ExecuteTemplate(w, r, "add-new-todo-swap", data)
		// h.GetToDos(w, r)
		return
	}

	// we have a bunch of todos and we only want to swap out the
	// content inside of the incomplete list
	h.todoTemplate.ExecuteTemplate(w, r, "swap-todo", enhancedRow)
}

func (h *BaseHandler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	row, err := h.todoService.ToggleTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.todoTemplate.ExecuteTemplate(w, r, "swap-todo", NewEnhancedToDo(row))
}

func (h *BaseHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	count, err := h.todoService.DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if count > 0 {
		w.Write([]byte(""))
		return
	}

	h.todoTemplate.ExecuteTemplate(w, r, "empty-list", nil)
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
	oldToDo, err := h.todoService.GetSingleToDo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	enhancedToDo := NewEnhancedToDo(oldToDo)
	newToDoDone := r.FormValue("todo-done")
	newToDoText := r.FormValue("todo-text")
	toDoDoneAsBool := newToDoDone == "true"

	if newToDoText == "error" {
		enhancedToDo.Error = true

		h.todoTemplate.ExecuteTemplate(
			w,
			r,
			"swap-single",
			enhancedToDo,
		)
		return
	}

	todo, err := h.todoService.UpdateSingleToDo(id, newToDoText, toDoDoneAsBool)
	enhancedToDo = NewEnhancedToDo(todo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if oldToDo.Done != todo.Done {
		h.todoTemplate.ExecuteTemplate(w, r, "swap-todo", enhancedToDo)
		return
	}

	h.todoTemplate.ExecuteTemplate(w, r, "swap-single", enhancedToDo)
}

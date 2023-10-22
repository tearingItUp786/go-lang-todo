package controllers

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
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

func (h *BaseHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.todoService.DeleteAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.todoTemplate.ExecuteTemplate(w, r, "empty-list", nil)
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

func (h *BaseHandler) BulkUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseMultipartForm(32 << 20) // Max memory to allocate for the form data
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, _, err := r.FormFile("csv") // Name of the file input field in the form
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	records, err := readCSV(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Process the CSV records
	bulkTodos := []models.ToDo{}
	for i := 1; i < len(records); i++ {
		record := records[i]
		fmt.Println(record)
		done := record[1] == "true"
		bulkTodos = append(bulkTodos, models.ToDo{Text: record[0], Done: done})
	}

	// You can also write the CSV data to a file if needed
	// writeFile(records)
	todos, err := h.todoService.BulkInsertToDos(bulkTodos)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	enhancedToDos := []EnhancedToDo{}
	for _, todo := range todos {
		todo := NewEnhancedToDo(todo)
		enhancedToDos = append(enhancedToDos, todo)
	}

	data := Data{
		ToDos: enhancedToDos,
	}
	h.todoTemplate.ExecuteTemplate(w, r, "add-new-todo-swap", data)
}

func readCSV(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

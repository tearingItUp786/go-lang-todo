package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/tearingItUp786/go-lang-todo/models"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	db *sql.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

type Data struct {
	ToDos []models.ToDo
}

func (h *BaseHandler) GetToDos(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	baseModel := models.NewBaseModel(h.db)

	todos, _ := baseModel.GetTodos()
	fmt.Println(todos)

	tmpl, err := template.ParseFiles("views/index.gohtml", "views/template.gohtml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := Data{
		ToDos: todos,
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

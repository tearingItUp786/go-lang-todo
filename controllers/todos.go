package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/tearingItUp786/go-lang-todo/models"
	"github.com/tearingItUp786/go-lang-todo/templates"
	"github.com/tearingItUp786/go-lang-todo/views"
)

// BaseHandler
type BaseHandler struct {
	db           *sql.DB
	homeTemplate Template
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
		homeTemplate: views.Must(views.ParseFS(
			templates.FS,
			"index.gohtml", "template.gohtml",
		)),
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

	data := Data{
		ToDos: todos,
	}
	h.homeTemplate.Execute(w, r, data)
}

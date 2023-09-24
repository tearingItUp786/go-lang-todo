package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tearingItUp786/go-lang-todo/controllers"
	"github.com/tearingItUp786/go-lang-todo/models"
)

func main() {
	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	h := controllers.NewBaseHandler(db)
	router := chi.NewRouter()
	router.Get("/", h.GetToDos)

	http.ListenAndServe(":3000", router)
}

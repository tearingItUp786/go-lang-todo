package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/tearingItUp786/go-lang-todo/controllers"
	"github.com/tearingItUp786/go-lang-todo/models"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := models.Open(models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	})
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error connecting to database")
	}

	defer db.Close()

	todoController := controllers.NewBaseHandler(db)

	router := chi.NewRouter()
	router.Get("/", todoController.GetToDos)

	http.ListenAndServe(":3000", router)
}

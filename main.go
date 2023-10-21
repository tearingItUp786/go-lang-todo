package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/tearingItUp786/go-lang-todo/controllers"
	"github.com/tearingItUp786/go-lang-todo/migrations"
	"github.com/tearingItUp786/go-lang-todo/models"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed dist/*
var jsComponents embed.FS

func main() {
	env := os.Getenv("FOO_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env

	db, err := models.Open(models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	})

	// todo: get this from env
	csrfKey := os.Getenv("CSRF_KEY")
	csrfSecure := os.Getenv("CSRF_SECURE") == "true"

	if err != nil {
		fmt.Println(err)
		log.Fatal("Error connecting to database")
	}

	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Secure(csrfSecure),
	)
	todoService := models.NewBaseModel(db)

	todoController := controllers.NewTodoController(
		controllers.BaseHandlerInput{TodoService: todoService},
	)

	router := chi.NewRouter()

	router.Use(csrfMw)
	router.Get("/", todoController.GetToDos)
	router.Post("/new", todoController.NewTodo)

	subRouter := chi.NewRouter()
	subRouter.Route("/{id}", func(r chi.Router) {
		r.Delete("/", todoController.DeleteTodo)
		r.Patch("/toggle", todoController.ToggleTodo)
		r.Get("/edit", todoController.GetEditToDo)
		r.Patch("/edit", todoController.PatchEditToDo)
	})

	router.Mount("/", subRouter)
	// Serve the embedded static files
	fileServer := http.FileServer(http.FS(staticFiles))
	router.Handle("/static/*", fileServer)

	// Serve the embedded static files
	distFileServer := http.FileServer(http.FS(jsComponents))
	router.Handle("/dist/*", distFileServer)

	port := os.Getenv("PORT")

	fmt.Println("Server running on port", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

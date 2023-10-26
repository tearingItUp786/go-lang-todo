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

	userService := models.NewUserService(db)
	sessionService := models.NewSessionService(db)
	todoService := models.NewBaseModel(db)

	// middleware setup
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	userController := controllers.NewUserController(
		userService,
		sessionService,
	)

	todoController := controllers.NewTodoController(
		controllers.BaseHandlerInput{TodoService: todoService},
	)

	router := chi.NewRouter()

	router.Use(csrfMw)
	router.Use(umw.SetUser)

	router.Get("/signin", userController.GetSignIn)
	router.Post("/signin", userController.ProcessSignIn)
	router.Post("/signout", userController.ProcessSignOut)
	router.Get("/signup", userController.GetSignUp)
	router.Post("/users", userController.ProcessSignUp)

	router.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userController.CurrentUser)
	})

	router.Route("/", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", todoController.GetToDos)
	})

	subRouter := chi.NewRouter()

	subRouter.Route("/", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Post("/bulk-upload", todoController.BulkUpload)
		r.Delete("/delete-all", todoController.DeleteAll)
		r.Post("/new", todoController.NewTodo)

		r.Delete("/{id}", todoController.DeleteTodo)
		r.Patch("/{id}/toggle", todoController.ToggleTodo)
		r.Get("/{id}/edit", todoController.GetEditToDo)
		r.Patch("/{id}/edit", todoController.PatchEditToDo)
	})

	router.Mount("/todo", subRouter)
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

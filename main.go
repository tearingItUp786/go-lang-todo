package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	todoService := models.NewTodoService(db)

	// middleware setup
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	// controller setup
	userController := controllers.NewUserController(
		userService,
		sessionService,
	)

	todoController := controllers.NewTodoController(
		controllers.ToDoBaseHandlerInput{TodoService: todoService},
	)

	googleAuthController := controllers.NewGoogleAuthController(
		os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		userService,
		sessionService,
	)

	// routing setup starts
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(csrfMw)
	router.Use(umw.SetUser)
	fmt.Println(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	router.Mount("/", baseRouter(umw, userController, todoController))
	router.Mount("/auth/google", oauthRouter(*googleAuthController))
	router.Mount("/todo", todoRouter(umw, todoController))
	router.Mount("/users", userRouter(umw, userController))

	// Serve the embedded static files
	fileServer := http.FileServer(http.FS(staticFiles))
	router.Handle("/static/*", fileServer)

	// Serve the embedded js files
	distFileServer := http.FileServer(http.FS(jsComponents))
	router.Handle("/dist/*", distFileServer)

	port := os.Getenv("PORT")

	fmt.Println("Server running on port", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}

func oauthRouter(googleAuthController controllers.MyGoogleOAuth) http.Handler {
	router := chi.NewRouter()
	router.Get("/login", googleAuthController.OauthGoogleLogin)
	router.Get("/callback", googleAuthController.OauthGoogleCallback)
	return router
}

func baseRouter(
	umw controllers.UserMiddleware,
	userController controllers.UserBaseHandler,
	todoController controllers.ToDoBaseHandler,
) http.Handler {
	router := chi.NewRouter()

	router.Get("/signin", userController.GetSignIn)
	router.Post("/signin", userController.ProcessSignIn)
	router.Post("/signout", userController.ProcessSignOut)
	router.Get("/signup", userController.GetSignUp)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/todo", http.StatusFound)
	})

	return router
}

func todoRouter(
	umw controllers.UserMiddleware,
	todoController controllers.ToDoBaseHandler,
) http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", todoController.GetToDos)
		r.Post("/bulk-upload", todoController.BulkUpload)
		r.Delete("/delete-all", todoController.DeleteAll)
		r.Post("/new", todoController.NewTodo)

		r.Delete("/{id}", todoController.DeleteTodo)
		r.Patch("/{id}/toggle", todoController.ToggleTodo)
		r.Get("/{id}/edit", todoController.GetEditToDo)
		r.Patch("/{id}/edit", todoController.PatchEditToDo)
	})

	return router
}

func userRouter(
	umw controllers.UserMiddleware,
	userController controllers.UserBaseHandler,
) http.Handler {
	router := chi.NewRouter()
	router.Post("/", userController.ProcessSignUp)
	router.Route("/", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/me", userController.CurrentUser)
	})

	return router
}

package controllers

import (
	"fmt"
	"net/http"

	"github.com/tearingItUp786/go-lang-todo/context"
	"github.com/tearingItUp786/go-lang-todo/models"
	"github.com/tearingItUp786/go-lang-todo/templates"
	"github.com/tearingItUp786/go-lang-todo/views"
)

type UserBaseHandler struct {
	userService    *models.UserService
	sessionService *models.SessionService
	signInTemplate Template
	signUpTemplate Template
}

type UserBaseHandlerInput struct {
	UserService    *models.UserService
	SessionService *models.SessionService
}

func NewUserController(
	userService *models.UserService,
	sessionService *models.SessionService,
) *UserBaseHandler {
	signInTemplate := views.Must(views.ParseFS(
		templates.FS,
		"signin.gohtml", "unauth-template.gohtml",
	))
	signUpTemplate := views.Must(views.ParseFS(
		templates.FS,
		"signup.gohtml", "unauth-template.gohtml",
	))
	return &UserBaseHandler{
		userService:    userService,
		sessionService: sessionService,
		signInTemplate: signInTemplate,
		signUpTemplate: signUpTemplate,
	}
}

func (ubh UserBaseHandler) GetSignIn(w http.ResponseWriter, r *http.Request) {
	ubh.signInTemplate.Execute(w, r, nil)
}

func (ubh UserBaseHandler) GetSignUp(w http.ResponseWriter, r *http.Request) {
	ubh.signUpTemplate.Execute(w, r, nil)
}

func (ubh UserBaseHandler) ProcessSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := ubh.userService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := ubh.sessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (ubh UserBaseHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (ubh UserBaseHandler) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := ubh.userService.Authenticate(data.Email, data.Password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	session, err := ubh.sessionService.Create(user.ID)
	fmt.Println(session, err)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (ubh UserBaseHandler) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = ubh.sessionService.Delete(token)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

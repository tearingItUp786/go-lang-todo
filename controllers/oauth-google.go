package controllers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tearingItUp786/go-lang-todo/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

type MyGoogleOAuth struct {
	*oauth2.Config
	userService    *models.UserService
	sessionService *models.SessionService
}

func NewGoogleAuthController(
	clientId, clientSecret, redirectUrl string,
	userService *models.UserService,
	sessionService *models.SessionService,
) *MyGoogleOAuth {
	oauthConfig := oauth2.Config{
		RedirectURL: redirectUrl,
		// RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return &MyGoogleOAuth{
		Config:         &oauthConfig,
		userService:    userService,
		sessionService: sessionService,
	}
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func (o MyGoogleOAuth) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := o.generateStateOauthCookie(w)
	u := o.AuthCodeURL(oauthState, oauth2.AccessTypeOffline)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (o MyGoogleOAuth) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := o.getUserDataFromGoogle(r.FormValue("code"), r.Context())
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	// GetOrCreate user using user model
	// the password hash can be null so need to update the table
	// also need to add a google_user_id to the table.
	// and the username will be the email.
	// if there is already a user with that email, send back the appropriate response

	str := string(data)
	fmt.Println(str)
	// Declare a variable of type User to hold the unmarshaled data
	var googleUser GoogleUser

	// Perform the unmarshal operation
	err = json.Unmarshal(data, &googleUser)
	if err != nil {
		log.Fatalf("Error unmarshaling data: %v", err)
	}

	dbUser, err := o.userService.GetUser(googleUser.Email)
	if err != nil {
		// check if the user doesn't exist
		if err == sql.ErrNoRows {
			dbUser, err = o.userService.CreateGoogleUser(googleUser.Email, googleUser.ID)
			if err != nil {
				fmt.Println("WTF one", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Println("WTF two", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	session, err := o.sessionService.Create(dbUser.ID)
	if err != nil {
		fmt.Println("WTF three", err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	fmt.Println("WTF four", session.Token)
	http.Redirect(w, r, "/", http.StatusFound)
	return
}

func (o MyGoogleOAuth) generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	setCookie(w, "oauthstate", state)

	return state
}

func (o MyGoogleOAuth) getUserDataFromGoogle(code string, ctx context.Context) ([]byte, error) {
	token, err := o.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	return io.ReadAll(response.Body)
}

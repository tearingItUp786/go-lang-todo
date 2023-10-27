package controllers

import (
	"fmt"
	"net/http"
	"os"
)

const (
	CookieSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	httpOnly := os.Getenv("COOKIE_SECURE") != "true"
	secure := os.Getenv("COOKIE_SECURE") == "true"

	fmt.Println("cookie secure", secure)
	fmt.Println("http only", httpOnly)

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: httpOnly,
		Secure:   secure,
	}
	return &cookie
}

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}

func readCookie(r *http.Request, name string) (string, error) {
	c, err := r.Cookie(name)
	if err != nil {
		return "", fmt.Errorf("%s: %w", name, err)
	}
	return c.Value, nil
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := newCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}

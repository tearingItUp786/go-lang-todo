package controllers

import "net/http"

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data any)
	ExecuteTemplate(w http.ResponseWriter, r *http.Request, name string, data any)
}

package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("GET /", app.home)

	mux.HandleFunc("GET /snippet/", app.snippetCreate)
	mux.HandleFunc("POST /snippet/", app.snippetCreatePost)
	mux.HandleFunc("GET /snippet/new", app.snippetCreate)
	mux.HandleFunc("GET /snippet/{id}", app.snippetView)
	return app.recoverPanic(mux)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(mux)
}

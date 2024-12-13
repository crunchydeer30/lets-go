package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/", dynamic.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/new", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("GET /snippet/{id}", dynamic.ThenFunc(app.snippetView))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}

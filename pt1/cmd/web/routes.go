package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.noSurf, app.authenticate)
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /", dynamic.ThenFunc(app.home))
	mux.Handle("GET /ping", http.HandlerFunc(ping))

	mux.Handle("GET /snippet/", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/new", protected.ThenFunc(app.snippetCreate))
	mux.Handle("GET /snippet/{id}", dynamic.ThenFunc(app.snippetView))

	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(mux)
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/crunchydeer30/lets-go/internal/models"
	"github.com/crunchydeer30/lets-go/internal/validator"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	latestSnippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = latestSnippets

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *app) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if (err != nil) || (id < 1) {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	flash := app.sessionManager.PopString(r.Context(), "flash")
	data := app.newTemplateData(r)
	data.Snippet = snippet
	data.Flash = flash

	app.render(w, http.StatusOK, "view.html", data)
}

func (app *app) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *app) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckValue(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckValue(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckValue(validator.MaxChars(form.Title, 100), "title", "Title must not be more than 100 characters long")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *app) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *app) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckValue(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckValue(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckValue(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckValue(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckValue(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddValidationError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintln(w, "userSignupPost")
}

func (app *app) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *app) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckValue(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckValue(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckValue(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.IsValid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)

	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddOtherError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/snippet/new", http.StatusSeeOther)
}

func (app *app) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

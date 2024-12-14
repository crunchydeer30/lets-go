package main

import (
	"database/sql"
	"html/template"
	"net/http"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/crunchydeer30/lets-go/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type app struct {
	logger         Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	cfg := LoadConfig()
	logger := NewLogger()

	db, err := sql.Open(cfg.DB.driver, cfg.DB.url)
	if err != nil {
		logger.error.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.error.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &app{
		logger:         *logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		ErrorLog:     app.logger.error,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.info.Printf("Starting server on https://localhost:%s", cfg.Port)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	app.logger.error.Fatal(err)
}

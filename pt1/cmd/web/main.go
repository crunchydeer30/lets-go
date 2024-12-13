package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/crunchydeer30/lets-go/internal/models"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type app struct {
	logger        Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
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

	app := &app{
		logger:        *logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	server := &http.Server{
		Addr:     ":" + cfg.Port,
		ErrorLog: app.logger.error,
		Handler:  app.routes(),
	}

	app.logger.info.Printf("Starting server on http://localhost:%s", cfg.Port)
	err = server.ListenAndServe()
	app.logger.error.Fatal(err)
}

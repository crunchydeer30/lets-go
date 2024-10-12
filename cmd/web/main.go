package main

import (
	"database/sql"
	"net/http"

	"github.com/crunchydeer30/lets-go/internal/models"
	_ "github.com/lib/pq"
)

type app struct {
	logger   Logger
	snippets *models.SnippetModel
}

func main() {
	cfg := LoadConfig()
	logger := NewLogger()

	db, err := sql.Open(cfg.DB.driver, cfg.DB.url)
	if err != nil {
		logger.error.Fatal(err)
	}
	defer db.Close()

	app := &app{
		logger:   *logger,
		snippets: &models.SnippetModel{DB: db},
	}

	db.Ping()

	server := &http.Server{
		Addr:     ":" + cfg.Port,
		ErrorLog: app.logger.error,
		Handler:  app.routes(),
	}

	app.logger.info.Printf("Starting server on http://localhost:%s", cfg.Port)
	err = server.ListenAndServe()
	app.logger.error.Fatal(err)
}

package main

import (
	"flag"
	"net/http"
)

type app struct {
	logger Logger
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	logger := NewLogger()

	app := &app{
		logger: *logger,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: app.logger.error,
		Handler:  app.routes(),
	}

	app.logger.info.Printf("Starting server on %s", *addr)
	err := server.ListenAndServe()
	app.logger.error.Fatal(err)
}

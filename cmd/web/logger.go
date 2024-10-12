package main

import (
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		error: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
	}
}

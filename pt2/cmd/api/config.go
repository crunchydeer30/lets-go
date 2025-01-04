package main

type Config struct {
	port int
	env  string
	db   DbConfig
}

type DbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

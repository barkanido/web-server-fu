package main

import "github.com/gorilla/mux"

type DB = map[string]string

type server struct {
	db     *DB
	router *mux.Router
}

func newServer(db *DB) *server {
	s := &server{}
	s.db = db
	s.router = mux.NewRouter()
	s.routes()
	return s
}

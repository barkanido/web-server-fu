package main

import (
	"math/rand"
	"time"
)

func (s *server) routes() {
	seed := rand.NewSource(time.Now().UnixNano())
	s.router.HandleFunc("/api/", s.handleAPI(rand.New(seed))).Methods("GET")
	s.router.HandleFunc("/greet", s.handleGreet("hello %s")).Methods("POST")
	s.router.HandleFunc("/template", s.handleTemplate()).Methods("POST")
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
	s.router.HandleFunc("/admin", s.adminOnly(s.handleIndex())).Methods("GET")
}

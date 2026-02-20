package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"sink.io/m/src/api"
	"sink.io/m/src/smtp"
	"sink.io/m/src/store"
	"sink.io/m/src/ws"
)

func main() {
	st := store.New()
	hub := ws.NewHub()

	smtpBackend := &smtp.Backend{Store: st, Hub: hub}
	go smtp.StartSMTP(smtpBackend, ":2525")

	h := &api.Handler{Store: st, Hub: hub}
	r := mux.NewRouter()
	h.RegisterRoutes(r)

	log.Println("API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

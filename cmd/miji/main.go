package main

import (
	"log"
	"net/http"

	"github.com/gocs/miji/postgres"
	"github.com/gocs/miji/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:secret@localhost/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	h := web.NewHandler(store)
	log.Fatal(http.ListenAndServe(":3000", h))
}

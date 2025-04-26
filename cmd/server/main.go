package main

import (
	"log"
	"net/http"

	myHttp "github.com/elvis-onobo/go-wallet-api/internal/http"
	"github.com/elvis-onobo/go-wallet-api/pkg/db"
)

func main() {
	db.Init()
	router := myHttp.NewRouter()

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

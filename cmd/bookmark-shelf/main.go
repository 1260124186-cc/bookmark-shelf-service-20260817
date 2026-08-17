package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/transport"
)

func main() {
	repository := store.NewMemoryRepository()
	library := service.NewLibrary(repository)
	handler := transport.NewHandler(library)

	address := os.Getenv("BOOKMARK_SHELF_ADDR")
	if address == "" {
		address = ":8080"
	}

	log.Printf("bookmark shelf listening on %s", address)
	if err := http.ListenAndServe(address, handler); err != nil {
		log.Fatal(err)
	}
}

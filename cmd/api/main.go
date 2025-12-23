package main

import (
	"log"
	"net/http"

	"example.com/prac11TIP/internal/http"
	"example.com/prac11TIP/internal/http/handlers"
	"example.com/prac11TIP/internal/repo"
)

func main() {
	repo := repo.NewNoteRepoMem()
	h := &handlers.Handler{Repo: repo}
	r := httpx.NewRouter(h)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

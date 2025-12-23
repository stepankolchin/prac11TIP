package httpx

import (
	"example.com/prac11TIP/internal/http/handlers"
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Post("/", h.CreateNote)
		r.Get("/", h.GetAllNotes)
		r.Get("/{id}", h.GetNoteByID)
		r.Put("/{id}", h.UpdateNote)
		r.Delete("/{id}", h.DeleteNote)
	})
	return r
}

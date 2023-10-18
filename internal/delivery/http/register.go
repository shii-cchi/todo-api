package http

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (h *Handler) RegisterHTTPEndpoints(r chi.Router) {
	r.Get("/", homeHandler)
	r.Mount("/users", h.userHandlers())
	r.Mount("/todo", h.todoHandlers())
}

func (h *Handler) todoHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Group(func(r chi.Router) {
		r.Get("/", h.handlerFetchTodos)
		r.Get("/{id}", h.handlerFetchTodo)
		r.Post("/", h.handlerCreateTodo)
		r.Put("/{id}", h.handlerUpdateTodo)
		r.Delete("/{id}", h.handlerDeleteTodo)
	})

	return rg
}

func (h *Handler) userHandlers() http.Handler {
	rg := chi.NewRouter()
	rg.Post("/", h.handlerCreateUser)
	return rg
}

package handler

import (
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Handler struct {
	mux     *mux.Router
	log     *slog.Logger
	service PersonService
}

func NewHandler(log *slog.Logger, service PersonService) *Handler {
	h := &Handler{
		mux:     mux.NewRouter(),
		log:     log,
		service: service,
	}

	return h
}

func (h *Handler) initRoutes() {
	h.mux.HandleFunc("/api/v1", h.getAllPerson).Methods(http.MethodGet)
	h.mux.HandleFunc("/api/v1/{id}", h.getPerson).Methods(http.MethodGet)
	h.mux.HandleFunc("/api/v1", h.createPerson).Methods(http.MethodPost)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

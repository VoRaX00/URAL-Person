package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"persons/internal/domain"
	"persons/internal/service"
)

type UserService interface {
	GetAll() ([]domain.GetUser, error)
	GetById(id uuid.UUID) (*domain.GetUser, error)
	Create(user domain.RegisterUser) (uuid.UUID, error)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getUser"
	log := h.log.With(
		slog.String("op", op),
	)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Error("invalid users id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("getting users")
	user, err := h.service.GetById(id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			log.Error("users not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Error("error with getting users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully got users")

	body, err := json.Marshal(user)
	if err != nil {
		log.Error("error marshaling users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) getAllUsers(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getAllUser"
	log := h.log.With(
		slog.String("op", op),
	)

	log.Info("getting users")
	users, err := h.service.GetAll()
	if err != nil {
		log.Error("error with getting users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Info("successfully got users")

	body, err := json.Marshal(users)
	if err != nil {
		log.Error("error with marshaling users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createUser"
	log := h.log.With(
		slog.String("op", op),
	)

	var user domain.RegisterUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Error("error unmarshalling users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("creating users")
	id, err := h.service.Create(user)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyExists) {
			log.Error("users already exists")
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Error("error with creating users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully created users")

	body, err := json.Marshal(id)
	if err != nil {
		log.Error("error marshaling users", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(body)
}

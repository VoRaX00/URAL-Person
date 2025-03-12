package handler

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"persons/internal/domain"
)

type PersonService interface {
	GetAll() ([]domain.GetPerson, error)
	GetById(id uuid.UUID) (*domain.GetPerson, error)
	Create(person domain.RegisterPerson) (uuid.UUID, error)
}

func (h *Handler) getPerson(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getPerson"
	log := h.log.With(
		slog.String("op", op),
	)

	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Error("invalid person id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("getting person")
	person, err := h.service.GetById(id)
	if err != nil {
		// TODO: Обработка ошибки, если пользователь не найден
		log.Error("error with getting person", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully got person")

	body, err := json.Marshal(person)
	if err != nil {
		log.Error("error marshaling person", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) getAllPerson(w http.ResponseWriter, r *http.Request) {
	const op = "handler.getAllPerson"
	log := h.log.With(
		slog.String("op", op),
	)

	log.Info("getting persons")
	persons, err := h.service.GetAll()
	if err != nil {
		log.Error("error with getting persons", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Info("successfully got persons")

	body, err := json.Marshal(persons)
	if err != nil {
		log.Error("error with marshaling persons", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (h *Handler) createPerson(w http.ResponseWriter, r *http.Request) {
	const op = "handler.createPerson"
	log := h.log.With(
		slog.String("op", op),
	)

	var person domain.RegisterPerson
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Error("error unmarshalling person", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("creating person")
	id, err := h.service.Create(person)
	if err != nil {
		log.Error("error with creating person", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Info("successfully created person")

	body, err := json.Marshal(id)
	if err != nil {
		log.Error("error marshaling person", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(body)
}

package person

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"persons/internal/domain"
	"persons/internal/domain/models"
	"persons/internal/service"
	"persons/internal/storage/postgres"
	"persons/pkg/mapper"
)

type Saver interface {
	Save(person models.Person) error
}

type Provider interface {
	GetById(id uuid.UUID) (*models.Person, error)
	GetAll() ([]models.Person, error)
}

type Service struct {
	log      *slog.Logger
	saver    Saver
	provider Provider
}

func NewService(log *slog.Logger,
	saver Saver,
	provider Provider) *Service {
	return &Service{
		log:      log,
		saver:    saver,
		provider: provider,
	}
}

func (s *Service) GetAll() ([]domain.GetPerson, error) {
	const op = "service.GetAll"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting all persons")
	result, err := s.provider.GetAll()
	if err != nil {
		log.Error("failed to get all persons", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("got all persons")

	persons := make([]domain.GetPerson, len(result))
	log.Info("mapping all persons")
	for i, person := range result {
		persons[i] = mapper.PersonToGet(person)
	}
	log.Info("got all persons")
	return persons, nil
}

func (s *Service) GetById(id uuid.UUID) (*domain.GetPerson, error) {
	const op = "service.GetById"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting person by id")
	res, err := s.provider.GetById(id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			log.Error("person not found")
			return nil, fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}
		log.Error("failed to get person by id", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("got person by id")

	person := mapper.PersonToGet(*res)
	return &person, nil
}

func (s *Service) Create(person domain.RegisterPerson) (uuid.UUID, error) {
	const op = "service.Create"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating person")
	result := mapper.RegisterToPerson(person)
	result.Id = uuid.New()

	passHash, err := bcrypt.GenerateFromPassword([]byte(result.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err.Error())
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	result.PasswordHash = string(passHash)

	err = s.saver.Save(result)
	if err != nil {
		if errors.Is(err, postgres.ErrAlreadyExists) {
			log.Error("conflict save person")
			return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrAlreadyExists)
		}
		log.Error("failed to save person", slog.String("err", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("person created")
	return result.Id, nil
}

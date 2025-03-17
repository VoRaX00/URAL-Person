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
	Save(user models.User) error
}

type Provider interface {
	GetById(id uuid.UUID) (*models.User, error)
	GetAll() ([]models.User, error)
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

func (s *Service) GetAll() ([]domain.GetUser, error) {
	const op = "service.GetAll"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting all users")
	result, err := s.provider.GetAll()
	if err != nil {
		log.Error("failed to get all users", slog.String("err", err.Error()))
		return nil, err
	}
	log.Info("got all users")

	users := make([]domain.GetUser, len(result))
	log.Info("mapping all users")
	for i, user := range result {
		users[i] = mapper.UserToGet(user)
	}
	log.Info("got all users")
	return users, nil
}

func (s *Service) GetById(id uuid.UUID) (*domain.GetUser, error) {
	const op = "service.GetById"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("getting users by id")
	res, err := s.provider.GetById(id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			log.Error("users not found")
			return nil, fmt.Errorf("%s: %w", op, service.ErrNotFound)
		}
		log.Error("failed to get users by id", slog.String("err", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("got users by id")

	user := mapper.UserToGet(*res)
	return &user, nil
}

func (s *Service) Create(user domain.RegisterUser) (uuid.UUID, error) {
	const op = "service.Create"
	log := s.log.With(
		slog.String("op", op),
	)

	log.Info("creating users")
	result := mapper.RegisterToUser(user)
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
			log.Error("conflict save users")
			return uuid.Nil, fmt.Errorf("%s: %w", op, service.ErrAlreadyExists)
		}
		log.Error("failed to save users", slog.String("err", err.Error()))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("users created")
	return result.Id, nil
}

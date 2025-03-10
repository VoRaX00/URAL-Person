package person

import "log/slog"

type Repository interface {
}

type Service struct {
	log        *slog.Logger
	repository Repository
}

func NewService(log *slog.Logger, repository Repository) *Service {
	return &Service{
		log:        log,
		repository: repository,
	}
}

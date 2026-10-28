package service

import (
	"context"
	"github.com/google/uuid"
	"kontakt/internal/core/domain"
)

type Repository interface {
	Create(context.Context, domain.OrderInput) (domain.Order, error)
	List(context.Context) ([]domain.Order, error)
	Status(context.Context, uuid.UUID, string) error
}
type Service struct{ Repo Repository }

func (s *Service) Create(c context.Context, in domain.OrderInput) (domain.Order, error) {
	return s.Repo.Create(c, in)
}
func (s *Service) List(c context.Context) ([]domain.Order, error) { return s.Repo.List(c) }
func (s *Service) Status(c context.Context, id uuid.UUID, status string) error {
	return s.Repo.Status(c, id, status)
}

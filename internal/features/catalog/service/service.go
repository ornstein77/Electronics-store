package service

import (
	"context"
	"github.com/google/uuid"
	"kontakt/internal/core/domain"
)

type Repository interface {
	List(context.Context) ([]domain.Product, error)
	Get(context.Context, uuid.UUID) (domain.Product, error)
	Save(context.Context, domain.Product, bool) error
}
type Service struct{ Repo Repository }

func (s *Service) List(c context.Context) ([]domain.Product, error) { return s.Repo.List(c) }
func (s *Service) Get(c context.Context, id uuid.UUID) (domain.Product, error) {
	return s.Repo.Get(c, id)
}
func (s *Service) Save(c context.Context, p domain.Product, create bool) (domain.Product, error) {
	if create {
		p.ID = uuid.New()
	}
	return p, s.Repo.Save(c, p, create)
}

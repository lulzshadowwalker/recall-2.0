package service

import (
	"context"

	"github.com/lulzshadowwalker/recall/internal"
)

type Repository interface {
	Create(ctx context.Context, params internal.CreateMemoryParams) (internal.Memory, error)
}

type service struct {
	Repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Create(ctx context.Context, params internal.CreateMemoryParams) (internal.Memory, error) {
	return s.Repository.Create(ctx, params)
}

package service

import (
	"context"

	"github.com/lulzshadowwalker/recall/internal"
)

type Repository interface {
	Create(ctx context.Context, params internal.CreateMemoryParams) (internal.Memory, error)
}

type service struct {
	r Repository
}

func (s *service) Create(ctx context.Context, params internal.CreateMemoryParams) (internal.Memory, error) {
	return s.r.Create(ctx, params)
}

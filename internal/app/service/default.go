// internal/app/service/default.go
package service

import (
	"context"

	"go-practice/internal/app/repository"

	"go.uber.org/zap"
)

type Service interface {
	GetUser(ctx context.Context, id int64) error
}

type service struct {
	log  *zap.Logger
	repo repository.Repository
}

func NewService(log *zap.Logger, repo repository.Repository) Service {
	return &service{log, repo}
}

func (s *service) GetUser(ctx context.Context, id int64) error {
	return nil
}

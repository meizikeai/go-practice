// internal/app/service/default.go
package service

import (
	"context"

	"go-practice/internal/app/repository"
	"go-practice/internal/config"

	"go.uber.org/zap"
)

type Service interface {
	GetUser(ctx context.Context, id int64) error
}

type service struct {
	env  config.App
	log  *zap.Logger
	repo repository.Repository
}

func NewService(env config.App, log *zap.Logger, repo repository.Repository) Service {
	return &service{env, log, repo}
}

func (s *service) GetUser(ctx context.Context, id int64) error {
	return nil
}

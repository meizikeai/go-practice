// internal/app/service/common.go
package service

import (
	"context"

	"go-practice/internal/app/repository"
)

type service struct {
	env  string
	repo repository.Repository
}

func NewService(env string, repo repository.Repository) Service {
	return &service{env, repo}
}

type Service interface {
	GetUser(ctx context.Context, id int64) error
}

func (s *service) GetUser(ctx context.Context, id int64) error {
	return nil
}

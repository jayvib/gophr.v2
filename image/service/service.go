package service

import (
	"context"
	"gophr.v2/image"
)

func New(repo image.Repository) image.Service {
	return &service{repo: repo}
}

type service struct {
	repo image.Repository
}

func (s *service) Save(ctx context.Context, image *image.Image) error {
	return nil
}

func (s *service) Find(ctx context.Context, id string) (*image.Image, error) {
	return s.repo.Find(ctx, id)
}

func (s *service) FindAll(ctx context.Context, offset int) ([]*image.Image, error) {
	return nil, nil
}

func (s *service) FindAllByUser(ctx context.Context, userId string, offset int) ([]*image.Image, error) {
	return nil, nil
}


package service

import (
	"context"
	"gophr.v2/image"
	"gophr.v2/image/imageutil"
	"gophr.v2/util/valueutil"
	"time"
)

func New(repo image.Repository) image.Service {
	return &service{repo: repo}
}

type service struct {
	repo image.Repository
}

func (s *service) Save(ctx context.Context, image *image.Image) error {
	fillNecessaryField(image)
	return s.repo.Save(ctx, image)
}

func fillNecessaryField(img *image.Image) {
	if img.CreatedAt == nil {
		img.CreatedAt = valueutil.TimePointer(time.Now().UTC())
	}
	if img.ImageID == "" {
		img.ImageID = imageutil.GenerateID()
	}
}

func (s *service) Find(ctx context.Context, id string) (*image.Image, error) {
	return s.repo.Find(ctx, id)
}

func (s *service) FindAll(ctx context.Context, offset int) ([]*image.Image, error) {
	return s.repo.FindAll(ctx, offset)
}

func (s *service) FindAllByUser(ctx context.Context, userId string, offset int) ([]*image.Image, error) {
	return s.repo.FindAllByUser(ctx, userId, offset)
}


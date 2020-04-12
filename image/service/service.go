package service

import (
	"context"
	"github.com/spf13/afero"
	"gophr.v2/image"
	"gophr.v2/image/file"
	"gophr.v2/image/imageutil"
	"gophr.v2/util/valueutil"
	"net/http"
	"time"
)

func New(repo image.Repository, fs afero.Fs) image.Service {
	return &service{
		repo: repo,
		client: http.DefaultClient,
		fs: fs,
	}
}

type service struct {
	repo image.Repository
	client *http.Client
	fs afero.Fs
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

func (s *service) CreateImageFromURL(ctx context.Context, imageUrl string) (*image.Image, error) {
	return nil, nil
}

func (s *service) CreateImageFromFile(ctx context.Context, f file.File, metadata *file.Metadata) (*image.Image, error) {
	return nil, nil
}
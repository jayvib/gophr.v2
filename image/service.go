package image

import (
	"context"
	"io"
)

//go:generate mockery --name=Service

type Service interface {
	Save(ctx context.Context, image *Image) error
	Find(ctx context.Context, id string) (*Image, error)
	FindAll(ctx context.Context, offset int) ([]*Image, error)
	FindAllByUser(ctx context.Context, userId string, offset int) ([]*Image, error)
	CreateImageFromURL(ctx context.Context, url, userId, description string) (*Image, error)
	CreateImageFromFile(ctx context.Context, r io.Reader, filename, description, userId string) (*Image, error)
}

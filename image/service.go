package image

import (
	"context"
	"gophr.v2/image/file"
)

//go:generate mockery --name=Service

type Service interface {
	Save(ctx context.Context, image *Image) error
	Find(ctx context.Context, id string) (*Image, error)
	FindAll(ctx context.Context, offset int) ([]*Image, error)
	FindAllByUser(ctx context.Context, userId string, offset int) ([]*Image, error)
	CreateImageFromURL(ctx context.Context, url string) (*Image, error)
	CreateImageFromFile(ctx context.Context, f file.File, meta *file.Metadata) (*Image, error)
}

package image

import "context"

type Repository interface {
	Save(ctx context.Context, image *Image) error
	Find(ctx context.Context, id string) (*Image, error)
	FindAll(ctx context.Context, offset int) ([]*Image, error)
	FindAllByUser(ctx context.Context, userId string, offset int) ([]*Image, error)
}
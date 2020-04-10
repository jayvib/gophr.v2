package mysql

import (
	"context"
	"gophr.v2/image"
)

type repository struct {}

func (r *repository) Save(ctx context.Context, image *image.Image) error {
	return nil
}

func (r *repository) Find(ctx context.Context, id string) (*image.Image, error) {
	return nil, nil
}

func (r *repository) FindAll(ctx context.Context, offset int) ([]*image.Image, error) {
	return nil, nil
}

func (r *repository) FindAllByUser(ctx context.Context, userId string, offset int) ([]*image.Image, error) {
	return nil, nil
}

package image

import "github.com/pkg/errors"

var (
	ErrNotFound           = errors.New("image: item not found")
	ErrInvalidImageType = errors.New("image: invalid image type")
)

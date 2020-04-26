package image

import "github.com/pkg/errors"

var (
	ErrNotFound           = errors.New("image: item not found")
	ErrInvalidImageType   = errors.New("image: invalid image type")
	ErrInvalidImageURL    = errors.New("image: invalid image url")
	ErrFailedRequest      = errors.New("image: failed while do a client request")
	ErrInvalidContentType = errors.New("image: invalid content-type")
)

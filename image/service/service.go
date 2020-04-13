package service

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"gophr.v2/image"
	"gophr.v2/image/file"
	"gophr.v2/image/imageutil"
	"gophr.v2/util/valueutil"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"time"
)

func New(repo image.Repository, fs afero.Fs, client *http.Client) image.Service {
	return &service{
		repo: repo,
		client: client,
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

func (s *service) CreateImageFromURL(ctx context.Context, imageUrl string, userId string, description string) (*image.Image, error) {
	resp, err := s.client.Get(imageUrl)
	if err != nil {
		return nil, image.ErrInvalidImageURL
	}

	if resp.StatusCode != http.StatusOK {
		return nil, image.ErrFailedRequest
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, image.ErrInvalidContentType
	}

	ext, ok := image.MimeExtensions[mimeType]
	if !ok {
		return nil, image.ErrInvalidContentType
	}
	imageName := filepath.Base(imageUrl)
	imageID := imageutil.GenerateID()
	imageLocation := fmt.Sprintf("%s%s", imageID, ext)

	savedFile, err := s.fs.Create(filepath.Join("./data/images/", imageLocation))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = savedFile.Close()
	}()

	size, err := io.Copy(savedFile, resp.Body)
	if err != nil {
		return nil, err
	}

	img := &image.Image{
		ImageID: imageID,
		UserID: userId,
		Name: imageName,
		Location: imageLocation,
		Description: description,
		Size: size,
	}

	err = s.Save(ctx, img)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (s *service) CreateImageFromFile(ctx context.Context, f file.File, metadata *file.Metadata) (*image.Image, error) {
	return nil, nil
}
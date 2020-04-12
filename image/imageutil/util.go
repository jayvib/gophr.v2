package imageutil

import (
	"gophr.v2/image"
	"gophr.v2/util/randutil"
	"mime"
	"net/http"
)

func GenerateID() string {
	return randutil.GenerateID("image")
}

func GetFileExtensionFromResponse(r *http.Response) (ext string, err error) {
	mimeType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return "", image.ErrInvalidImageType
	}

	ext, ok := image.MimeExtensions[mimeType]
	if !ok {
		return "", image.ErrInvalidImageType
	}
	return ext, nil
}

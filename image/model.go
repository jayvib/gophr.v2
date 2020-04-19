package image

import "time"

var MimeExtensions = map[string]string{
	"image/png":  ".png",
	"image/jpeg": ".jpg",
	"image/gif":  ".gif",
}

type Image struct {
	ID        uint       `json:"id,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" sql:"index"`

	UserID      string `json:"userId,omitempty"`
	ImageID     string `json:"imageId,omitempty"`
	Name        string `json:"name,omitempty"`
	Location    string `json:"location,omitempty"`
	Size        int64  `json:"size,omitempty"`
	Description string `json:"description,omitempty"`
}

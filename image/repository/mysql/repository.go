package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/prometheus/common/log"
	"gophr.v2/image"
)

const pageSize = 25

func New(db *sql.DB) image.Repository {
	return &repository{
		conn: db,
	}
}

type repository struct {
	conn *sql.DB
}

func (r *repository) Save(ctx context.Context, image *image.Image) error {
	query := "INSERT INTO images(userId, imageId, name, location, description, size, created_at, updated_at, deleted_at) VALUES(?,?,?,?,?,?,?,?,?)"
	return r.doSave(func(tx *sql.Tx)error{
		res, err := tx.ExecContext(ctx, query,
			image.UserID,
			image.ImageID,
			image.Name,
			image.Location,
			image.Description,
			image.Size,
			image.CreatedAt,
			image.UpdatedAt,
			image.DeletedAt,
		)
		if err != nil {
			return r.checkError(err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			return r.checkError(err)
		}

		image.ID = uint(id)
		return nil
	})
}

func (r *repository) doSave(fn func(tx *sql.Tx) error) (err error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return r.checkError(err)
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			if e := tx.Rollback(); e != nil {
				// TODO: I think its better to cascade the error
				// For now just log the error.
				log.Error("error while rolling back data in sql:", e)
			}
		}
	}()

	err = fn(tx)
	if err != nil {
		return
	}

	return nil
}

func (r *repository) Find(ctx context.Context, id string) (*image.Image, error) {
	query := "SELECT id, userId, imageId, name, location, description, size, created_at, updated_at, deleted_at FROM images WHERE imageId = ?"
	return r.doQuerySingleReturn(ctx, query, id)
}

func (r *repository) FindAll(ctx context.Context, offset int) ([]*image.Image, error) {
	query := `SELECT id, userId, imageId, name, location, description, size, created_at, updated_at, deleted_at 
						FROM images
						ORDER BY created_at DESC
						LIMIT ?
						OFFSET ?`
	return r.doQuery(ctx, query, pageSize, offset)
}

func (r *repository) FindAllByUser(ctx context.Context, userId string, offset int) ([]*image.Image, error) {
	return nil, nil
}

func (r *repository) checkError(err error) error {
	var cerr error
	switch err {
	case sql.ErrNoRows:
		cerr = image.ErrNotFound
	case nil:
		cerr = nil
	default:
		cerr = fmt.Errorf("mysql: unexpected error %w", err)
	}
	return cerr
}

func (r *repository) doQuery(ctx context.Context, query string, args ...interface{}) (images []*image.Image, err error) {
	row, err := r.conn.QueryContext(ctx, query, args...)
	if err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	defer func() {
		if e := row.Close(); err == nil && e != nil {
			err = e
		}
	}()

	images = make([]*image.Image, 0)
	for row.Next() {
		var img image.Image
		err = row.Scan(&img.ID, &img.UserID, &img.ImageID, &img.Name, &img.Location, &img.Description, &img.Size, &img.CreatedAt, &img.UpdatedAt, &img.DeletedAt)
		if err != nil {
			return nil, r.checkError(err)
		}
		images = append(images, &img)
	}
	if err = row.Err(); err != nil {
		log.Debug(err)
		return nil, r.checkError(err)
	}
	return images, nil
}

func (r *repository) doQuerySingleReturn(ctx context.Context, query string, value interface{}) (img *image.Image, err error) {
	images, err := r.doQuery(ctx, query, value)
	if err != nil {
		return nil, r.checkError(err)
	}

	if len(images) == 0 {
		return nil, image.ErrNotFound
	}
	return images[0], nil
}

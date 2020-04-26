//+build integration

package mysql_test

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gophr.v2/config"
	"gophr.v2/driver/mysql"
	"gophr.v2/image"
	"gophr.v2/image/imageutil"
	mysqlrepo "gophr.v2/image/repository/mysql"
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"os"
	"testing"
	"time"
)

var db *sql.DB

func setup() {
	conf, err := config.New(config.DevelopmentEnv)
	if err != nil {
		panic(err)
	}
	db, err = mysql.InitializeDriver(conf)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	deleteAllInDB()
	err := db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestRepository_Save(t *testing.T) {
	input := &image.Image{
		CreatedAt:   valueutil.TimePointer(time.Now()),
		UserID:      userutil.GenerateID(),
		ImageID:     imageutil.GenerateID(),
		Name:        "Luffy Monkey",
		Location:    "East Blue",
		Size:        1024,
		Description: "A Pirate King from East Blue",
	}

	repo := mysqlrepo.New(db)
	err := repo.Save(context.Background(), input)
	require.NoError(t, err)
	assert.NotEmpty(t, input.ID)
	assertSavedImage(t, input)
}

func TestRepository_Find(t *testing.T) {
	repo := mysqlrepo.New(db)

	t.Run("Image Found", func(t *testing.T) {
		want := &image.Image{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		}
		err := repo.Save(context.Background(), want)
		require.NoError(t, err)

		got, err := repo.Find(context.Background(), want.ImageID)
		require.NoError(t, err)
		assertImage(t, want, got)
	})

	t.Run("Image Not Found", func(t *testing.T) {
		_, err := repo.Find(context.Background(), "notfoundid")
		assert.Error(t, err)
		assert.Equal(t, image.ErrNotFound, err)
	})
}

func TestRepository_FindAll(t *testing.T) {
	// Delete the existing contents
	deleteAllInDB()
	images := []*image.Image{
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Roronoa Zoro",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Sanji Vinsmoke",
			Location:    "West Blue",
			Size:        1024,
			Description: "A Cook from West Blue",
		},
	}

	repo := mysqlrepo.New(db)
	storeImages(t, repo, images)

	got, err := repo.FindAll(context.Background(), 0)
	assert.NoError(t, err)
	assert.Len(t, got, 3)
}

func TestRepository_FindAllByUser(t *testing.T) {
	userId := userutil.GenerateID()
	images := []*image.Image{
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userId,
			ImageID:     imageutil.GenerateID(),
			Name:        "Luffy Monkey",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Pirate King from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userId,
			ImageID:     imageutil.GenerateID(),
			Name:        "Roronoa Zoro",
			Location:    "East Blue",
			Size:        1024,
			Description: "A Swordsman from East Blue",
		},
		{
			CreatedAt:   valueutil.TimePointer(time.Now()),
			UserID:      userutil.GenerateID(),
			ImageID:     imageutil.GenerateID(),
			Name:        "Sanji Vinsmoke",
			Location:    "West Blue",
			Size:        1024,
			Description: "A Cook from West Blue",
		},
	}
	repo := mysqlrepo.New(db)
	storeImages(t, repo, images)
	got, err := repo.FindAllByUser(context.Background(), userId, 0)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
}

func storeImages(t *testing.T, repo image.Repository, images []*image.Image) {
	for _, img := range images {
		err := repo.Save(context.Background(), img)
		require.NoError(t, err)
	}
}

func deleteAllInDB() {
	query := "DELETE FROM images"
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func assertSavedImage(t *testing.T, input *image.Image) {
	query := "SELECT id, userId, imageId, name, location, description, size, created_at, updated_at, deleted_at FROM images WHERE id = ?"
	row, err := db.QueryContext(context.Background(), query, input.ID)
	require.NoError(t, err)
	defer func() {
		err = row.Close()
		require.NoError(t, err)
	}()
	var img image.Image
	for row.Next() {
		err = row.Scan(&img.ID, &img.UserID, &img.ImageID, &img.Name, &img.Location, &img.Description, &img.Size, &img.CreatedAt, &img.UpdatedAt, &img.DeletedAt)
		require.NoError(t, err)
		break
	}
	assertImage(t, input, &img)
}

func assertImage(t *testing.T, want *image.Image, got *image.Image) {
	want.CreatedAt = nil
	got.CreatedAt = nil
	assert.Equal(t, want, got)
}

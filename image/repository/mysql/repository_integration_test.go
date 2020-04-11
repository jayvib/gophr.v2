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
	"gophr.v2/user/userutil"
	"gophr.v2/util/valueutil"
	"os"
	"testing"
	"time"
	mysqlrepo "gophr.v2/image/repository/mysql"
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
	err := db.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestRepository_Save(t *testing.T) {
	input := &image.Image{
		CreatedAt: valueutil.TimePointer(time.Now()),
		UserID: userutil.GenerateID(),
		ImageID: imageutil.GenerateID(),
		Name: "Luffy Monkey",
		Location: "East Blue",
		Size: 1024,
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

	t.Run("Image Found", func(t *testing.T){
		want := &image.Image{
			CreatedAt: valueutil.TimePointer(time.Now()),
			UserID: userutil.GenerateID(),
			ImageID: imageutil.GenerateID(),
			Name: "Luffy Monkey",
			Location: "East Blue",
			Size: 1024,
			Description: "A Pirate King from East Blue",
		}
		err := repo.Save(context.Background(), want)
		require.NoError(t, err)

		got, err := repo.Find(context.Background(), want.ImageID)
		require.NoError(t, err)
		assertImage(t, want, got)
	})

	t.Run("Image Not Found", func(t *testing.T){
		_, err := repo.Find(context.Background(), "notfoundid")
		assert.Error(t, err)
		assert.Equal(t, image.ErrNotFound, err)
	})
}

func assertSavedImage(t *testing.T, input *image.Image) {
	query := "SELECT id, userId, imageId, name, location, description, size, created_at, updated_at, deleted_at FROM images WHERE id = ?"
	row, err := db.QueryContext(context.Background(), query, input.ID)
	require.NoError(t, err)
	defer row.Close()
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


//+build unit

package mysql

import (
	"database/sql"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
)

func setup() (*sql.DB, sqlmock.Sqlmock, *sqlmock.Rows, error) {
	return nil, nil, nil, nil
}

func TestRepository_GetByEmail(t *testing.T) {
}

func TestRepository_GetByID(t *testing.T) {
}

func TestRepository_GetByUsername(t *testing.T) {
}

func TestRepository_Save(t *testing.T) {
}



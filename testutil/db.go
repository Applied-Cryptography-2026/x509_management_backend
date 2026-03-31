package testutil

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDBMock creates a GORM instance backed by a SQL mock.
// Call this in test files to avoid needing a real MySQL instance.
func NewDBMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	gdb.LogMode(true)
	return gdb, mock, nil
}

package mocks

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// NewMockDB initializes a new mock MySQL database instance
func NewMockDB() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		return nil, nil, err
	}
	return sqlx.NewDb(db, "mysql"), mock, nil
}

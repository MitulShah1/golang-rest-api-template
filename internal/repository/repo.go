package repository

import "github.com/MitulShah1/golang-rest-api-template/package/database"

type DBRepository interface {
	// Product Repository
	ProductRepository
	// Category Repository
	CategoryRepository
}

type NewRepository struct {
	db *database.Database
}

// NewDBRepository creates a new instance of the DBRepository interface using the provided database.
// The returned DBRepository implementation is the NewRepository struct, which wraps the provided database.
func NewDBRepository(db *database.Database) DBRepository {
	return &NewRepository{
		db: db,
	}
}

// Package model provides data structures for database entities.
// It includes models for categories, products, and other database objects.
package model

import "time"

type Category struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	ParentID    *int      `db:"parent_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

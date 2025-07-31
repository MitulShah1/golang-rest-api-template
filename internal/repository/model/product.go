// Package model provides data structures for database entities.
// It includes models for categories, products, and other database objects.
package model

import "time"

// Product represents a product entity
type Product struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Stock       int       `db:"stock"`
	CategoryID  int       `db:"category_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

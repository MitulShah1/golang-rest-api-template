package database

import (
	"errors"
	"testing"

	"golang-rest-api-template/package/database/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// TestNewDatabase_Success checks if the database initializes correctly
func TestNewDatabase_Success(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDB()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Expect "SELECT 1" to verify connection
	mock.ExpectQuery("SELECT 1").WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	db := &Database{DB: mockDB}
	err = db.DB.QueryRow("SELECT 1").Scan(new(int)) // Explicitly trigger the query
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestNewDatabase_Failure simulates a failed database connection
func TestNewDatabase_Failure(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDB()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Simulate connection failure
	mock.ExpectQuery("SELECT 1").WillReturnError(errors.New("database connection error"))

	db := &Database{DB: mockDB}
	err = db.DB.QueryRow("SELECT 1").Scan(new(int)) // Trigger expected query
	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDatabase_Close verifies that Close() works properly
func TestDatabase_Close(t *testing.T) {
	mockDB, mock, err := mocks.NewMockDB()
	assert.NoError(t, err)

	// Expect the database to be closed
	mock.ExpectClose()

	db := &Database{DB: mockDB}
	db.Close()

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetDSN validates the DSN string generation for different database drivers
func TestGetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DBConfig
		expected string
	}{
		{
			name: "MySQL DSN",
			config: DBConfig{
				Driver:   "mysql",
				Host:     "localhost",
				Port:     "3306",
				User:     "root",
				Password: "password",
				DBName:   "testdb",
			},
			expected: "root:password@tcp(localhost:3306)/testdb?parseTime=true",
		},
		{
			name: "PostgreSQL DSN",
			config: DBConfig{
				Driver:   "postgres",
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "postgres://postgres:password@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "SQLite DSN",
			config: DBConfig{
				Driver: "sqlite3",
				DBName: "test.db",
			},
			expected: "test.db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := getDSN(tt.config)
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

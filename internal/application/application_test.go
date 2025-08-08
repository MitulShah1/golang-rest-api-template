package application

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApplication(t *testing.T) {
	app := NewApplication()

	assert.NotNil(t, app)
	assert.Equal(t, "go-rest-api-template", app.Name)
	assert.NotNil(t, app.ShutdownChan)
}

func TestApplication_GetLogger(t *testing.T) {
	app := NewApplication()

	// Logger should be nil before initialization
	assert.Nil(t, app.GetLogger())
}

func TestApplication_GetConfig(t *testing.T) {
	app := NewApplication()

	// Config should be nil before initialization
	assert.Nil(t, app.GetConfig())
}

func TestApplication_GetDatabase(t *testing.T) {
	app := NewApplication()

	// Database should be nil before initialization
	assert.Nil(t, app.GetDatabase())
}

func TestApplication_GetCache(t *testing.T) {
	app := NewApplication()

	// Cache should be nil before initialization
	assert.Nil(t, app.GetCache())
}

func TestApplication_GetServer(t *testing.T) {
	app := NewApplication()

	// Server should be nil before initialization
	assert.Nil(t, app.GetServer())
}

// TestApplication_Shutdown_WithoutComponents tests shutdown with nil components
func TestApplication_Shutdown_WithoutComponents(t *testing.T) {
	app := NewApplication()

	// Shutdown should not panic even with nil components
	err := app.Shutdown()
	assert.NoError(t, err)
}

package health

import (
	"testing"

	"golang-rest-api-template/package/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthAPI(t *testing.T) {
	tests := []struct {
		name        string
		inputLogger *logger.Logger
		wantNil     bool
	}{
		{
			name:        "Success with valid logger",
			inputLogger: logger.NewLogger(logger.DefaultOptions()),
			wantNil:     false,
		},
		{
			name:        "Success with nil logger",
			inputLogger: nil,
			wantNil:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := NewHealthAPI(tt.inputLogger)

			if tt.wantNil {
				assert.Nil(t, api)
			} else {
				assert.NotNil(t, api)
				assert.Equal(t, tt.inputLogger, api.logger)
			}
		})
	}
}

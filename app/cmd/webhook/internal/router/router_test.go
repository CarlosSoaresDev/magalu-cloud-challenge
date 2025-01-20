package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInit(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	Init(router, logger)

	tests := []struct {
		method   string
		endpoint string
		expected int
	}{
		{"POST", "/api/v1/paypal/webhook", http.StatusBadRequest},
		{"POST", "/api/v1/stripe/webhook", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		// Action
		req, _ := http.NewRequest(tt.method, tt.endpoint, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, tt.expected, w.Code)
	}
}

package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
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
		{"GET", "/api/v1/currencies", http.StatusBadRequest},
		{"POST", "/api/v1/currencies/convert", http.StatusBadRequest},
		{"GET", "/api/v1/gateways/avaiables", http.StatusOK},
		{"GET", "/api/v1/gateways/transactions", http.StatusOK},
		{"POST", "/api/v1/gateways", http.StatusBadRequest},
		{"GET", "/ping", http.StatusOK},
	}

	for _, tt := range tests {
		// Action
		req, _ := http.NewRequest(tt.method, tt.endpoint, nil)
		req.Header.Set("x-mgc-correlationId", utils.GenerateGUID())
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, tt.expected, w.Code)
	}
}

package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestApiResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should return 200 status code with data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := gin.H{"message": "success"}
		utils.ApiResponse(c, http.StatusOK, data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
	})

	t.Run("should return 400 status code and abort", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		data := gin.H{"error": "bad request"}
		utils.ApiResponse(c, http.StatusBadRequest, data)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"error": "bad request"}`, w.Body.String())
	})
}

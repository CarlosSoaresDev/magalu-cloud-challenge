package currency_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/handlers/currency"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type CurrencyServiceMock struct {
	mock.Mock
}

func (m *CurrencyServiceMock) GetAllCurrency() (*[]string, error) {
	args := m.Called()
	var result *[]string
	if args.Get(0) != nil {
		res := args.Get(0).([]string)
		result = &res
	}
	return result, args.Error(1)
}

func (m *CurrencyServiceMock) ConvertExchangeRate(currency models.CurrencyConvert) (*float64, error) {
	args := m.Called(currency.FromCurrency, currency.ToCurrency, currency.Amount)
	result := args.Get(0).(float64)
	return &result, args.Error(1)
}

func TestGetAllCurrencyHandler_Success(t *testing.T) {

	// Arrange
	gin.SetMode(gin.TestMode)
	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)
	mockCurrencies := []string{"USD", "EUR"}
	mockCurrencyService.On("GetAllCurrency").Return(mockCurrencies, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/currencies", nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.GetAllCurrencyHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestGetAllCurrencyHandler_Failure_GetCorrelationId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/currencies", nil)

	// Action
	handler.GetAllCurrencyHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestGetAllCurrencyHandler_Failure_GetAllCurrency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)
	mockCurrencyService.On("GetAllCurrency").Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/currencies", nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.GetAllCurrencyHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestConvertExchangeRateHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)

	payload := models.CurrencyConvert{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       100,
	}
	expectedResult := 85.0

	mockCurrencyService.On("ConvertExchangeRate", payload.FromCurrency, payload.ToCurrency, payload.Amount).Return(expectedResult, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/currency/convert", utils.ToJSONReader(payload))
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.ConvertExchangeRateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestConvertExchangeRateHandler_Failure_GetCorrelationId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/currency/convert", nil)

	// Action
	handler.ConvertExchangeRateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestConvertExchangeRateHandler_Failure_BindJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Arrange
	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/currency/convert", nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.ConvertExchangeRateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

func TestConvertExchangeRateHandler_Failure_ConvertExchangeRate(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockCurrencyService := new(CurrencyServiceMock)
	mockLogger := zap.NewNop()
	handler := currency.New(mockLogger, mockCurrencyService)

	payload := models.CurrencyConvert{
		FromCurrency: "USD",
		ToCurrency:   "EUR",
		Amount:       100,
	}

	mockCurrencyService.On("ConvertExchangeRate", payload.FromCurrency, payload.ToCurrency, payload.Amount).Return(float64(0), errors.New("conversion error"))

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodPost, "/currency/convert", utils.ToJSONReader(payload))
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.ConvertExchangeRateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCurrencyService.AssertExpectations(t)
}

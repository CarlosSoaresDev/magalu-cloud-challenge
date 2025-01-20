package gateway_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/handlers/gateway"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type GatewayServiceMock struct {
	mock.Mock
}

func (m *GatewayServiceMock) GetAllAvaiablesGateways() []string {
	args := m.Called()
	var result []string
	if args.Get(0) != nil {
		result = args.Get(0).([]string)
	}
	return result
}

func (m *GatewayServiceMock) GetAllTransactionsByDate(date string) (*[]models.Transaction, error) {
	args := m.Called(date)
	var result []models.Transaction
	if args.Get(0) != nil {
		result = args.Get(0).([]models.Transaction)
	}
	return &result, args.Error(1)
}

func (m *GatewayServiceMock) AddTransaction(id string, payment models.Gateway) error {
	args := m.Called(id, payment)
	return args.Error(0)
}

func TestGetAllAvaiablesGateways_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockGatewayService := new(GatewayServiceMock)
	mockLogger := zap.NewNop()
	handler := gateway.New(mockLogger, mockGatewayService)
	mockGateways := []string{"Stripe", "Paypal"}
	mockGatewayService.On("GetAllAvaiablesGateways").Return(mockGateways, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/gateways", nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.GetAllAvaiablesGateways(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockGatewayService.AssertExpectations(t)
}

func TestGetAllAvaiablesGateways_Failure_GetCorrelationId(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)
	mockGatewayService := new(GatewayServiceMock)
	mockLogger := zap.NewNop()
	handler := gateway.New(mockLogger, mockGatewayService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/gateways", nil)

	// Action
	handler.GetAllAvaiablesGateways(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockGatewayService.AssertExpectations(t)
}

func TestGetAllTransactionsByDateHandler_Success(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockGatewayService := new(GatewayServiceMock)
	mockLogger := zap.NewNop()
	handler := gateway.New(mockLogger, mockGatewayService)

	date := "20/01/2025"
	mockTransactions := []models.Transaction{
		{
			Id:       "1",
			Amount:   100,
			Currency: "USD",
			TransactionStatus: []models.TransactionStatus{
				{Status: "pending"},
				{Status: "success"},
			}},
		{
			Id:       "2",
			Amount:   200,
			Currency: "EUR",
			TransactionStatus: []models.TransactionStatus{
				{Status: "pending"},
			}},
	}

	mockGatewayService.On("GetAllTransactionsByDate", date).Return(mockTransactions, nil)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/transactions?date="+date, nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.GetAllTransactionsByDateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockGatewayService.AssertExpectations(t)
}

func TestGetAllTransactionsByDateHandler_Failure_GetCorrelationId(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockGatewayService := new(GatewayServiceMock)
	mockLogger := zap.NewNop()
	handler := gateway.New(mockLogger, mockGatewayService)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/transactions", nil)

	// Action
	handler.GetAllTransactionsByDateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockGatewayService.AssertExpectations(t)
}

func TestGetAllTransactionsByDateHandler_Failure_GetAllTransactionsByDate(t *testing.T) {
	// Arrange
	gin.SetMode(gin.TestMode)

	mockGatewayService := new(GatewayServiceMock)
	mockLogger := zap.NewNop()
	handler := gateway.New(mockLogger, mockGatewayService)

	date := "01_01_2023"
	mockGatewayService.On("GetAllTransactionsByDate", date).Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest(http.MethodGet, "/transactions?date="+date, nil)
	ctx.Request.Header.Set("x-mgc-correlationId", utils.GenerateGUID())

	// Action
	handler.GetAllTransactionsByDateHandler(ctx)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockGatewayService.AssertExpectations(t)
}

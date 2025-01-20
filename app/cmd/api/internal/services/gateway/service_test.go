package gateway

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/services/gateway/provider"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCacheClient struct {
	mock.Mock
}

func (m *MockCacheClient) Get(key string) ([]byte, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return []byte(args.String(0)), args.Error(1)
}

func (m *MockCacheClient) Set(key string, item interface{}, expiration time.Duration) error {
	args := m.Called(key, item, expiration)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (m *MockCacheClient) CheckCache() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockCacheClient) Delete(key string) (*int64, error) {
	args := m.Called(key)
	return args.Get(0).(*int64), args.Error(1)
}

func TestNew(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheClient)

	// Action
	service := New(mockCache)

	// Assert
	if service == nil {
		t.Errorf("Expected service to be non-nil")
	}

	if service.cache != mockCache {
		t.Errorf("Expected cache to be set correctly")
	}
}

func TestGetAllAvaiablesGateways(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	expectedGateways := []string{"gateway1", "gateway2", "gateway3"}
	provider.Providers = map[provider.ProviderType]provider.PaymentGateway{
		provider.ProviderType("gateway1"): nil,
		provider.ProviderType("gateway2"): nil,
		provider.ProviderType("gateway3"): nil,
	}

	// Action
	actualGateways := service.GetAllAvaiablesGateways()

	// Assert
	if len(actualGateways) != len(expectedGateways) {
		t.Errorf("Expected %d gateways, got %d", len(expectedGateways), len(actualGateways))
	}

	for i, gateway := range actualGateways {
		if gateway != expectedGateways[i] {
			t.Errorf("Expected gateway %s, got %s", expectedGateways[i], gateway)
		}
	}
}

func TestAddTransaction_Success(t *testing.T) {

	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	id := "transaction1"
	payment := models.Gateway{
		Amount:   100.0,
		Currency: "USD",
	}

	now := time.Now()
	transactionsByDate := fmt.Sprintf("%s_%s", cache.TransactionsKey, now.Format("02_01_2006"))

	mockCache.On("Get", transactionsByDate).Return(nil, errors.New(cache.ErrCacheMiss.Error()))
	mockCache.On("Set", transactionsByDate, mock.Anything, time.Duration(0)).Return(nil)

	// Action
	err := service.AddTransaction(id, payment)

	// Assert
	assert.NoError(t, err)
	mockCache.AssertExpectations(t)
}

func TestAddTransaction_CacheSetError(t *testing.T) {

	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	id := "transaction1"
	payment := models.Gateway{
		Amount:   100.0,
		Currency: "USD",
	}

	now := time.Now()
	transactionsByDate := fmt.Sprintf("%s_%s", cache.TransactionsKey, now.Format("02_01_2006"))

	mockCache.On("Get", transactionsByDate).Return(nil, errors.New(cache.ErrCacheMiss.Error()))
	mockCache.On("Set", transactionsByDate, mock.Anything, time.Duration(0)).Return(errors.New("cache set error"))

	// Action
	err := service.AddTransaction(id, payment)

	// Assert
	assert.Error(t, err)
	mockCache.AssertExpectations(t)
}

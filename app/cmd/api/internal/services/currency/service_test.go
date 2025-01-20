package currency

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
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
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	assert.NotNil(t, service)
	assert.Equal(t, mockCache, service.cache)
}

func TestGetAllCurrency_Success(t *testing.T) {
	// Arrange
	mockResponseStr := `{
		"rates": {
			"USD": 1.0,
			"EUR": 0.85,
			"JPY": 110.0
		}
	}`
	server := mockServer(mockResponseStr)
	os.Setenv("OPEN_EXCHANGE_RATES_SECRET_KEY", "test_api_key")
	os.Setenv("OPEN_EXCHANGE_RATES_URL", server.URL+"/%s")

	mockCache := new(MockCacheClient)
	service := New(mockCache)

	mockCache.On("Get", cache.ExchangeRateKey).Return(nil, errors.New("cache miss"))
	mockCache.On("Set", cache.ExchangeRateKey, mock.Anything, time.Minute*5).Return(nil)

	// Action
	currencies, err := service.GetAllCurrency()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, currencies)
	assert.Equal(t, []string{"EUR", "JPY", "USD"}, *currencies)
	mockCache.AssertExpectations(t)
}

func TestGetAllCurrency_CacheHint(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	mockResponse := &models.CurrencyDataResponse{
		Rates: map[string]float64{
			"USD": 1.0,
			"EUR": 0.85,
			"JPY": 110.0,
		},
	}

	mockCache.On("Get", cache.ExchangeRateKey).Return(utils.ToJSON(mockResponse), nil)

	// Action
	currencies, err := service.GetAllCurrency()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, currencies)
	assert.Equal(t, []string{"EUR", "JPY", "USD"}, *currencies)
	mockCache.AssertExpectations(t)
}

func TestGetAllCurrency_Failure_GetAndSerializerData(t *testing.T) {

	// Arrange
	os.Setenv("OPEN_EXCHANGE_RATES_SECRET_KEY", "test_api_key")
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	mockCache.On("Get", cache.ExchangeRateKey).Return("", errors.New("cache miss"))
	mockCache.On("Set", cache.ExchangeRateKey, mock.Anything, time.Minute*5).Return(errors.New("cache set error"))

	// Action
	currencies, err := service.GetAllCurrency()

	// Assert
	assert.Error(t, err)
	assert.Nil(t, currencies)
	mockCache.AssertExpectations(t)
}

func TestGetRates_Success(t *testing.T) {
	// Arrange
	mockResponse := `{
		"rates": {
			"USD": 1.0,
			"EUR": 0.9
		}
	}`
	server := mockServer(mockResponse)

	apiKey := "test_api_key"
	// Action
	data, err := getRates(server.URL+"/%s", apiKey)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 1.0, data.Rates["USD"])
	assert.Equal(t, 0.9, data.Rates["EUR"])
}

func TestConvertExchangeRate_SuccessfulConversion(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	mockCache.On("Get", cache.ExchangeRateKey).Return(`{"rates":{"USD":1.0,"EUR":0.85}}`, nil)

	currency := models.CurrencyConvert{
		Amount:       100,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	}

	// Action
	result, err := service.ConvertExchangeRate(currency)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 85.0, *result)
}

func TestConvertExchangeRate_MissingCurrencyKey(t *testing.T) {
	// Arrange
	mockCache := new(MockCacheClient)
	service := New(mockCache)

	mockCache.On("Get", cache.ExchangeRateKey).Return(`{"rates":{"USD":1.0}}`, nil)

	currency := models.CurrencyConvert{
		Amount:       100,
		FromCurrency: "USD",
		ToCurrency:   "EUR",
	}

	// Action
	result, err := service.ConvertExchangeRate(currency)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "missing or unavailable currency keys: [EUR]", err.Error())
}

func TestGetRates_Failure(t *testing.T) {
	// Arrange
	server := mockServer("")
	apiKey := "test_api_key"

	// Action
	data, err := getRates(server.URL+"/%s", apiKey)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, data)
}

func mockServer(mockResponse string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))

	return server
}

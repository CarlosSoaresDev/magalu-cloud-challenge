package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/CarlosSoaresDev/magalu-cloud-challenge/cmd/api/internal/models"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/cache"
	"github.com/CarlosSoaresDev/magalu-cloud-challenge/pkg/utils"
)

type CurrencyService interface {
	GetAllCurrency() (*[]string, error)
	ConvertExchangeRate(currency models.CurrencyConvert) (*float64, error)
}

type currencyService struct {
	cache cache.CacheClient
}

// New creates a new instance of currencyService with the provided cache client.
// It returns a pointer to the newly created currencyService.
//
// Parameters:
//   - cache: an instance of cache.CacheClient used for caching.
//
// Returns:
//   - *currencyService: a pointer to the initialized currencyService.
func New(cache cache.CacheClient) *currencyService {
	return &currencyService{
		cache: cache,
	}
}

// GetAllCurrency retrieves all available currency codes.
// It fetches and deserializes the data, then extracts the currency codes
// from the response and sorts them alphabetically.
//
// Returns a pointer to a slice of currency codes or an error if the operation fails.
func (p *currencyService) GetAllCurrency() (*[]string, error) {
	res, err := p.getAndSerializerData()
	if err != nil {
		return nil, err
	}

	currencies := make([]string, 0, len(res.Rates))
	for currency := range res.Rates {
		currencies = append(currencies, currency)
	}

	sort.Strings(currencies)
	return &currencies, nil
}

// ConvertExchangeRate converts the amount from one currency to another based on the exchange rates.
// It takes a CurrencyConvert model as input which contains the amount to be converted and the source and target currencies.
// It returns the converted amount as a pointer to float64 and an error if any occurs during the conversion process.
//
// The function performs the following steps:
// 1. Retrieves and serializes the exchange rate data.
// 2. Checks for missing keys in the exchange rate data for the source and target currencies.
// 3. Converts the amount using the exchange rates for the source and target currencies.
//
// Parameters:
// - currency: A models.CurrencyConvert struct containing the amount, source currency, and target currency.
//
// Returns:
// - A pointer to the converted amount as float64.
// - An error if any issue occurs during the retrieval of exchange rates or the conversion process.
func (p *currencyService) ConvertExchangeRate(currency models.CurrencyConvert) (*float64, error) {

	res, err := p.getAndSerializerData()
	if err != nil {
		return nil, err
	}

	err = checkMissingKeys(res.Rates, currency.FromCurrency, currency.ToCurrency)
	if err != nil {
		return nil, err
	}

	amount := convert(currency.Amount, res.Rates[currency.FromCurrency], res.Rates[currency.ToCurrency])
	return &amount, nil
}

// GetAndSerializerData retrieves currency data from the cache or an external service,
// serializes it, and stores it in the cache if not already present.
// It returns the currency data response or an error if any operation fails.
//
// Returns:
//   - *models.CurrencyDataResponse: The currency data response.
//   - error: An error if any operation fails.
func (p currencyService) getAndSerializerData() (*models.CurrencyDataResponse, error) {

	c, err := p.cache.Get(cache.ExchangeRateKey)
	var res *models.CurrencyDataResponse

	if err != nil {
		secretKey, err := getSecretKey()
		if err != nil {
			return nil, err
		}

		ulr, err := getUlr()
		if err != nil {
			return nil, err
		}

		res, err = getRates(ulr, secretKey)
		if err != nil {
			return nil, err
		}

		ratesSerializer, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		if err = p.cache.Set(cache.ExchangeRateKey, ratesSerializer, time.Minute*5); err != nil {
			return nil, err
		}
	} else {
		if err = json.Unmarshal([]byte(c), &res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// getRates fetches the latest currency exchange rates from the Open Exchange Rates API.
// It takes an API key as a parameter and returns a CurrencyDataResponse object containing
// the exchange rates, or an error if the request fails or the response cannot be decoded.
//
// Parameters:
//   - apiKey: A string containing the API key for accessing the Open Exchange Rates API.
//
// Returns:
//   - *models.CurrencyDataResponse: A pointer to a CurrencyDataResponse object containing
//     the latest exchange rates.
//   - error: An error object if the request fails or the response cannot be decoded.
func getRates(urlBase, apiKey string) (*models.CurrencyDataResponse, error) {

	url := fmt.Sprintf(urlBase, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	var data models.CurrencyDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// convert converts an amount from one currency to another using the provided exchange rates.
// It takes three parameters:
// - amount: the amount of money to be converted.
// - rateFrom: the exchange rate of the original currency.
// - rateTo: the exchange rate of the target currency.
// It returns the converted amount in the target currency.
func convert(amount float64, rateFrom, rateTo float64) float64 {
	return amount * (rateTo / rateFrom)
}

// getSecretKey retrieves the Open Exchange Rates secret key from the environment variables.
// It returns the secret key as a string and an error if the key is not set or is empty.
//
// Returns:
//   - string: The secret key retrieved from the environment variables.
//   - error: An error if the secret key is not set or is empty.
func getSecretKey() (string, error) {
	secretKey := os.Getenv("OPEN_EXCHANGE_RATES_SECRET_KEY")
	if utils.IsEmptyOrNull(secretKey) {
		return "", fmt.Errorf("secret key is not set or is empty")
	}
	return secretKey, nil
}

// getUlr retrieves the URL for the Open Exchange Rates service from the environment variables.
// It returns the URL as a string and an error if the URL is not set or is empty.
//
// Returns:
//   - string: The URL for the Open Exchange Rates service.
//   - error: An error if the URL is not set or is empty.
func getUlr() (string, error) {

	secretKey := os.Getenv("OPEN_EXCHANGE_RATES_URL")
	if utils.IsEmptyOrNull(secretKey) {
		return "", fmt.Errorf("secret key is not set or is empty")
	}
	return secretKey, nil
}

// checkMissingKeys checks if the provided map contains all the specified keys.
// It returns an error if any of the keys are missing.
//
// Parameters:
//   - m: A map where the keys are strings and the values are float64.
//   - keys: A variadic parameter representing the keys to check in the map.
//
// Returns:
//   - error: An error indicating which keys are missing, or nil if all keys are present.
func checkMissingKeys(m map[string]float64, keys ...string) error {
	var missingKeys []string
	for _, key := range keys {
		if _, exists := m[key]; !exists {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing or unavailable currency keys: %v", missingKeys)
	}

	return nil
}
